package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"api/internal/models"
	"github.com/google/uuid"
)

var (
	ErrActivationTokenInvalid = errors.New("activation token is invalid")
	ErrActivationTokenExpired = errors.New("activation token is expired")
)

type ActivationService struct {
	mu      sync.RWMutex
	byUser  map[string]models.ActivationToken
	byHash  map[string]models.ActivationToken
}

func NewActivationService() *ActivationService {
	return &ActivationService{
		byUser: map[string]models.ActivationToken{},
		byHash: map[string]models.ActivationToken{},
	}
}

func (s *ActivationService) Issue(_ context.Context, tenantID, userID string) (string, models.ActivationToken, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	version := 1
	if existing, ok := s.byUser[userID]; ok {
		version = existing.TokenVersion + 1
		delete(s.byHash, existing.TokenHash)
	}

	rawToken := fmt.Sprintf("%s:%s:%d:%d", tenantID, userID, version, time.Now().UnixNano())
	hash := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hash[:])

	model := models.ActivationToken{
		ID:           uuid.NewString(),
		TenantID:     tenantID,
		UserID:       userID,
		TokenHash:    tokenHash,
		TokenVersion: version,
		ExpiresAt:    time.Now().UTC().Add(24 * time.Hour),
	}

	s.byUser[userID] = model
	s.byHash[tokenHash] = model
	return rawToken, model, nil
}

func (s *ActivationService) Verify(_ context.Context, rawToken string) (models.ActivationToken, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	hash := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hash[:])
	token, ok := s.byHash[tokenHash]
	if !ok {
		return models.ActivationToken{}, ErrActivationTokenInvalid
	}
	if time.Now().UTC().After(token.ExpiresAt) {
		return models.ActivationToken{}, ErrActivationTokenExpired
	}
	return token, nil
}

func (s *ActivationService) Consume(_ context.Context, rawToken string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	hash := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hash[:])
	token, ok := s.byHash[tokenHash]
	if !ok {
		return ErrActivationTokenInvalid
	}
	now := time.Now().UTC()
	token.ConsumedAt = &now
	s.byHash[tokenHash] = token
	s.byUser[token.UserID] = token
	return nil
}
