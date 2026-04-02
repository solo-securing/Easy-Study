package services

import (
	"context"
	"sync"
)

type TenantBranding struct {
	DisplayName string `json:"displayName"`
	LogoURL     string `json:"logoUrl,omitempty"`
	FaviconURL  string `json:"faviconUrl,omitempty"`
	PrimaryColor string `json:"primaryColor"`
}

type TenantBrandingPatch struct {
	DisplayName *string `json:"displayName,omitempty"`
	LogoURL     *string `json:"logoUrl,omitempty"`
	FaviconURL  *string `json:"faviconUrl,omitempty"`
	PrimaryColor *string `json:"primaryColor,omitempty"`
}

type TenantSettings struct {
	DefaultLocale               string `json:"defaultLocale"`
	Timezone                    string `json:"timezone"`
	SessionTTLMinutes           int    `json:"sessionTtlMinutes"`
	AllowStudentSelfRegistration bool   `json:"allowStudentSelfRegistration"`
}

type TenantSettingsPatch struct {
	DefaultLocale               *string `json:"defaultLocale,omitempty"`
	Timezone                    *string `json:"timezone,omitempty"`
	SessionTTLMinutes           *int    `json:"sessionTtlMinutes,omitempty"`
	AllowStudentSelfRegistration *bool  `json:"allowStudentSelfRegistration,omitempty"`
}

type TenantSettingsService struct {
	mu       sync.RWMutex
	branding map[string]TenantBranding
	settings map[string]TenantSettings
}

func NewTenantSettingsService() *TenantSettingsService {
	return &TenantSettingsService{
		branding: map[string]TenantBranding{},
		settings: map[string]TenantSettings{},
	}
}

func (s *TenantSettingsService) GetBranding(_ context.Context, tenantID string) TenantBranding {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if branding, ok := s.branding[tenantID]; ok {
		return branding
	}
	return TenantBranding{DisplayName: "", PrimaryColor: "#1D4ED8"}
}

func (s *TenantSettingsService) PatchBranding(_ context.Context, tenantID string, patch TenantBrandingPatch) TenantBranding {
	s.mu.Lock()
	defer s.mu.Unlock()
	current := s.branding[tenantID]
	if patch.DisplayName != nil {
		current.DisplayName = *patch.DisplayName
	}
	if patch.LogoURL != nil {
		current.LogoURL = *patch.LogoURL
	}
	if patch.FaviconURL != nil {
		current.FaviconURL = *patch.FaviconURL
	}
	if patch.PrimaryColor != nil {
		current.PrimaryColor = *patch.PrimaryColor
	}
	if current.PrimaryColor == "" {
		current.PrimaryColor = "#1D4ED8"
	}
	s.branding[tenantID] = current
	return current
}

func (s *TenantSettingsService) GetSettings(_ context.Context, tenantID string) TenantSettings {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if settings, ok := s.settings[tenantID]; ok {
		return settings
	}
	return TenantSettings{
		DefaultLocale:               "en",
		Timezone:                    "UTC",
		SessionTTLMinutes:           60,
		AllowStudentSelfRegistration: false,
	}
}

func (s *TenantSettingsService) PatchSettings(_ context.Context, tenantID string, patch TenantSettingsPatch) TenantSettings {
	s.mu.Lock()
	defer s.mu.Unlock()
	current := s.settings[tenantID]
	if current.DefaultLocale == "" {
		current.DefaultLocale = "en"
	}
	if current.Timezone == "" {
		current.Timezone = "UTC"
	}
	if current.SessionTTLMinutes == 0 {
		current.SessionTTLMinutes = 60
	}
	if patch.DefaultLocale != nil {
		current.DefaultLocale = *patch.DefaultLocale
	}
	if patch.Timezone != nil {
		current.Timezone = *patch.Timezone
	}
	if patch.SessionTTLMinutes != nil {
		current.SessionTTLMinutes = *patch.SessionTTLMinutes
	}
	if patch.AllowStudentSelfRegistration != nil {
		current.AllowStudentSelfRegistration = *patch.AllowStudentSelfRegistration
	}
	s.settings[tenantID] = current
	return current
}
