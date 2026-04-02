package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"api/internal/server"
)

func TestImpersonationDoesNotAppearInLearnerCourseProjections(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	t.Setenv("JWT_SECRET", "test-secret")

	srv := server.NewServer()
	handler := srv.Handler

	superToken := signedTokenUS2(t, "test-secret", map[string]any{"role": "super_admin", "sub": "super-admin-learner-check"})
	startBody, _ := json.Marshal(map[string]any{
		"tenantId": "tenant-001",
		"reason":   "support",
	})
	startReq := httptest.NewRequest(http.MethodPost, "/api/v1/admin/impersonations", bytes.NewBuffer(startBody))
	startReq.Header.Set("Content-Type", "application/json")
	startReq.Header.Set("Authorization", "Bearer "+superToken)
	startReq.Header.Set("X-Tenant-ID", "platform")
	startRes := httptest.NewRecorder()
	handler.ServeHTTP(startRes, startReq)
	if startRes.Code != http.StatusCreated {
		t.Fatalf("expected 201 start impersonation, got %d body=%s", startRes.Code, startRes.Body.String())
	}

	studentToken := signedTokenUS2(t, "test-secret", map[string]any{"role": "student", "sub": "student-visibility-001"})
	myCoursesReq := httptest.NewRequest(http.MethodGet, "/api/v1/me/courses", nil)
	myCoursesReq.Header.Set("Authorization", "Bearer "+studentToken)
	myCoursesReq.Header.Set("X-Tenant-ID", "tenant-001")
	myCoursesRes := httptest.NewRecorder()
	handler.ServeHTTP(myCoursesRes, myCoursesReq)
	if myCoursesRes.Code != http.StatusOK {
		t.Fatalf("expected 200 my courses, got %d body=%s", myCoursesRes.Code, myCoursesRes.Body.String())
	}

	var body map[string]any
	_ = json.Unmarshal(myCoursesRes.Body.Bytes(), &body)
	items, _ := body["items"].([]any)
	for _, raw := range items {
		item, _ := raw.(map[string]any)
		if title, _ := item["title"].(string); title == "super_admin" {
			t.Fatalf("super admin principal must not appear in learner projections")
		}
	}
}
