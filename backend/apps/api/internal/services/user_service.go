package services

import (
	"context"
	"errors"
	"strings"

	"api/internal/database"
	"api/internal/models"
)

type UserService struct {
	users      database.UserRepository
	activation *ActivationService
	audit      *AuditService
}

func NewUserService(users database.UserRepository, activation *ActivationService, audit *AuditService) *UserService {
	return &UserService{
		users:      users,
		activation: activation,
		audit:      audit,
	}
}

func (s *UserService) List(ctx context.Context, filter models.UserListFilter) (models.UserList, error) {
	return s.users.List(ctx, filter)
}

func (s *UserService) Create(ctx context.Context, actorID string, input models.UserCreateInput) (models.User, *models.ActivationToken, error) {
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.FullName = strings.TrimSpace(input.FullName)
	if input.Email == "" || input.FullName == "" || input.Role == "" {
		return models.User{}, nil, errors.New("missing required user fields")
	}

	user, err := s.users.Create(ctx, input)
	if err != nil {
		return models.User{}, nil, err
	}

	var tokenModel *models.ActivationToken
	if input.SendInvitation {
		_, issued, err := s.activation.Issue(ctx, input.TenantID, user.ID)
		if err != nil {
			return models.User{}, nil, err
		}
		tokenModel = &issued
	}

	if s.audit != nil {
		s.audit.Write(AuditRecord{
			TenantID: input.TenantID,
			ActorID:  actorID,
			Action:   "user.create",
			Resource: "user",
			Result:   "success",
		})
	}

	return user, tokenModel, nil
}

func (s *UserService) Get(ctx context.Context, tenantID, userID string) (models.User, error) {
	return s.users.Get(ctx, tenantID, userID)
}

func (s *UserService) Patch(ctx context.Context, actorID, tenantID, userID string, patch models.UserPatchInput) (models.User, error) {
	user, err := s.users.Update(ctx, tenantID, userID, patch)
	if err != nil {
		return models.User{}, err
	}

	if s.audit != nil {
		s.audit.Write(AuditRecord{
			TenantID: tenantID,
			ActorID:  actorID,
			Action:   "user.patch",
			Resource: "user",
			Result:   "success",
		})
	}

	return user, nil
}

func (s *UserService) ResendInvitation(ctx context.Context, tenantID, userID string) error {
	_, err := s.users.Get(ctx, tenantID, userID)
	if err != nil {
		return err
	}
	_, _, err = s.activation.Issue(ctx, tenantID, userID)
	return err
}
