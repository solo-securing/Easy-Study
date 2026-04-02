package database

import (
	"context"
	"errors"
	"sync"

	"api/internal/models"
	"github.com/google/uuid"
)

var (
	ErrGroupNotFound   = errors.New("group not found")
	ErrGroupNameExists = errors.New("group name already exists")
)

type GroupRepository interface {
	List(ctx context.Context, tenantID string) (models.GroupList, error)
	Create(ctx context.Context, input models.GroupCreateInput) (models.Group, error)
	Patch(ctx context.Context, tenantID, groupID string, patch models.GroupPatchInput) (models.Group, error)
	ReplaceMembers(ctx context.Context, input models.GroupMembersUpsert) (models.Group, error)
	MemberIDs(ctx context.Context, tenantID, groupID string) ([]string, error)
}

type inMemoryGroup struct {
	models.Group
	memberIDs map[string]struct{}
}

type InMemoryGroupRepository struct {
	mu     sync.RWMutex
	groups map[string]inMemoryGroup
}

func NewInMemoryGroupRepository() *InMemoryGroupRepository {
	return &InMemoryGroupRepository{groups: map[string]inMemoryGroup{}}
}

func (r *InMemoryGroupRepository) List(_ context.Context, tenantID string) (models.GroupList, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]models.Group, 0)
	for _, group := range r.groups {
		if group.TenantID == tenantID {
			items = append(items, group.Group)
		}
	}
	return models.GroupList{Items: items}, nil
}

func (r *InMemoryGroupRepository) Create(_ context.Context, input models.GroupCreateInput) (models.Group, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existing := range r.groups {
		if existing.TenantID == input.TenantID && existing.Name == input.Name {
			return models.Group{}, ErrGroupNameExists
		}
	}

	group := models.Group{
		ID:          uuid.NewString(),
		TenantID:    input.TenantID,
		Name:        input.Name,
		Description: input.Description,
		MemberCount: 0,
	}
	r.groups[group.ID] = inMemoryGroup{
		Group:     group,
		memberIDs: map[string]struct{}{},
	}
	return group, nil
}

func (r *InMemoryGroupRepository) Patch(_ context.Context, tenantID, groupID string, patch models.GroupPatchInput) (models.Group, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	group, ok := r.groups[groupID]
	if !ok || group.TenantID != tenantID {
		return models.Group{}, ErrGroupNotFound
	}
	if patch.Name != nil {
		group.Name = *patch.Name
	}
	if patch.Description != nil {
		group.Description = *patch.Description
	}
	r.groups[groupID] = group
	return group.Group, nil
}

func (r *InMemoryGroupRepository) ReplaceMembers(_ context.Context, input models.GroupMembersUpsert) (models.Group, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	group, ok := r.groups[input.GroupID]
	if !ok || group.TenantID != input.TenantID {
		return models.Group{}, ErrGroupNotFound
	}

	group.memberIDs = make(map[string]struct{}, len(input.UserIDs))
	for _, userID := range input.UserIDs {
		group.memberIDs[userID] = struct{}{}
	}
	group.MemberCount = len(group.memberIDs)
	r.groups[input.GroupID] = group
	return group.Group, nil
}

func (r *InMemoryGroupRepository) MemberIDs(_ context.Context, tenantID, groupID string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	group, ok := r.groups[groupID]
	if !ok || group.TenantID != tenantID {
		return nil, ErrGroupNotFound
	}

	members := make([]string, 0, len(group.memberIDs))
	for userID := range group.memberIDs {
		members = append(members, userID)
	}
	return members, nil
}
