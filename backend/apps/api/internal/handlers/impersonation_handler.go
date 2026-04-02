package handlers

import (
	"errors"
	"net/http"

	"api/internal/database"
	"api/internal/models"
	"api/internal/services"

	"github.com/gin-gonic/gin"
)

type ImpersonationHandler struct {
	service *services.ImpersonationService
}

func NewImpersonationHandler(service *services.ImpersonationService) *ImpersonationHandler {
	return &ImpersonationHandler{service: service}
}

func (h *ImpersonationHandler) Start(c *gin.Context) {
	var req models.ImpersonationStartInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	session, err := h.service.Start(c.Request.Context(), userIDFromClaims(c), req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"sessionId": session.SessionID,
		"tenantId":  session.TenantID,
		"startedAt": session.StartedAt,
		"mode":      session.Mode,
	})
}

func (h *ImpersonationHandler) End(c *gin.Context) {
	err := h.service.End(c.Request.Context(), userIDFromClaims(c), c.Param("sessionId"))
	if err != nil {
		switch {
		case errors.Is(err, database.ErrImpersonationSessionNotFound):
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrImpersonationClosedSession):
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ImpersonationHandler) ListAudit(c *gin.Context) {
	list, err := h.service.GetAudit(c.Request.Context(), c.Param("tenantId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}
