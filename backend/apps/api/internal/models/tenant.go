package models

import "time"

type TenantStatus string

const (
	TenantStatusActive      TenantStatus = "active"
	TenantStatusSuspended   TenantStatus = "suspended"
	TenantStatusDeactivated TenantStatus = "deactivated"
)

type Tenant struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	Subdomain     string       `json:"subdomain"`
	Status        TenantStatus `json:"status"`
	PlanCode      string       `json:"planCode"`
	CreatedAt     time.Time    `json:"createdAt"`
	DefaultLocale string       `json:"defaultLocale"`
	Timezone      string       `json:"timezone"`
}

type TenantCreateInput struct {
	Name               string `json:"name"`
	Subdomain          string `json:"subdomain"`
	OwnerEmail         string `json:"ownerEmail"`
	PlanCode           string `json:"planCode"`
	UserLimitOverride  *int   `json:"userLimitOverride,omitempty"`
	CourseLimitOverride *int  `json:"courseLimitOverride,omitempty"`
}

type TenantStatusPatch struct {
	Status TenantStatus `json:"status"`
	Reason string       `json:"reason,omitempty"`
}

type TenantListFilter struct {
	Status TenantStatus
	Page   int
	Size   int
}

type PageMeta struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

type TenantList struct {
	Items []Tenant  `json:"items"`
	Meta  PageMeta `json:"meta"`
}
