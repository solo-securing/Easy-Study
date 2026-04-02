package handlers

import (
	"net/http"

	"api/internal/services"

	"github.com/gin-gonic/gin"
)

type TenantSettingsHandler struct {
	service *services.TenantSettingsService
}

func NewTenantSettingsHandler(service *services.TenantSettingsService) *TenantSettingsHandler {
	return &TenantSettingsHandler{service: service}
}

func (h *TenantSettingsHandler) GetBranding(c *gin.Context) {
	c.JSON(http.StatusOK, h.service.GetBranding(c.Request.Context(), c.Param("tenantId")))
}

func (h *TenantSettingsHandler) PatchBranding(c *gin.Context) {
	var patch services.TenantBrandingPatch
	if err := c.ShouldBindJSON(&patch); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	c.JSON(http.StatusOK, h.service.PatchBranding(c.Request.Context(), c.Param("tenantId"), patch))
}

func (h *TenantSettingsHandler) GetSettings(c *gin.Context) {
	c.JSON(http.StatusOK, h.service.GetSettings(c.Request.Context(), c.Param("tenantId")))
}

func (h *TenantSettingsHandler) PatchSettings(c *gin.Context) {
	var patch services.TenantSettingsPatch
	if err := c.ShouldBindJSON(&patch); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	c.JSON(http.StatusOK, h.service.PatchSettings(c.Request.Context(), c.Param("tenantId"), patch))
}
