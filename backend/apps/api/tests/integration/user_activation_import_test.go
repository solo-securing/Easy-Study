package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"api/internal/server"

	"github.com/golang-jwt/jwt/v5"
)

func TestUserActivationAndImportFlow(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	t.Setenv("JWT_SECRET", "test-secret")

	srv := server.NewServer()
	handler := srv.Handler
	token := signedTokenUS2(t, "test-secret", map[string]any{"role": "tenant_admin"})

	userReqBody := map[string]any{
		"email":          "student1@example.com",
		"fullName":       "Student One",
		"role":           "student",
		"sendInvitation": true,
	}
	userPayload, _ := json.Marshal(userReqBody)
	userReq := httptest.NewRequest(http.MethodPost, "/api/v1/tenants/tenant-001/users", bytes.NewBuffer(userPayload))
	userReq.Header.Set("Content-Type", "application/json")
	userReq.Header.Set("Authorization", "Bearer "+token)
	userReq.Header.Set("X-Tenant-ID", "tenant-001")
	userRes := httptest.NewRecorder()
	handler.ServeHTTP(userRes, userReq)
	if userRes.Code != http.StatusCreated {
		t.Fatalf("expected 201 create user, got %d body=%s", userRes.Code, userRes.Body.String())
	}

	importReqBody := map[string]any{
		"fileAssetId":    "62d4f95b-6254-4f5a-95be-bb9f43f4e3f4",
		"defaultRole":    "student",
		"sendInvitation": true,
	}
	importPayload, _ := json.Marshal(importReqBody)
	importReq := httptest.NewRequest(http.MethodPost, "/api/v1/tenants/tenant-001/users/import-jobs", bytes.NewBuffer(importPayload))
	importReq.Header.Set("Content-Type", "application/json")
	importReq.Header.Set("Authorization", "Bearer "+token)
	importReq.Header.Set("X-Tenant-ID", "tenant-001")
	importRes := httptest.NewRecorder()
	handler.ServeHTTP(importRes, importReq)
	if importRes.Code != http.StatusAccepted {
		t.Fatalf("expected 202 create import job, got %d body=%s", importRes.Code, importRes.Body.String())
	}

	activationResendReq := map[string]any{
		"tenantId": "tenant-001",
		"userId":   "user-001",
	}
	resendPayload, _ := json.Marshal(activationResendReq)
	resendReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/activation/resend", bytes.NewBuffer(resendPayload))
	resendReq.Header.Set("Content-Type", "application/json")
	resendReq.Header.Set("Authorization", "Bearer "+token)
	resendReq.Header.Set("X-Tenant-ID", "tenant-001")
	resendRes := httptest.NewRecorder()
	handler.ServeHTTP(resendRes, resendReq)
	if resendRes.Code != http.StatusAccepted {
		t.Fatalf("expected 202 activation resend, got %d body=%s", resendRes.Code, resendRes.Body.String())
	}
}

func signedTokenUS2(t *testing.T, secret string, claims map[string]any) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	out, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return out
}
