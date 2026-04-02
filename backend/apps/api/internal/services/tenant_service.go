package services

import (
	"context"
	"errors"
	"strings"

	"api/internal/database"
	"api/internal/models"
)

var (
	ErrInvalidTenantStatusTransition = errors.New("invalid tenant status transition")
)

type TenantService struct {
	repo      database.TenantRepository
	subdomain *SubdomainService
	audit     *AuditService
}

func NewTenantService(repo database.TenantRepository, subdomain *SubdomainService, audit *AuditService) *TenantService {
	return &TenantService{
		repo:      repo,
		subdomain: subdomain,
		audit:     audit,
	}
}

func (s *TenantService) List(ctx context.Context, filter models.TenantListFilter) (models.TenantList, error) {
	return s.repo.List(ctx, filter)
}

func (s *TenantService) Create(ctx context.Context, actorID string, input models.TenantCreateInput) (models.Tenant, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Subdomain = strings.ToLower(strings.TrimSpace(input.Subdomain))
	input.OwnerEmail = strings.ToLower(strings.TrimSpace(input.OwnerEmail))
	input.PlanCode = strings.TrimSpace(input.PlanCode)

	if input.Name == "" || input.Subdomain == "" || input.OwnerEmail == "" || input.PlanCode == "" {
		return models.Tenant{}, errors.New("missing required tenant fields")
	}

	if err := s.subdomain.Validate(ctx, input.Subdomain); err != nil {
		return models.Tenant{}, err
	}

	tenant, err := s.repo.Create(ctx, input)
	if err != nil {
		return models.Tenant{}, err
	}

	if s.audit != nil {
		s.audit.Write(AuditRecord{
			TenantID: tenant.ID,
			ActorID:  actorID,
			Action:   "tenant.create",
			Resource: "tenant",
			Result:   "success",
		})
	}

	return tenant, nil
}

func (s *TenantService) GetByID(ctx context.Context, tenantID string) (models.Tenant, error) {
	return s.repo.GetByID(ctx, tenantID)
}

func (s *TenantService) UpdateStatus(ctx context.Context, actorID, tenantID string, patch models.TenantStatusPatch) (models.Tenant, error) {
	current, err := s.repo.GetByID(ctx, tenantID)
	if err != nil {
		return models.Tenant{}, err
	}

	if !isValidTenantTransition(current.Status, patch.Status) {
		return models.Tenant{}, ErrInvalidTenantStatusTransition
	}

	updated, err := s.repo.UpdateStatus(ctx, tenantID, patch)
	if err != nil {
		return models.Tenant{}, err
	}

	if s.audit != nil {
		s.audit.Write(AuditRecord{
			TenantID: tenantID,
			ActorID:  actorID,
			Action:   "tenant.update_status",
			Resource: "tenant",
			Result:   "success",
		})
	}

	return updated, nil
}

func isValidTenantTransition(from, to models.TenantStatus) bool {
	if from == to {
		return true
	}

	switch from {
	case models.TenantStatusActive:
		return to == models.TenantStatusSuspended || to == models.TenantStatusDeactivated
	case models.TenantStatusSuspended:
		return to == models.TenantStatusActive || to == models.TenantStatusDeactivated
	case models.TenantStatusDeactivated:
		return false
	default:
		return false
	}
}
