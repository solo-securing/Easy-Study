package models

import "time"

type ImpersonationMode string

const (
	ImpersonationModeReadConfigNonDestructive ImpersonationMode = "read_and_config_non_destructive"
)

type ImpersonationStartInput struct {
	TenantID string `json:"tenantId"`
	Reason   string `json:"reason,omitempty"`
}

type ImpersonationSession struct {
	SessionID    string           `json:"sessionId"`
	TenantID     string           `json:"tenantId"`
	SuperAdminID string           `json:"superAdminId"`
	Mode         ImpersonationMode `json:"mode"`
	Reason       string           `json:"reason,omitempty"`
	StartedAt    time.Time        `json:"startedAt"`
	EndedAt      *time.Time       `json:"endedAt,omitempty"`
	Actions      []string         `json:"actions,omitempty"`
}

type ImpersonationAuditList struct {
	Items []ImpersonationSession `json:"items"`
}
