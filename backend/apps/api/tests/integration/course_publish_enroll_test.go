package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"api/internal/server"
)

func TestCoursePublishValidationAndEnrollmentAssignment(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	t.Setenv("JWT_SECRET", "test-secret")

	srv := server.NewServer()
	handler := srv.Handler
	token := signedTokenUS2(t, "test-secret", map[string]any{"role": "instructor"})

	userPayload, _ := json.Marshal(map[string]any{
		"email":          "student-course@example.com",
		"fullName":       "Course Student",
		"role":           "student",
		"sendInvitation": false,
	})
	userReq := httptest.NewRequest(http.MethodPost, "/api/v1/tenants/tenant-001/users", bytes.NewBuffer(userPayload))
	userReq.Header.Set("Content-Type", "application/json")
	userReq.Header.Set("Authorization", "Bearer "+token)
	userReq.Header.Set("X-Tenant-ID", "tenant-001")
	userRes := httptest.NewRecorder()
	handler.ServeHTTP(userRes, userReq)
	if userRes.Code != http.StatusCreated {
		t.Fatalf("expected create user 201, got %d body=%s", userRes.Code, userRes.Body.String())
	}
	var user map[string]any
	_ = json.Unmarshal(userRes.Body.Bytes(), &user)
	userID, _ := user["id"].(string)

	groupPayload, _ := json.Marshal(map[string]any{
		"name":        "US3 Group",
		"description": "for course enrollment assignment",
	})
	groupReq := httptest.NewRequest(http.MethodPost, "/api/v1/tenants/tenant-001/groups", bytes.NewBuffer(groupPayload))
	groupReq.Header.Set("Content-Type", "application/json")
	groupReq.Header.Set("Authorization", "Bearer "+token)
	groupReq.Header.Set("X-Tenant-ID", "tenant-001")
	groupRes := httptest.NewRecorder()
	handler.ServeHTTP(groupRes, groupReq)
	if groupRes.Code != http.StatusCreated {
		t.Fatalf("expected create group 201, got %d body=%s", groupRes.Code, groupRes.Body.String())
	}
	var group map[string]any
	_ = json.Unmarshal(groupRes.Body.Bytes(), &group)
	groupID, _ := group["id"].(string)

	memberPayload, _ := json.Marshal(map[string]any{
		"userIds": []string{userID},
	})
	memberReq := httptest.NewRequest(http.MethodPut, "/api/v1/tenants/tenant-001/groups/"+groupID+"/members", bytes.NewBuffer(memberPayload))
	memberReq.Header.Set("Content-Type", "application/json")
	memberReq.Header.Set("Authorization", "Bearer "+token)
	memberReq.Header.Set("X-Tenant-ID", "tenant-001")
	memberRes := httptest.NewRecorder()
	handler.ServeHTTP(memberRes, memberReq)
	if memberRes.Code != http.StatusOK {
		t.Fatalf("expected replace members 200, got %d body=%s", memberRes.Code, memberRes.Body.String())
	}

	coursePayload, _ := json.Marshal(map[string]any{
		"title":       "US3 Course",
		"description": "publish validation and enrollments",
	})
	courseReq := httptest.NewRequest(http.MethodPost, "/api/v1/tenants/tenant-001/courses", bytes.NewBuffer(coursePayload))
	courseReq.Header.Set("Content-Type", "application/json")
	courseReq.Header.Set("Authorization", "Bearer "+token)
	courseReq.Header.Set("X-Tenant-ID", "tenant-001")
	courseRes := httptest.NewRecorder()
	handler.ServeHTTP(courseRes, courseReq)
	if courseRes.Code != http.StatusCreated {
		t.Fatalf("expected create course 201, got %d body=%s", courseRes.Code, courseRes.Body.String())
	}
	var course map[string]any
	_ = json.Unmarshal(courseRes.Body.Bytes(), &course)
	courseID, _ := course["id"].(string)

	publishReq := httptest.NewRequest(http.MethodPost, "/api/v1/tenants/tenant-001/courses/"+courseID+"/publish", nil)
	publishReq.Header.Set("Authorization", "Bearer "+token)
	publishReq.Header.Set("X-Tenant-ID", "tenant-001")
	publishRes := httptest.NewRecorder()
	handler.ServeHTTP(publishRes, publishReq)
	if publishRes.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected publish validation 422, got %d body=%s", publishRes.Code, publishRes.Body.String())
	}

	structurePayload, _ := json.Marshal(map[string]any{
		"sections": []map[string]any{
			{
				"title":     "Section 1",
				"sortOrder": 0,
				"subsections": []map[string]any{
					{
						"title":     "Subsection 1.1",
						"sortOrder": 0,
						"units": []map[string]any{
							{
								"title":     "Unit 1",
								"type":      "video",
								"sortOrder": 0,
								"contentRef": "asset://video-001",
							},
						},
					},
				},
			},
		},
	})
	structureReq := httptest.NewRequest(http.MethodPut, "/api/v1/tenants/tenant-001/courses/"+courseID+"/structure", bytes.NewBuffer(structurePayload))
	structureReq.Header.Set("Content-Type", "application/json")
	structureReq.Header.Set("Authorization", "Bearer "+token)
	structureReq.Header.Set("X-Tenant-ID", "tenant-001")
	structureRes := httptest.NewRecorder()
	handler.ServeHTTP(structureRes, structureReq)
	if structureRes.Code != http.StatusOK {
		t.Fatalf("expected put structure 200, got %d body=%s", structureRes.Code, structureRes.Body.String())
	}

	publishReq2 := httptest.NewRequest(http.MethodPost, "/api/v1/tenants/tenant-001/courses/"+courseID+"/publish", nil)
	publishReq2.Header.Set("Authorization", "Bearer "+token)
	publishReq2.Header.Set("X-Tenant-ID", "tenant-001")
	publishRes2 := httptest.NewRecorder()
	handler.ServeHTTP(publishRes2, publishReq2)
	if publishRes2.Code != http.StatusOK {
		t.Fatalf("expected publish 200, got %d body=%s", publishRes2.Code, publishRes2.Body.String())
	}

	assignPayload, _ := json.Marshal(map[string]any{
		"assignmentMode": "group_ids",
		"groupIds":       []string{groupID},
	})
	assignReq := httptest.NewRequest(http.MethodPost, "/api/v1/tenants/tenant-001/courses/"+courseID+"/enrollments", bytes.NewBuffer(assignPayload))
	assignReq.Header.Set("Content-Type", "application/json")
	assignReq.Header.Set("Authorization", "Bearer "+token)
	assignReq.Header.Set("X-Tenant-ID", "tenant-001")
	assignRes := httptest.NewRecorder()
	handler.ServeHTTP(assignRes, assignReq)
	if assignRes.Code != http.StatusOK {
		t.Fatalf("expected assign enrollments 200, got %d body=%s", assignRes.Code, assignRes.Body.String())
	}
	var enrollments map[string]any
	_ = json.Unmarshal(assignRes.Body.Bytes(), &enrollments)
	items, _ := enrollments["items"].([]any)
	if len(items) == 0 {
		t.Fatalf("expected non-empty enrollments after assignment")
	}
}
