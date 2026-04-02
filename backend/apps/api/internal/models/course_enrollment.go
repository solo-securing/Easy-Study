package models

import "time"

type EnrollmentSourceType string

const (
	EnrollmentSourceDirect EnrollmentSourceType = "direct"
	EnrollmentSourceGroup  EnrollmentSourceType = "group"
	EnrollmentSourceRole   EnrollmentSourceType = "role"
)

type EnrollmentAssignmentMode string

const (
	EnrollmentAssignmentUserIDs  EnrollmentAssignmentMode = "user_ids"
	EnrollmentAssignmentGroupIDs EnrollmentAssignmentMode = "group_ids"
	EnrollmentAssignmentRoles    EnrollmentAssignmentMode = "roles"
)

type CourseEnrollment struct {
	UserID     string               `json:"userId"`
	SourceType EnrollmentSourceType `json:"sourceType"`
	AssignedAt time.Time            `json:"assignedAt"`
}

type EnrollmentList struct {
	CourseID string             `json:"courseId"`
	Items    []CourseEnrollment `json:"items"`
}

type UserCourseEnrollment struct {
	CourseID string
}

type EnrollmentAssignInput struct {
	TenantID       string
	CourseID       string
	AssignmentMode EnrollmentAssignmentMode `json:"assignmentMode"`
	UserIDs        []string                 `json:"userIds,omitempty"`
	GroupIDs       []string                 `json:"groupIds,omitempty"`
	Roles          []UserRole               `json:"roles,omitempty"`
}

type CourseProgressSnapshot struct {
	CourseID           string    `json:"courseId"`
	UserID             string    `json:"userId"`
	CompletionPercent  float64   `json:"completionPercent"`
	CompletedUnitCount int       `json:"completedUnitCount"`
	UpdatedAt          time.Time `json:"updatedAt"`
}
