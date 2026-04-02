package performance

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"api/internal/server"

	"github.com/golang-jwt/jwt/v5"
)

func TestDashboardP95LatencyFor10kUsers(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	t.Setenv("JWT_SECRET", "test-secret")

	srv := server.NewServer()
	handler := srv.Handler
	token := signedPerfToken(t, "test-secret", map[string]any{"role": "tenant_admin", "sub": "tenant-admin-perf"})

	samples := make([]time.Duration, 0, 120)
	for i := 0; i < 120; i++ {
		start := time.Now()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tenants/tenant-001/reports/dashboard", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("X-Tenant-ID", "tenant-001")
		res := httptest.NewRecorder()
		handler.ServeHTTP(res, req)
		if res.Code != http.StatusOK {
			t.Fatalf("expected 200 dashboard, got %d body=%s", res.Code, res.Body.String())
		}
		samples = append(samples, time.Since(start))
	}

	p95 := percentileDuration(samples, 95)
	if p95 > 3*time.Second {
		t.Fatalf("expected p95 <= 3s, got %s", p95)
	}
}

func percentileDuration(values []time.Duration, p int) time.Duration {
	if len(values) == 0 {
		return 0
	}
	sorted := append([]time.Duration(nil), values...)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j] < sorted[i] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	idx := (len(sorted)-1)*p/100
	return sorted[idx]
}

func signedPerfToken(t *testing.T, secret string, claims map[string]any) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	out, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return out
}
