package services

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"api/internal/database"
	"api/internal/models"
)

var (
	ErrCoursePublishValidation = errors.New("course structure is invalid for publish")
)

type CourseService struct {
	courses database.CourseRepository
	audit   *AuditService
}

func NewCourseService(courses database.CourseRepository, audit *AuditService) *CourseService {
	return &CourseService{
		courses: courses,
		audit:   audit,
	}
}

func (s *CourseService) List(ctx context.Context, filter models.CourseListFilter) (models.CourseList, error) {
	return s.courses.List(ctx, filter)
}

func (s *CourseService) Create(ctx context.Context, actorID string, input models.CourseCreateInput) (models.Course, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)
	if len(input.Title) < 3 {
		return models.Course{}, errors.New("title must be at least 3 characters")
	}

	course, err := s.courses.Create(ctx, input)
	if err != nil {
		return models.Course{}, err
	}

	if s.audit != nil {
		s.audit.Write(AuditRecord{
			TenantID: input.TenantID,
			ActorID:  actorID,
			Action:   "course.create",
			Resource: "course",
			Result:   "success",
		})
	}
	return course, nil
}

func (s *CourseService) Get(ctx context.Context, tenantID, courseID string) (models.Course, error) {
	return s.courses.Get(ctx, tenantID, courseID)
}

func (s *CourseService) Patch(ctx context.Context, actorID, tenantID, courseID string, patch models.CoursePatchInput) (models.Course, error) {
	if patch.Title != nil {
		trimmed := strings.TrimSpace(*patch.Title)
		if len(trimmed) < 3 {
			return models.Course{}, errors.New("title must be at least 3 characters")
		}
		patch.Title = &trimmed
	}
	if patch.Description != nil {
		trimmed := strings.TrimSpace(*patch.Description)
		patch.Description = &trimmed
	}

	course, err := s.courses.Patch(ctx, tenantID, courseID, patch)
	if err != nil {
		return models.Course{}, err
	}

	if s.audit != nil {
		s.audit.Write(AuditRecord{
			TenantID: tenantID,
			ActorID:  actorID,
			Action:   "course.patch",
			Resource: "course",
			Result:   "success",
		})
	}

	return course, nil
}

func (s *CourseService) ReplaceStructure(ctx context.Context, actorID, tenantID, courseID string, input models.CourseStructureUpsertInput) (models.Course, error) {
	course, err := s.courses.ReplaceStructure(ctx, tenantID, courseID, input)
	if err != nil {
		return models.Course{}, err
	}
	if s.audit != nil {
		s.audit.Write(AuditRecord{
			TenantID: tenantID,
			ActorID:  actorID,
			Action:   "course.replace_structure",
			Resource: "course",
			Result:   "success",
		})
	}
	return course, nil
}

func (s *CourseService) Publish(ctx context.Context, actorID, tenantID, courseID string) (models.Course, []models.PublishValidationIssue, error) {
	course, err := s.courses.Get(ctx, tenantID, courseID)
	if err != nil {
		return models.Course{}, nil, err
	}

	issues := validateCourseStructure(course.Sections)
	if len(issues) > 0 {
		return models.Course{}, issues, ErrCoursePublishValidation
	}

	updated, err := s.courses.UpdateStatus(ctx, tenantID, courseID, models.CourseStatusPublished)
	if err != nil {
		return models.Course{}, nil, err
	}

	if s.audit != nil {
		s.audit.Write(AuditRecord{
			TenantID: tenantID,
			ActorID:  actorID,
			Action:   "course.publish",
			Resource: "course",
			Result:   "success",
		})
	}
	return updated, nil, nil
}

func (s *CourseService) Archive(ctx context.Context, actorID, tenantID, courseID string) (models.Course, error) {
	updated, err := s.courses.UpdateStatus(ctx, tenantID, courseID, models.CourseStatusArchived)
	if err != nil {
		return models.Course{}, err
	}
	if s.audit != nil {
		s.audit.Write(AuditRecord{
			TenantID: tenantID,
			ActorID:  actorID,
			Action:   "course.archive",
			Resource: "course",
			Result:   "success",
		})
	}
	return updated, nil
}

func validateCourseStructure(sections []models.Section) []models.PublishValidationIssue {
	issues := make([]models.PublishValidationIssue, 0)
	if len(sections) == 0 {
		issues = append(issues, models.PublishValidationIssue{
			Path:   "/sections",
			Reason: "at least one section is required",
		})
		return issues
	}

	for i, section := range sections {
		secPath := "/sections/" + strconvItoa(i)
		if strings.TrimSpace(section.Title) == "" {
			issues = append(issues, models.PublishValidationIssue{
				Path:   secPath + "/title",
				Reason: "title is required",
			})
		}
		if len(section.Subsections) == 0 {
			issues = append(issues, models.PublishValidationIssue{
				Path:   secPath + "/subsections",
				Reason: "at least one subsection is required",
			})
		}
		for j, subsection := range section.Subsections {
			subPath := secPath + "/subsections/" + strconvItoa(j)
			if strings.TrimSpace(subsection.Title) == "" {
				issues = append(issues, models.PublishValidationIssue{
					Path:   subPath + "/title",
					Reason: "title is required",
				})
			}
			if len(subsection.Units) == 0 {
				issues = append(issues, models.PublishValidationIssue{
					Path:   subPath + "/units",
					Reason: "at least one unit is required",
				})
			}
			for k, unit := range subsection.Units {
				unitPath := subPath + "/units/" + strconvItoa(k)
				if strings.TrimSpace(unit.Title) == "" {
					issues = append(issues, models.PublishValidationIssue{
						Path:   unitPath + "/title",
						Reason: "title is required",
					})
				}
			}
		}
	}
	return issues
}

func strconvItoa(v int) string {
	return strconv.Itoa(v)
}
