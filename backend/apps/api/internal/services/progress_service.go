package services

import (
	"context"

	"api/internal/database"
	"api/internal/models"
)

type ProgressService struct {
	enrollments database.EnrollmentRepository
}

func NewProgressService(enrollments database.EnrollmentRepository) *ProgressService {
	return &ProgressService{enrollments: enrollments}
}

func (s *ProgressService) GetMyCourses(ctx context.Context, tenantID, userID string) (models.CourseList, error) {
	links, err := s.enrollments.ListByUser(ctx, tenantID, userID)
	if err != nil {
		return models.CourseList{}, err
	}

	items := make([]models.CourseSummary, 0, len(links))
	for _, link := range links {
		items = append(items, models.CourseSummary{
			ID:     link.CourseID,
			Title:  "Course " + link.CourseID,
			Status: models.CourseStatusPublished,
		})
	}
	return models.CourseList{
		Items: items,
		Meta: models.PageMeta{
			Page:  1,
			Size:  len(items),
			Total: len(items),
		},
	}, nil
}

func (s *ProgressService) GetMyProgress(ctx context.Context, tenantID, courseID, userID string) (models.CourseProgressSnapshot, error) {
	_, err := s.enrollments.List(ctx, tenantID, courseID)
	if err != nil {
		return models.CourseProgressSnapshot{}, err
	}
	return models.CourseProgressSnapshot{
		CourseID:           courseID,
		UserID:             userID,
		CompletionPercent:  0,
		CompletedUnitCount: 0,
	}, nil
}
