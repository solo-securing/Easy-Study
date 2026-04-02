package models

import "time"

type UserRole string

const (
	UserRoleTenantAdmin UserRole = "tenant_admin"
	UserRoleInstructor  UserRole = "instructor"
	UserRoleStudent     UserRole = "student"
)

type UserStatus string

const (
	UserStatusInvited     UserStatus = "invited"
	UserStatusActive      UserStatus = "active"
	UserStatusSuspended   UserStatus = "suspended"
	UserStatusDeactivated UserStatus = "deactivated"
)

type User struct {
	ID        string     `json:"id"`
	TenantID  string     `json:"tenantId"`
	Email     string     `json:"email"`
	FullName  string     `json:"fullName"`
	Role      UserRole   `json:"role"`
	Status    UserStatus `json:"status"`
	CreatedAt time.Time  `json:"createdAt"`
}

type UserCreateInput struct {
	TenantID        string   `json:"tenantId"`
	Email           string   `json:"email"`
	FullName        string   `json:"fullName"`
	Role            UserRole `json:"role"`
	SendInvitation  bool     `json:"sendInvitation"`
}

type UserPatchInput struct {
	FullName *string     `json:"fullName,omitempty"`
	Role     *UserRole   `json:"role,omitempty"`
	Status   *UserStatus `json:"status,omitempty"`
}

type UserListFilter struct {
	TenantID string
	Role     UserRole
	Status   UserStatus
	Page     int
	Size     int
}

type UserList struct {
	Items []User   `json:"items"`
	Meta  PageMeta `json:"meta"`
}

type ActivationToken struct {
	ID           string
	TenantID     string
	UserID       string
	TokenHash    string
	TokenVersion int
	ExpiresAt    time.Time
	ConsumedAt   *time.Time
}

type Group struct {
	ID          string `json:"id"`
	TenantID    string `json:"tenantId"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MemberCount int    `json:"memberCount"`
}

type GroupCreateInput struct {
	TenantID    string `json:"tenantId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GroupPatchInput struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type GroupMembersUpsert struct {
	TenantID string   `json:"tenantId"`
	GroupID  string   `json:"groupId"`
	UserIDs  []string `json:"userIds"`
}

type GroupList struct {
	Items []Group `json:"items"`
}

type CSVImportJob struct {
	JobID        string               `json:"jobId"`
	TenantID     string               `json:"tenantId"`
	Status       string               `json:"status"`
	TotalRows    int                  `json:"totalRows"`
	SuccessCount int                  `json:"successCount"`
	FailedCount  int                  `json:"failedCount"`
	Errors       []CSVImportJobError  `json:"errors,omitempty"`
}

type CSVImportJobError struct {
	Row     int    `json:"row"`
	Message string `json:"message"`
}

type CSVImportJobRequest struct {
	FileAssetID    string   `json:"fileAssetId"`
	DefaultRole    UserRole `json:"defaultRole"`
	SendInvitation bool     `json:"sendInvitation"`
}
