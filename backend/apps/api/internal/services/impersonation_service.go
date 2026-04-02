package services

import (
	"context"
	"errors"
	"strings"

	"api/internal/database"
	"api/internal/models"
)

var (
	ErrImpersonationClosedSession = errors.New("impersonation session already closed")
)

type ImpersonationService struct {
	repo  database.ImpersonationRepository
	audit *AuditService
}

func NewImpersonationService(repo database.ImpersonationRepository, audit *AuditService) *ImpersonationService {
	return &ImpersonationService{
		repo:  repo,
		audit: audit,
	}
}

func (s *ImpersonationService) Start(ctx context.Context, superAdminID string, input models.ImpersonationStartInput) (models.ImpersonationSession, error) {
	input.TenantID = strings.TrimSpace(input.TenantID)
	input.Reason = strings.TrimSpace(input.Reason)
	if input.TenantID == "" {
		return models.ImpersonationSession{}, errors.New("tenantId is required")
	}
	session, err := s.repo.Start(ctx, models.ImpersonationSession{
		TenantID:     input.TenantID,
		SuperAdminID: superAdminID,
		Mode:         models.ImpersonationModeReadConfigNonDestructive,
		Reason:       input.Reason,
		Actions:      []string{},
	})
	if err != nil {
		return models.ImpersonationSession{}, err
	}
	if s.audit != nil {
		s.audit.Write(AuditRecord{
			TenantID: input.TenantID,
			ActorID:  superAdminID,
			Action:   "impersonation.start",
			Resource: "impersonation_session",
			Result:   "success",
		})
	}
	return session, nil
}

func (s *ImpersonationService) End(ctx context.Context, superAdminID, sessionID string) error {
	session, err := s.repo.Get(ctx, sessionID)
	if err != nil {
		return err
	}
	if session.EndedAt != nil {
		return ErrImpersonationClosedSession
	}
	if err := s.repo.End(ctx, sessionID); err != nil {
		return err
	}
	if s.audit != nil {
		s.audit.Write(AuditRecord{
			TenantID: session.TenantID,
			ActorID:  superAdminID,
			Action:   "impersonation.end",
			Resource: "impersonation_session",
			Result:   "success",
		})
	}
	return nil
}

func (s *ImpersonationService) TrackAction(ctx context.Context, sessionID, action string) error {
	if strings.TrimSpace(action) == "" {
		return errors.New("action is required")
	}
	return s.repo.TrackAction(ctx, sessionID, action)
}

func (s *ImpersonationService) GetAudit(ctx context.Context, tenantID string) (models.ImpersonationAuditList, error) {
	return s.repo.ListByTenant(ctx, tenantID)
}
