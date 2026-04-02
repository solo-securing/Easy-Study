package models

import "time"

type UsageMetricDaily struct {
	TenantID         string    `json:"tenantId"`
	MetricDate       time.Time `json:"metricDate"`
	ActiveUserCount  int       `json:"activeUserCount"`
	ActiveCourseCount int      `json:"activeCourseCount"`
	CompletionRate   float64   `json:"completionRate"`
	LoginCount       int       `json:"loginCount"`
	LearningMinutes  int       `json:"learningMinutes"`
}

type TenantDashboard struct {
	TenantID         string  `json:"tenantId"`
	ActiveUserCount  int     `json:"activeUserCount"`
	ActiveCourseCount int    `json:"activeCourseCount"`
	CompletionRate   float64 `json:"completionRate"`
}

type AdminUsageRow struct {
	TenantID       string `json:"tenantId"`
	TenantName     string `json:"tenantName"`
	UserCount      int    `json:"userCount"`
	CourseCount    int    `json:"courseCount"`
	LoginCount     int    `json:"loginCount"`
	LearningMinutes int   `json:"learningMinutes"`
}

type AdminUsageReport struct {
	Items []AdminUsageRow `json:"items"`
	Meta  PageMeta        `json:"meta"`
}

type CourseReport struct {
	CourseID       string                `json:"courseId"`
	LearnerCount   int                   `json:"learnerCount"`
	CompletionRate float64               `json:"completionRate"`
	Items          []StudentProgressItem `json:"items"`
}

type StudentProgressItem struct {
	UserID             string  `json:"userId"`
	CompletionPercent  float64 `json:"completionPercent"`
	CompletedUnitCount int     `json:"completedUnitCount"`
	OfficialQuizScore  float64 `json:"officialQuizScore,omitempty"`
}
