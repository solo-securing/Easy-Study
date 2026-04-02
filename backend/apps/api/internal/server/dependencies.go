package server

import (
	"log/slog"
	"os"

	"api/internal/database"
	"api/internal/handlers"
	"api/internal/services"
)

type Dependencies struct {
	TenantHandler *handlers.TenantHandler
}

func NewDependencies(logger *slog.Logger) *Dependencies {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	repo := database.NewInMemoryTenantRepository()
	subdomainService := services.NewSubdomainService([]string{"www", "api", "admin"})
	auditService := services.NewAuditService(logger)
	tenantService := services.NewTenantService(repo, subdomainService, auditService)

	return &Dependencies{
		TenantHandler: handlers.NewTenantHandler(tenantService),
	}
}
