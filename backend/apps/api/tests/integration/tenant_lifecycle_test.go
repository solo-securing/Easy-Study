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

func TestTenantLifecycleFlow(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	t.Setenv("JWT_SECRET", "test-secret")

	httpServer := server.NewServer()
	handler := httpServer.Handler

	token := signedToken(t, "test-secret", map[string]any{"role": "super_admin"})

	createPayload := map[string]any{
		"name":       "Tenant One",
		"subdomain":  "tenant-one",
		"ownerEmail": "owner@example.com",
		"planCode":   "pro",
	}
	createBody, _ := json.Marshal(createPayload)

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/super-admin/tenants", bytes.NewBuffer(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("Authorization", "Bearer "+token)
	createReq.Header.Set("X-Tenant-ID", "platform")
	createRes := httptest.NewRecorder()
	handler.ServeHTTP(createRes, createReq)
	if createRes.Code != http.StatusCreated {
		t.Fatalf("expected create status 201, got %d body=%s", createRes.Code, createRes.Body.String())
	}

	var created map[string]any
	if err := json.Unmarshal(createRes.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal create response: %v", err)
	}
	tenantID, _ := created["id"].(string)
	if tenantID == "" {
		t.Fatalf("expected tenant id in create response")
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/super-admin/tenants", nil)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listReq.Header.Set("X-Tenant-ID", "platform")
	listRes := httptest.NewRecorder()
	handler.ServeHTTP(listRes, listReq)
	if listRes.Code != http.StatusOK {
		t.Fatalf("expected list status 200, got %d body=%s", listRes.Code, listRes.Body.String())
	}

	statusPatch := map[string]any{"status": "suspended", "reason": "manual review"}
	statusBody, _ := json.Marshal(statusPatch)
	statusReq := httptest.NewRequest(http.MethodPatch, "/api/v1/super-admin/tenants/"+tenantID+"/status", bytes.NewBuffer(statusBody))
	statusReq.Header.Set("Content-Type", "application/json")
	statusReq.Header.Set("Authorization", "Bearer "+token)
	statusReq.Header.Set("X-Tenant-ID", "platform")
	statusRes := httptest.NewRecorder()
	handler.ServeHTTP(statusRes, statusReq)
	if statusRes.Code != http.StatusOK {
		t.Fatalf("expected status patch 200, got %d body=%s", statusRes.Code, statusRes.Body.String())
	}
}

func signedToken(t *testing.T, secret string, claims map[string]any) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	out, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return out
}
