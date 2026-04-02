package server

import (
	"context"
	"log/slog"
	"os"

	"api/internal/database"
	"api/internal/handlers"
	"api/internal/models"
	"api/internal/services"
)

type Dependencies struct {
	TenantHandler         *handlers.TenantHandler
	UserHandler           *handlers.UserHandler
	CourseHandler         *handlers.CourseHandler
	QuizHandler           *handlers.QuizHandler
	StudentHandler        *handlers.StudentHandler
	ImpersonationHandler  *handlers.ImpersonationHandler
	ReportHandler         *handlers.ReportHandler
	TenantSettingsHandler *handlers.TenantSettingsHandler
	ActivationService     *services.ActivationService
}

func NewDependencies(logger *slog.Logger) *Dependencies {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	repo := database.NewInMemoryTenantRepository()
	userRepo := database.NewInMemoryUserRepository()
	groupRepo := database.NewInMemoryGroupRepository()
	courseRepo := database.NewInMemoryCourseRepository()
	enrollmentRepo := database.NewInMemoryEnrollmentRepository()
	quizRepo := database.NewInMemoryQuizRepository()
	impersonationRepo := database.NewInMemoryImpersonationRepository()
	usageRepo := database.NewInMemoryUsageMetricsRepository()
	subdomainService := services.NewSubdomainService([]string{"www", "api", "admin"})
	auditService := services.NewAuditService(logger)
	tenantService := services.NewTenantService(repo, subdomainService, auditService)
	activationService := services.NewActivationService()
	userService := services.NewUserService(userRepo, activationService, auditService)
	importService := services.NewUserImportService()
	courseService := services.NewCourseService(courseRepo, auditService)
	enrollmentService := services.NewEnrollmentService(enrollmentRepo, userRepo, groupRepo, courseRepo, auditService)
	quizService := services.NewQuizService(quizRepo, auditService)
	progressService := services.NewProgressService(enrollmentRepo)
	impersonationService := services.NewImpersonationService(impersonationRepo, auditService)
	reportingService := services.NewReportingService(usageRepo)
	tenantSettingsService := services.NewTenantSettingsService()

	_ = quizService.SeedQuiz(context.Background(), models.Quiz{
		QuizID:           "quiz-001",
		TenantID:         "tenant-001",
		CourseID:         "course-001",
		Title:            "Sample Quiz",
		MaxScore:         100,
		PassingScore:     60,
		TimeLimitMinutes: 30,
		AttemptLimit:     3,
		ScoringPolicy:    models.QuizScoringPolicyHighestValidAttempt,
	})
	_ = quizService.SeedQuiz(context.Background(), models.Quiz{
		QuizID:           "quiz-timeout",
		TenantID:         "tenant-001",
		CourseID:         "course-001",
		Title:            "Timeout Quiz",
		MaxScore:         100,
		PassingScore:     60,
		TimeLimitMinutes: 0,
		AttemptLimit:     1,
		ScoringPolicy:    models.QuizScoringPolicyHighestValidAttempt,
	})
	_ = reportingService.SeedTenantMetric(context.Background(), models.UsageMetricDaily{
		TenantID:          "tenant-001",
		ActiveUserCount:   120,
		ActiveCourseCount: 8,
		CompletionRate:    73.5,
		LoginCount:        460,
		LearningMinutes:   12800,
	})

	return &Dependencies{
		TenantHandler:         handlers.NewTenantHandler(tenantService),
		UserHandler:           handlers.NewUserHandler(userService, groupRepo, importService),
		CourseHandler:         handlers.NewCourseHandler(courseService, enrollmentService),
		QuizHandler:           handlers.NewQuizHandler(quizService),
		StudentHandler:        handlers.NewStudentHandler(progressService),
		ImpersonationHandler:  handlers.NewImpersonationHandler(impersonationService),
		ReportHandler:         handlers.NewReportHandler(reportingService),
		TenantSettingsHandler: handlers.NewTenantSettingsHandler(tenantSettingsService),
		ActivationService:     activationService,
	}
}
