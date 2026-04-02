package services

import (
	"context"
	"errors"
	"time"

	"api/internal/database"
	"api/internal/models"
)

type EnrollmentService struct {
	enrollments database.EnrollmentRepository
	users       database.UserRepository
	groups      database.GroupRepository
	courses     database.CourseRepository
	audit       *AuditService
}

func NewEnrollmentService(
	enrollments database.EnrollmentRepository,
	users database.UserRepository,
	groups database.GroupRepository,
	courses database.CourseRepository,
	audit *AuditService,
) *EnrollmentService {
	return &EnrollmentService{
		enrollments: enrollments,
		users:       users,
		groups:      groups,
		courses:     courses,
		audit:       audit,
	}
}

func (s *EnrollmentService) List(ctx context.Context, tenantID, courseID string) (models.EnrollmentList, error) {
	if _, err := s.courses.Get(ctx, tenantID, courseID); err != nil {
		return models.EnrollmentList{}, err
	}
	return s.enrollments.List(ctx, tenantID, courseID)
}

func (s *EnrollmentService) Assign(ctx context.Context, actorID string, input models.EnrollmentAssignInput) (models.EnrollmentList, error) {
	if _, err := s.courses.Get(ctx, input.TenantID, input.CourseID); err != nil {
		return models.EnrollmentList{}, err
	}

	rows := make([]models.CourseEnrollment, 0)
	switch input.AssignmentMode {
	case models.EnrollmentAssignmentUserIDs:
		if len(input.UserIDs) == 0 {
			return models.EnrollmentList{}, errors.New("userIds is required for user_ids mode")
		}
		for _, userID := range input.UserIDs {
			if _, err := s.users.Get(ctx, input.TenantID, userID); err != nil {
				return models.EnrollmentList{}, err
			}
			rows = append(rows, models.CourseEnrollment{
				UserID:     userID,
				SourceType: models.EnrollmentSourceDirect,
				AssignedAt: time.Now().UTC(),
			})
		}
	case models.EnrollmentAssignmentGroupIDs:
		if len(input.GroupIDs) == 0 {
			return models.EnrollmentList{}, errors.New("groupIds is required for group_ids mode")
		}
		seen := map[string]struct{}{}
		for _, groupID := range input.GroupIDs {
			memberIDs, err := s.groups.MemberIDs(ctx, input.TenantID, groupID)
			if err != nil {
				return models.EnrollmentList{}, err
			}
			for _, userID := range memberIDs {
				if _, ok := seen[userID]; ok {
					continue
				}
				if _, err := s.users.Get(ctx, input.TenantID, userID); err != nil {
					return models.EnrollmentList{}, err
				}
				seen[userID] = struct{}{}
				rows = append(rows, models.CourseEnrollment{
					UserID:     userID,
					SourceType: models.EnrollmentSourceGroup,
					AssignedAt: time.Now().UTC(),
				})
			}
		}
	case models.EnrollmentAssignmentRoles:
		if len(input.Roles) == 0 {
			return models.EnrollmentList{}, errors.New("roles is required for roles mode")
		}
		for _, role := range input.Roles {
			usersList, err := s.users.List(ctx, models.UserListFilter{
				TenantID: input.TenantID,
				Role:     role,
				Page:     1,
				Size:     1000,
			})
			if err != nil {
				return models.EnrollmentList{}, err
			}
			for _, user := range usersList.Items {
				rows = append(rows, models.CourseEnrollment{
					UserID:     user.ID,
					SourceType: models.EnrollmentSourceRole,
					AssignedAt: time.Now().UTC(),
				})
			}
		}
	default:
		return models.EnrollmentList{}, errors.New("invalid assignmentMode")
	}

	out, err := s.enrollments.Replace(ctx, input, dedupeEnrollments(rows))
	if err != nil {
		return models.EnrollmentList{}, err
	}

	if s.audit != nil {
		s.audit.Write(AuditRecord{
			TenantID: input.TenantID,
			ActorID:  actorID,
			Action:   "course.assign_enrollments",
			Resource: "course",
			Result:   "success",
		})
	}
	return out, nil
}

func dedupeEnrollments(rows []models.CourseEnrollment) []models.CourseEnrollment {
	seen := map[string]struct{}{}
	out := make([]models.CourseEnrollment, 0, len(rows))
	for _, row := range rows {
		if _, ok := seen[row.UserID]; ok {
			continue
		}
		seen[row.UserID] = struct{}{}
		out = append(out, row)
	}
	return out
}
