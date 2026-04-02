package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"api/internal/server"
)

func TestQuizAutosaveRecoveryAndTimeoutAutoSubmit(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	t.Setenv("JWT_SECRET", "test-secret")

	srv := server.NewServer()
	handler := srv.Handler
	token := signedTokenUS2(t, "test-secret", map[string]any{
		"role": "student",
		"sub":  "student-quiz-001",
	})

	startReq := httptest.NewRequest(http.MethodPost, "/api/v1/tenants/tenant-001/quizzes/quiz-001/attempts", bytes.NewBufferString("{}"))
	startReq.Header.Set("Content-Type", "application/json")
	startReq.Header.Set("Authorization", "Bearer "+token)
	startReq.Header.Set("X-Tenant-ID", "tenant-001")
	startRes := httptest.NewRecorder()
	handler.ServeHTTP(startRes, startReq)
	if startRes.Code != http.StatusCreated {
		t.Fatalf("expected 201 start attempt, got %d body=%s", startRes.Code, startRes.Body.String())
	}

	var started map[string]any
	_ = json.Unmarshal(startRes.Body.Bytes(), &started)
	attemptID, _ := started["attemptId"].(string)
	if attemptID == "" {
		t.Fatalf("expected attempt id")
	}

	autosavePayload, _ := json.Marshal(map[string]any{
		"answers": []map[string]any{
			{"questionId": "question-001", "textAnswer": "answer-v1"},
		},
		"clientSavedAt": time.Now().UTC().Format(time.RFC3339),
	})
	autosaveReq := httptest.NewRequest(http.MethodPatch, "/api/v1/tenants/tenant-001/quizzes/quiz-001/attempts/"+attemptID+"/autosave", bytes.NewBuffer(autosavePayload))
	autosaveReq.Header.Set("Content-Type", "application/json")
	autosaveReq.Header.Set("Authorization", "Bearer "+token)
	autosaveReq.Header.Set("X-Tenant-ID", "tenant-001")
	autosaveRes := httptest.NewRecorder()
	handler.ServeHTTP(autosaveRes, autosaveReq)
	if autosaveRes.Code != http.StatusOK {
		t.Fatalf("expected 200 autosave, got %d body=%s", autosaveRes.Code, autosaveRes.Body.String())
	}

	recoveryPayload, _ := json.Marshal(map[string]any{
		"resumeAttemptId": attemptID,
	})
	recoveryReq := httptest.NewRequest(http.MethodPost, "/api/v1/tenants/tenant-001/quizzes/quiz-001/attempts", bytes.NewBuffer(recoveryPayload))
	recoveryReq.Header.Set("Content-Type", "application/json")
	recoveryReq.Header.Set("Authorization", "Bearer "+token)
	recoveryReq.Header.Set("X-Tenant-ID", "tenant-001")
	recoveryRes := httptest.NewRecorder()
	handler.ServeHTTP(recoveryRes, recoveryReq)
	if recoveryRes.Code != http.StatusCreated {
		t.Fatalf("expected 201 recovery start, got %d body=%s", recoveryRes.Code, recoveryRes.Body.String())
	}
	var recovered map[string]any
	_ = json.Unmarshal(recoveryRes.Body.Bytes(), &recovered)
	recoveredAttemptID, _ := recovered["attemptId"].(string)
	if recoveredAttemptID != attemptID {
		t.Fatalf("expected resumed same attempt, got %s vs %s", recoveredAttemptID, attemptID)
	}

	detailReq := httptest.NewRequest(http.MethodGet, "/api/v1/tenants/tenant-001/quizzes/quiz-001/attempts/"+attemptID, nil)
	detailReq.Header.Set("Authorization", "Bearer "+token)
	detailReq.Header.Set("X-Tenant-ID", "tenant-001")
	detailRes := httptest.NewRecorder()
	handler.ServeHTTP(detailRes, detailReq)
	if detailRes.Code != http.StatusOK {
		t.Fatalf("expected 200 attempt detail, got %d body=%s", detailRes.Code, detailRes.Body.String())
	}
	var detail map[string]any
	_ = json.Unmarshal(detailRes.Body.Bytes(), &detail)
	answers, _ := detail["answers"].([]any)
	if len(answers) == 0 {
		t.Fatalf("expected recovered answers")
	}

	submitReq := httptest.NewRequest(http.MethodPost, "/api/v1/tenants/tenant-001/quizzes/quiz-001/attempts/"+attemptID+"/submit", nil)
	submitReq.Header.Set("Authorization", "Bearer "+token)
	submitReq.Header.Set("X-Tenant-ID", "tenant-001")
	submitRes := httptest.NewRecorder()
	handler.ServeHTTP(submitRes, submitReq)
	if submitRes.Code != http.StatusOK {
		t.Fatalf("expected 200 submit, got %d body=%s", submitRes.Code, submitRes.Body.String())
	}
}

func TestQuizTimeoutAutoSubmitFlow(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	t.Setenv("JWT_SECRET", "test-secret")

	srv := server.NewServer()
	handler := srv.Handler
	token := signedTokenUS2(t, "test-secret", map[string]any{
		"role": "student",
		"sub":  "student-timeout-001",
	})

	startReq := httptest.NewRequest(http.MethodPost, "/api/v1/tenants/tenant-001/quizzes/quiz-timeout/attempts", bytes.NewBufferString("{}"))
	startReq.Header.Set("Content-Type", "application/json")
	startReq.Header.Set("Authorization", "Bearer "+token)
	startReq.Header.Set("X-Tenant-ID", "tenant-001")
	startRes := httptest.NewRecorder()
	handler.ServeHTTP(startRes, startReq)
	if startRes.Code != http.StatusCreated {
		t.Fatalf("expected 201 start timeout attempt, got %d body=%s", startRes.Code, startRes.Body.String())
	}
	var started map[string]any
	_ = json.Unmarshal(startRes.Body.Bytes(), &started)
	attemptID, _ := started["attemptId"].(string)
	if attemptID == "" {
		t.Fatalf("expected timeout attempt id")
	}

	time.Sleep(1200 * time.Millisecond)

	detailReq := httptest.NewRequest(http.MethodGet, "/api/v1/tenants/tenant-001/quizzes/quiz-timeout/attempts/"+attemptID, nil)
	detailReq.Header.Set("Authorization", "Bearer "+token)
	detailReq.Header.Set("X-Tenant-ID", "tenant-001")
	detailRes := httptest.NewRecorder()
	handler.ServeHTTP(detailRes, detailReq)
	if detailRes.Code != http.StatusOK {
		t.Fatalf("expected 200 timeout attempt detail, got %d body=%s", detailRes.Code, detailRes.Body.String())
	}
	var detail map[string]any
	_ = json.Unmarshal(detailRes.Body.Bytes(), &detail)
	status, _ := detail["status"].(string)
	if status != "auto_submitted" {
		t.Fatalf("expected auto_submitted status, got %s", status)
	}
}
