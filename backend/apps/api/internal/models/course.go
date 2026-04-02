package models

import "time"

type CourseStatus string

const (
	CourseStatusDraft     CourseStatus = "draft"
	CourseStatusPublished CourseStatus = "published"
	CourseStatusArchived  CourseStatus = "archived"
)

type UnitType string

const (
	UnitTypeVideo    UnitType = "video"
	UnitTypeDocument UnitType = "document"
	UnitTypeQuiz     UnitType = "quiz"
)

type Course struct {
	ID          string       `json:"id"`
	TenantID    string       `json:"tenantId"`
	Title       string       `json:"title"`
	Description string       `json:"description,omitempty"`
	Status      CourseStatus `json:"status"`
	Sections    []Section    `json:"sections,omitempty"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

type Section struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	SortOrder   int          `json:"sortOrder"`
	Subsections []Subsection `json:"subsections"`
}

type Subsection struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	SortOrder int    `json:"sortOrder"`
	Units     []Unit `json:"units"`
}

type Unit struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Type       UnitType `json:"type"`
	SortOrder  int      `json:"sortOrder"`
	ContentRef string   `json:"contentRef,omitempty"`
	QuizID     string   `json:"quizId,omitempty"`
}

type CourseCreateInput struct {
	TenantID    string `json:"tenantId"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CoursePatchInput struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
}

type CourseListFilter struct {
	TenantID string
	Status   CourseStatus
	Page     int
	Size     int
}

type CourseList struct {
	Items []CourseSummary `json:"items"`
	Meta  PageMeta        `json:"meta"`
}

type CourseSummary struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description,omitempty"`
	Status      CourseStatus `json:"status"`
}

type CourseStructureUpsertInput struct {
	Sections []SectionInput `json:"sections"`
}

type SectionInput struct {
	Title       string            `json:"title"`
	SortOrder   int               `json:"sortOrder"`
	Subsections []SubsectionInput `json:"subsections"`
}

type SubsectionInput struct {
	Title     string      `json:"title"`
	SortOrder int         `json:"sortOrder"`
	Units     []UnitInput `json:"units"`
}

type UnitInput struct {
	Title      string   `json:"title"`
	Type       UnitType `json:"type"`
	SortOrder  int      `json:"sortOrder"`
	ContentRef string   `json:"contentRef,omitempty"`
	QuizID     string   `json:"quizId,omitempty"`
}

type PublishValidationIssue struct {
	Path   string `json:"path"`
	Reason string `json:"reason"`
}
