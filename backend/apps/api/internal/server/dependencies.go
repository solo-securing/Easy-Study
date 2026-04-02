package server

import (
	"log/slog"
	"os"

	"api/internal/database"
	"api/internal/handlers"
	"api/internal/services"
)

type Dependencies struct {
	TenantHandler         *handlers.TenantHandler
	UserHandler           *handlers.UserHandler
	CourseHandler         *handlers.CourseHandler
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
	subdomainService := services.NewSubdomainService([]string{"www", "api", "admin"})
	auditService := services.NewAuditService(logger)
	tenantService := services.NewTenantService(repo, subdomainService, auditService)
	activationService := services.NewActivationService()
	userService := services.NewUserService(userRepo, activationService, auditService)
	importService := services.NewUserImportService()
	courseService := services.NewCourseService(courseRepo, auditService)
	enrollmentService := services.NewEnrollmentService(enrollmentRepo, userRepo, groupRepo, courseRepo, auditService)
	tenantSettingsService := services.NewTenantSettingsService()

	return &Dependencies{
		TenantHandler:         handlers.NewTenantHandler(tenantService),
		UserHandler:           handlers.NewUserHandler(userService, groupRepo, importService),
		CourseHandler:         handlers.NewCourseHandler(courseService, enrollmentService),
		TenantSettingsHandler: handlers.NewTenantSettingsHandler(tenantSettingsService),
		ActivationService:     activationService,
	}
}
