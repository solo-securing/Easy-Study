package handlers

import (
	"net/http"
	"strconv"

	"api/internal/services"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reports *services.ReportingService
}

func NewReportHandler(reports *services.ReportingService) *ReportHandler {
	return &ReportHandler{reports: reports}
}

func (h *ReportHandler) TenantDashboard(c *gin.Context) {
	data, err := h.reports.GetTenantDashboard(c.Request.Context(), c.Param("tenantId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *ReportHandler) CourseReport(c *gin.Context) {
	data, err := h.reports.GetCourseReport(c.Request.Context(), c.Param("tenantId"), c.Param("courseId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *ReportHandler) AdminUsage(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	data, err := h.reports.GetAdminUsage(c.Request.Context(), page, size)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
