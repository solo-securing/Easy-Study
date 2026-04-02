package handlers

import (
	"net/http"

	"api/internal/middleware"
	"api/internal/services"

	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	progress *services.ProgressService
}

func NewStudentHandler(progress *services.ProgressService) *StudentHandler {
	return &StudentHandler{progress: progress}
}

func (h *StudentHandler) ListMyCourses(c *gin.Context) {
	tenantID, ok := middleware.TenantID(c)
	if !ok || tenantID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing tenant context"})
		return
	}
	list, err := h.progress.GetMyCourses(c.Request.Context(), tenantID, userIDFromClaims(c))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *StudentHandler) GetMyCourseProgress(c *gin.Context) {
	snapshot, err := h.progress.GetMyProgress(c.Request.Context(), c.Param("tenantId"), c.Param("courseId"), userIDFromClaims(c))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"courseId":           snapshot.CourseID,
		"completionPercent":  snapshot.CompletionPercent,
		"completedUnitCount": snapshot.CompletedUnitCount,
		"officialQuizScore":  0,
	})
}
