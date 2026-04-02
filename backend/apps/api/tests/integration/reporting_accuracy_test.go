package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"api/internal/server"
)

func TestReportingAccuracyWithSeededData(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	t.Setenv("JWT_SECRET", "test-secret")

	srv := server.NewServer()
	handler := srv.Handler

	tenantAdminToken := signedTokenUS2(t, "test-secret", map[string]any{"role": "tenant_admin", "sub": "tenant-admin-001"})
	dashboardReq := httptest.NewRequest(http.MethodGet, "/api/v1/tenants/tenant-001/reports/dashboard", nil)
	dashboardReq.Header.Set("Authorization", "Bearer "+tenantAdminToken)
	dashboardReq.Header.Set("X-Tenant-ID", "tenant-001")
	dashboardRes := httptest.NewRecorder()
	handler.ServeHTTP(dashboardRes, dashboardReq)
	if dashboardRes.Code != http.StatusOK {
		t.Fatalf("expected 200 dashboard, got %d body=%s", dashboardRes.Code, dashboardRes.Body.String())
	}

	var dashboard map[string]any
	_ = json.Unmarshal(dashboardRes.Body.Bytes(), &dashboard)
	if dashboard["activeUserCount"] == nil || dashboard["activeCourseCount"] == nil || dashboard["completionRate"] == nil {
		t.Fatalf("expected dashboard fields present")
	}

	instructorToken := signedTokenUS2(t, "test-secret", map[string]any{"role": "instructor", "sub": "instructor-001"})
	courseReq := httptest.NewRequest(http.MethodGet, "/api/v1/tenants/tenant-001/reports/courses/course-001", nil)
	courseReq.Header.Set("Authorization", "Bearer "+instructorToken)
	courseReq.Header.Set("X-Tenant-ID", "tenant-001")
	courseRes := httptest.NewRecorder()
	handler.ServeHTTP(courseRes, courseReq)
	if courseRes.Code != http.StatusOK {
		t.Fatalf("expected 200 course report, got %d body=%s", courseRes.Code, courseRes.Body.String())
	}

	var courseReport map[string]any
	_ = json.Unmarshal(courseRes.Body.Bytes(), &courseReport)
	if courseReport["courseId"] != "course-001" {
		t.Fatalf("expected course report for course-001")
	}
	if _, ok := courseReport["items"].([]any); !ok {
		t.Fatalf("expected course report items array")
	}
}
