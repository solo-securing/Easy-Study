package database

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"api/internal/models"
	"github.com/google/uuid"
)

var (
	ErrCourseNotFound = errors.New("course not found")
)

type CourseRepository interface {
	List(ctx context.Context, filter models.CourseListFilter) (models.CourseList, error)
	Create(ctx context.Context, input models.CourseCreateInput) (models.Course, error)
	Get(ctx context.Context, tenantID, courseID string) (models.Course, error)
	Patch(ctx context.Context, tenantID, courseID string, patch models.CoursePatchInput) (models.Course, error)
	ReplaceStructure(ctx context.Context, tenantID, courseID string, input models.CourseStructureUpsertInput) (models.Course, error)
	UpdateStatus(ctx context.Context, tenantID, courseID string, status models.CourseStatus) (models.Course, error)
}

type InMemoryCourseRepository struct {
	mu      sync.RWMutex
	courses map[string]models.Course
}

func NewInMemoryCourseRepository() *InMemoryCourseRepository {
	return &InMemoryCourseRepository{
		courses: map[string]models.Course{},
	}
}

func (r *InMemoryCourseRepository) List(_ context.Context, filter models.CourseListFilter) (models.CourseList, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Size < 1 {
		filter.Size = 20
	}

	all := make([]models.CourseSummary, 0)
	for _, course := range r.courses {
		if course.TenantID != filter.TenantID {
			continue
		}
		if filter.Status != "" && course.Status != filter.Status {
			continue
		}
		all = append(all, models.CourseSummary{
			ID:          course.ID,
			Title:       course.Title,
			Description: course.Description,
			Status:      course.Status,
		})
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].Title < all[j].Title
	})

	start := (filter.Page - 1) * filter.Size
	if start > len(all) {
		start = len(all)
	}
	end := start + filter.Size
	if end > len(all) {
		end = len(all)
	}

	return models.CourseList{
		Items: all[start:end],
		Meta: models.PageMeta{
			Page:  filter.Page,
			Size:  filter.Size,
			Total: len(all),
		},
	}, nil
}

func (r *InMemoryCourseRepository) Create(_ context.Context, input models.CourseCreateInput) (models.Course, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	course := models.Course{
		ID:          uuid.NewString(),
		TenantID:    input.TenantID,
		Title:       input.Title,
		Description: input.Description,
		Status:      models.CourseStatusDraft,
		Sections:    []models.Section{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	r.courses[course.ID] = course
	return course, nil
}

func (r *InMemoryCourseRepository) Get(_ context.Context, tenantID, courseID string) (models.Course, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	course, ok := r.courses[courseID]
	if !ok || course.TenantID != tenantID {
		return models.Course{}, ErrCourseNotFound
	}
	return course, nil
}

func (r *InMemoryCourseRepository) Patch(_ context.Context, tenantID, courseID string, patch models.CoursePatchInput) (models.Course, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	course, ok := r.courses[courseID]
	if !ok || course.TenantID != tenantID {
		return models.Course{}, ErrCourseNotFound
	}
	if patch.Title != nil {
		course.Title = *patch.Title
	}
	if patch.Description != nil {
		course.Description = *patch.Description
	}
	course.UpdatedAt = time.Now().UTC()
	r.courses[courseID] = course
	return course, nil
}

func (r *InMemoryCourseRepository) ReplaceStructure(_ context.Context, tenantID, courseID string, input models.CourseStructureUpsertInput) (models.Course, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	course, ok := r.courses[courseID]
	if !ok || course.TenantID != tenantID {
		return models.Course{}, ErrCourseNotFound
	}

	sections := make([]models.Section, 0, len(input.Sections))
	for _, sec := range input.Sections {
		subsections := make([]models.Subsection, 0, len(sec.Subsections))
		for _, sub := range sec.Subsections {
			units := make([]models.Unit, 0, len(sub.Units))
			for _, unit := range sub.Units {
				units = append(units, models.Unit{
					ID:         uuid.NewString(),
					Title:      unit.Title,
					Type:       unit.Type,
					SortOrder:  unit.SortOrder,
					ContentRef: unit.ContentRef,
					QuizID:     unit.QuizID,
				})
			}
			subsections = append(subsections, models.Subsection{
				ID:        uuid.NewString(),
				Title:     sub.Title,
				SortOrder: sub.SortOrder,
				Units:     units,
			})
		}
		sections = append(sections, models.Section{
			ID:          uuid.NewString(),
			Title:       sec.Title,
			SortOrder:   sec.SortOrder,
			Subsections: subsections,
		})
	}

	course.Sections = sections
	course.UpdatedAt = time.Now().UTC()
	r.courses[courseID] = course
	return course, nil
}

func (r *InMemoryCourseRepository) UpdateStatus(_ context.Context, tenantID, courseID string, status models.CourseStatus) (models.Course, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	course, ok := r.courses[courseID]
	if !ok || course.TenantID != tenantID {
		return models.Course{}, ErrCourseNotFound
	}
	course.Status = status
	course.UpdatedAt = time.Now().UTC()
	r.courses[courseID] = course
	return course, nil
}
