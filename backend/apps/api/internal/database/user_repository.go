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
	ErrUserNotFound     = errors.New("user not found")
	ErrUserEmailExists  = errors.New("user email already exists")
)

type UserRepository interface {
	List(ctx context.Context, filter models.UserListFilter) (models.UserList, error)
	Create(ctx context.Context, input models.UserCreateInput) (models.User, error)
	Get(ctx context.Context, tenantID, userID string) (models.User, error)
	Update(ctx context.Context, tenantID, userID string, patch models.UserPatchInput) (models.User, error)
}

type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]models.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{users: map[string]models.User{}}
}

func (r *InMemoryUserRepository) List(_ context.Context, filter models.UserListFilter) (models.UserList, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Size < 1 {
		filter.Size = 20
	}

	filtered := make([]models.User, 0)
	for _, user := range r.users {
		if user.TenantID != filter.TenantID {
			continue
		}
		if filter.Role != "" && user.Role != filter.Role {
			continue
		}
		if filter.Status != "" && user.Status != filter.Status {
			continue
		}
		filtered = append(filtered, user)
	}

	start := (filter.Page - 1) * filter.Size
	if start > len(filtered) {
		start = len(filtered)
	}
	end := start + filter.Size
	if end > len(filtered) {
		end = len(filtered)
	}

	return models.UserList{
		Items: filtered[start:end],
		Meta: models.PageMeta{
			Page:  filter.Page,
			Size:  filter.Size,
			Total: len(filtered),
		},
	}, nil
}

func (r *InMemoryUserRepository) Create(_ context.Context, input models.UserCreateInput) (models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existing := range r.users {
		if existing.TenantID == input.TenantID && existing.Email == input.Email {
			return models.User{}, ErrUserEmailExists
		}
	}

	user := models.User{
		ID:        uuid.NewString(),
		TenantID:  input.TenantID,
		Email:     input.Email,
		FullName:  input.FullName,
		Role:      input.Role,
		Status:    models.UserStatusInvited,
		CreatedAt: time.Now().UTC(),
	}

	r.users[user.ID] = user
	return user, nil
}

func (r *InMemoryUserRepository) Get(_ context.Context, tenantID, userID string) (models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[userID]
	if !ok || user.TenantID != tenantID {
		return models.User{}, ErrUserNotFound
	}
	return user, nil
}

func (r *InMemoryUserRepository) Update(_ context.Context, tenantID, userID string, patch models.UserPatchInput) (models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[userID]
	if !ok || user.TenantID != tenantID {
		return models.User{}, ErrUserNotFound
	}

	if patch.FullName != nil {
		user.FullName = *patch.FullName
	}
	if patch.Role != nil {
		user.Role = *patch.Role
	}
	if patch.Status != nil {
		user.Status = *patch.Status
	}
	r.users[userID] = user
	return user, nil
}
