package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"api/internal/database"
	"api/internal/models"
	"api/internal/services"

	"github.com/gin-gonic/gin"
)

type TenantHandler struct {
	service *services.TenantService
}

func NewTenantHandler(service *services.TenantService) *TenantHandler {
	return &TenantHandler{service: service}
}

func (h *TenantHandler) ListTenants(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	status := models.TenantStatus(c.Query("status"))

	result, err := h.service.List(c.Request.Context(), models.TenantListFilter{
		Status: status,
		Page:   page,
		Size:   size,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *TenantHandler) CreateTenant(c *gin.Context) {
	var input models.TenantCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	tenant, err := h.service.Create(c.Request.Context(), actorID(c), input)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, services.ErrReservedSubdomain) {
			status = http.StatusConflict
		}
		c.AbortWithStatusJSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tenant)
}

func (h *TenantHandler) GetTenantByID(c *gin.Context) {
	tenant, err := h.service.GetByID(c.Request.Context(), c.Param("tenantId"))
	if err != nil {
		if errors.Is(err, database.ErrTenantNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

func (h *TenantHandler) UpdateTenantStatus(c *gin.Context) {
	var patch models.TenantStatusPatch
	if err := c.ShouldBindJSON(&patch); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	tenant, err := h.service.UpdateStatus(c.Request.Context(), actorID(c), c.Param("tenantId"), patch)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrTenantNotFound):
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrInvalidTenantStatusTransition):
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, tenant)
}

func actorID(c *gin.Context) string {
	requestID := c.GetHeader("X-Request-ID")
	if requestID == "" {
		return "system"
	}

	return requestID
}
