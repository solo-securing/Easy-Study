package database

import (
	"context"
	"errors"
	"sync"
	"time"

	"api/internal/models"
	"github.com/google/uuid"
)

var (
	ErrImpersonationSessionNotFound = errors.New("impersonation session not found")
)

type ImpersonationRepository interface {
	Start(ctx context.Context, session models.ImpersonationSession) (models.ImpersonationSession, error)
	End(ctx context.Context, sessionID string) error
	TrackAction(ctx context.Context, sessionID, action string) error
	Get(ctx context.Context, sessionID string) (models.ImpersonationSession, error)
	ListByTenant(ctx context.Context, tenantID string) (models.ImpersonationAuditList, error)
}

type InMemoryImpersonationRepository struct {
	mu       sync.RWMutex
	sessions map[string]models.ImpersonationSession
}

func NewInMemoryImpersonationRepository() *InMemoryImpersonationRepository {
	return &InMemoryImpersonationRepository{
		sessions: map[string]models.ImpersonationSession{},
	}
}

func (r *InMemoryImpersonationRepository) Start(_ context.Context, session models.ImpersonationSession) (models.ImpersonationSession, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if session.SessionID == "" {
		session.SessionID = uuid.NewString()
	}
	if session.StartedAt.IsZero() {
		session.StartedAt = time.Now().UTC()
	}
	r.sessions[session.SessionID] = session
	return session, nil
}

func (r *InMemoryImpersonationRepository) End(_ context.Context, sessionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	session, ok := r.sessions[sessionID]
	if !ok {
		return ErrImpersonationSessionNotFound
	}
	now := time.Now().UTC()
	session.EndedAt = &now
	r.sessions[sessionID] = session
	return nil
}

func (r *InMemoryImpersonationRepository) TrackAction(_ context.Context, sessionID, action string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	session, ok := r.sessions[sessionID]
	if !ok {
		return ErrImpersonationSessionNotFound
	}
	session.Actions = append(session.Actions, action)
	r.sessions[sessionID] = session
	return nil
}

func (r *InMemoryImpersonationRepository) Get(_ context.Context, sessionID string) (models.ImpersonationSession, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	session, ok := r.sessions[sessionID]
	if !ok {
		return models.ImpersonationSession{}, ErrImpersonationSessionNotFound
	}
	return session, nil
}

func (r *InMemoryImpersonationRepository) ListByTenant(_ context.Context, tenantID string) (models.ImpersonationAuditList, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]models.ImpersonationSession, 0)
	for _, session := range r.sessions {
		if session.TenantID == tenantID {
			items = append(items, session)
		}
	}
	return models.ImpersonationAuditList{Items: items}, nil
}
