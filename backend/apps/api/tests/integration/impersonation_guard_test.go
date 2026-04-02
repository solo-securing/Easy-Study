package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"api/internal/server"
)

func TestImpersonationGuardBlocksDestructiveActions(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	t.Setenv("JWT_SECRET", "test-secret")

	srv := server.NewServer()
	handler := srv.Handler
	token := signedTokenUS2(t, "test-secret", map[string]any{"role": "super_admin", "sub": "super-admin-001"})

	startBody, _ := json.Marshal(map[string]any{
		"tenantId": "tenant-001",
		"reason":   "support",
	})
	startReq := httptest.NewRequest(http.MethodPost, "/api/v1/admin/impersonations", bytes.NewBuffer(startBody))
	startReq.Header.Set("Content-Type", "application/json")
	startReq.Header.Set("Authorization", "Bearer "+token)
	startReq.Header.Set("X-Tenant-ID", "platform")
	startRes := httptest.NewRecorder()
	handler.ServeHTTP(startRes, startReq)
	if startRes.Code != http.StatusCreated {
		t.Fatalf("expected 201 start impersonation, got %d body=%s", startRes.Code, startRes.Body.String())
	}

	var started map[string]any
	_ = json.Unmarshal(startRes.Body.Bytes(), &started)
	sessionID, _ := started["sessionId"].(string)
	if sessionID == "" {
		t.Fatalf("expected session id")
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/impersonations/"+sessionID, nil)
	deleteReq.Header.Set("Authorization", "Bearer "+token)
	deleteReq.Header.Set("X-Tenant-ID", "platform")
	deleteReq.Header.Set("X-Impersonation-Session-ID", sessionID)
	deleteRes := httptest.NewRecorder()
	handler.ServeHTTP(deleteRes, deleteReq)
	if deleteRes.Code != http.StatusForbidden {
		t.Fatalf("expected 403 blocked destructive action, got %d body=%s", deleteRes.Code, deleteRes.Body.String())
	}
}
