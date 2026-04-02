package services

import (
	"log/slog"
	"time"
)

type AuditRecord struct {
	TenantID string
	ActorID  string
	Action   string
	Resource string
	Result   string
	At       time.Time
}

type AuditService struct {
	logger *slog.Logger
}

func NewAuditService(logger *slog.Logger) *AuditService {
	return &AuditService{logger: logger}
}

func (s *AuditService) Write(record AuditRecord) {
	at := record.At
	if at.IsZero() {
		at = time.Now().UTC()
	}

	s.logger.Info(
		"audit_log",
		"tenant_id", record.TenantID,
		"actor_id", record.ActorID,
		"action", record.Action,
		"resource", record.Resource,
		"result", record.Result,
		"at", at.Format(time.RFC3339Nano),
	)
}
