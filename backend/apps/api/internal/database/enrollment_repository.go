package database

import (
	"context"
	"sync"
	"time"

	"api/internal/models"
)

type EnrollmentRepository interface {
	List(ctx context.Context, tenantID, courseID string) (models.EnrollmentList, error)
	Replace(ctx context.Context, input models.EnrollmentAssignInput, rows []models.CourseEnrollment) (models.EnrollmentList, error)
	ListByUser(ctx context.Context, tenantID, userID string) ([]models.UserCourseEnrollment, error)
}

type InMemoryEnrollmentRepository struct {
	mu    sync.RWMutex
	items map[string][]models.CourseEnrollment
}

func NewInMemoryEnrollmentRepository() *InMemoryEnrollmentRepository {
	return &InMemoryEnrollmentRepository{
		items: map[string][]models.CourseEnrollment{},
	}
}

func enrollmentKey(tenantID, courseID string) string {
	return tenantID + ":" + courseID
}

func (r *InMemoryEnrollmentRepository) List(_ context.Context, tenantID, courseID string) (models.EnrollmentList, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := enrollmentKey(tenantID, courseID)
	rows := r.items[key]
	out := make([]models.CourseEnrollment, len(rows))
	copy(out, rows)
	return models.EnrollmentList{
		CourseID: courseID,
		Items:    out,
	}, nil
}

func (r *InMemoryEnrollmentRepository) Replace(_ context.Context, input models.EnrollmentAssignInput, rows []models.CourseEnrollment) (models.EnrollmentList, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := enrollmentKey(input.TenantID, input.CourseID)
	normalized := make([]models.CourseEnrollment, 0, len(rows))
	for _, row := range rows {
		if row.AssignedAt.IsZero() {
			row.AssignedAt = time.Now().UTC()
		}
		normalized = append(normalized, row)
	}
	r.items[key] = normalized
	return models.EnrollmentList{
		CourseID: input.CourseID,
		Items:    normalized,
	}, nil
}

func (r *InMemoryEnrollmentRepository) ListByUser(_ context.Context, tenantID, userID string) ([]models.UserCourseEnrollment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]models.UserCourseEnrollment, 0)
	for key, rows := range r.items {
		prefix := tenantID + ":"
		if len(key) < len(prefix) || key[:len(prefix)] != prefix {
			continue
		}
		courseID := key[len(prefix):]
		for _, row := range rows {
			if row.UserID == userID {
				out = append(out, models.UserCourseEnrollment{CourseID: courseID})
				break
			}
		}
	}
	return out, nil
}
