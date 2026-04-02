package server

import (
	"net/http"
	"os"

	"api/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	jwtSecret := os.Getenv("JWT_SECRET")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))
	r.Use(middleware.RequestLogger(s.logger))

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	api := r.Group("/api/v1")
	api.Use(middleware.TenantContext())
	api.Use(middleware.JWTAuth(jwtSecret))
	api.Use(middleware.BlockDestructiveWhenImpersonating())

	superAdmin := api.Group("/super-admin")
	superAdmin.Use(middleware.RequireRoles("super_admin"))
	superAdmin.GET("/health", s.healthHandler)
	superAdmin.GET("/tenants", s.deps.TenantHandler.ListTenants)
	superAdmin.POST("/tenants", s.deps.TenantHandler.CreateTenant)
	superAdmin.GET("/tenants/:tenantId", s.deps.TenantHandler.GetTenantByID)
	superAdmin.PATCH("/tenants/:tenantId/status", s.deps.TenantHandler.UpdateTenantStatus)

	tenantAdmin := api.Group("/tenant-admin")
	tenantAdmin.Use(middleware.RequireRoles("tenant_admin"))
	tenantAdmin.GET("/health", s.healthHandler)

	tenantScoped := api.Group("/tenants/:tenantId")
	tenantScoped.Use(middleware.RequireRoles("tenant_admin", "super_admin", "instructor"))
	tenantScoped.GET("/users", s.deps.UserHandler.ListUsers)
	tenantScoped.POST("/users", s.deps.UserHandler.CreateUser)
	tenantScoped.GET("/users/:userId", s.deps.UserHandler.GetUser)
	tenantScoped.PATCH("/users/:userId", s.deps.UserHandler.PatchUser)
	tenantScoped.POST("/users/:userId/invite-resend", s.deps.UserHandler.ResendInvitation)
	tenantScoped.POST("/users/import-jobs", s.deps.UserHandler.CreateImportJob)
	tenantScoped.GET("/users/import-jobs/:jobId", s.deps.UserHandler.GetImportJob)
	tenantScoped.GET("/groups", s.deps.UserHandler.ListGroups)
	tenantScoped.POST("/groups", s.deps.UserHandler.CreateGroup)
	tenantScoped.PATCH("/groups/:groupId", s.deps.UserHandler.PatchGroup)
	tenantScoped.PUT("/groups/:groupId/members", s.deps.UserHandler.PutGroupMembers)
	tenantScoped.GET("/courses", s.deps.CourseHandler.ListCourses)
	tenantScoped.POST("/courses", s.deps.CourseHandler.CreateCourse)
	tenantScoped.GET("/courses/:courseId", s.deps.CourseHandler.GetCourse)
	tenantScoped.PATCH("/courses/:courseId", s.deps.CourseHandler.PatchCourse)
	tenantScoped.PUT("/courses/:courseId/structure", s.deps.CourseHandler.PutCourseStructure)
	tenantScoped.POST("/courses/:courseId/publish", s.deps.CourseHandler.PublishCourse)
	tenantScoped.POST("/courses/:courseId/archive", s.deps.CourseHandler.ArchiveCourse)
	tenantScoped.GET("/courses/:courseId/enrollments", s.deps.CourseHandler.ListEnrollments)
	tenantScoped.POST("/courses/:courseId/enrollments", s.deps.CourseHandler.AssignEnrollments)
	tenantScoped.GET("/reports/dashboard", s.deps.ReportHandler.TenantDashboard)
	tenantScoped.GET("/reports/courses/:courseId", s.deps.ReportHandler.CourseReport)
	tenantScoped.GET("/branding", s.deps.TenantSettingsHandler.GetBranding)
	tenantScoped.PATCH("/branding", s.deps.TenantSettingsHandler.PatchBranding)
	tenantScoped.GET("/settings", s.deps.TenantSettingsHandler.GetSettings)
	tenantScoped.PATCH("/settings", s.deps.TenantSettingsHandler.PatchSettings)

	quizReaders := api.Group("/tenants/:tenantId")
	quizReaders.Use(middleware.RequireRoles("student", "instructor", "tenant_admin", "super_admin"))
	quizReaders.GET("/quizzes/:quizId", s.deps.QuizHandler.GetQuiz)
	quizReaders.GET("/quizzes/:quizId/attempts/:attemptId", s.deps.QuizHandler.GetAttempt)

	studentScoped := api.Group("/tenants/:tenantId")
	studentScoped.Use(middleware.RequireRoles("student"))
	studentScoped.POST("/quizzes/:quizId/attempts", s.deps.QuizHandler.StartAttempt)
	studentScoped.PATCH("/quizzes/:quizId/attempts/:attemptId/autosave", s.deps.QuizHandler.Autosave)
	studentScoped.POST("/quizzes/:quizId/attempts/:attemptId/submit", s.deps.QuizHandler.Submit)
	studentScoped.GET("/courses/:courseId/progress/me", s.deps.StudentHandler.GetMyCourseProgress)

	me := api.Group("/me")
	me.Use(middleware.RequireRoles("student", "instructor", "tenant_admin"))
	me.GET("/courses", s.deps.StudentHandler.ListMyCourses)

	adminImpersonation := api.Group("/admin")
	adminImpersonation.Use(middleware.RequireRoles("super_admin"))
	adminImpersonation.POST("/impersonations", s.deps.ImpersonationHandler.Start)
	adminImpersonation.DELETE("/impersonations/:sessionId", s.deps.ImpersonationHandler.End)

	adminReports := api.Group("/admin")
	adminReports.Use(middleware.RequireRoles("super_admin"))
	adminReports.GET("/reports/tenants-usage", s.deps.ReportHandler.AdminUsage)

	tenantAudit := api.Group("/tenants/:tenantId")
	tenantAudit.Use(middleware.RequireRoles("tenant_admin"))
	tenantAudit.GET("/impersonation-audit", s.deps.ImpersonationHandler.ListAudit)

	auth := api.Group("/auth")
	auth.POST("/activation/verify", func(c *gin.Context) {
		var req struct {
			Token string `json:"token"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || req.Token == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
			return
		}
		token, err := s.deps.ActivationService.Verify(c.Request.Context(), req.Token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusGone, gin.H{"error": err.Error()})
			return
		}
		if err := s.deps.ActivationService.Consume(c.Request.Context(), req.Token); err != nil {
			c.AbortWithStatusJSON(http.StatusGone, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"userId":   token.UserID,
			"tenantId": token.TenantID,
			"status":   "activated",
		})
	})
	auth.POST("/activation/resend", func(c *gin.Context) {
		var req struct {
			TenantID string `json:"tenantId"`
			UserID   string `json:"userId"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || req.TenantID == "" || req.UserID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
			return
		}
		_, _, err := s.deps.ActivationService.Issue(c.Request.Context(), req.TenantID, req.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"status": "queued"})
	})

	instructor := api.Group("/instructor")
	instructor.Use(middleware.RequireRoles("instructor"))
	instructor.GET("/health", s.healthHandler)

	student := api.Group("/student")
	student.Use(middleware.RequireRoles("student"))
	student.GET("/health", s.healthHandler)

	public := r.Group("/public")
	if s.rateLimiter != nil {
		public.Use(s.rateLimiter.Middleware("public"))
	}
	public.GET("/health", s.healthHandler)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
