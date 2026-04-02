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
