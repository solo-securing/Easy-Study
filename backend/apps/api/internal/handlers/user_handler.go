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

type UserHandler struct {
	users   *services.UserService
	groups  database.GroupRepository
	imports *services.UserImportService
}

func NewUserHandler(users *services.UserService, groups database.GroupRepository, imports *services.UserImportService) *UserHandler {
	return &UserHandler{
		users:   users,
		groups:  groups,
		imports: imports,
	}
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	role := models.UserRole(c.Query("role"))
	status := models.UserStatus(c.Query("status"))

	users, err := h.users.List(c.Request.Context(), models.UserListFilter{
		TenantID: c.Param("tenantId"),
		Role:     role,
		Status:   status,
		Page:     page,
		Size:     size,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.UserCreateInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	req.TenantID = c.Param("tenantId")

	user, _, err := h.users.Create(c.Request.Context(), actorID(c), req)
	if err != nil {
		if errors.Is(err, database.ErrUserEmailExists) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	user, err := h.users.Get(c.Request.Context(), c.Param("tenantId"), c.Param("userId"))
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) PatchUser(c *gin.Context) {
	var patch models.UserPatchInput
	if err := c.ShouldBindJSON(&patch); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	user, err := h.users.Patch(c.Request.Context(), actorID(c), c.Param("tenantId"), c.Param("userId"), patch)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) ResendInvitation(c *gin.Context) {
	if err := h.users.ResendInvitation(c.Request.Context(), c.Param("tenantId"), c.Param("userId")); err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "queued", "deliveryChannel": "email"})
}

func (h *UserHandler) CreateImportJob(c *gin.Context) {
	var req models.CSVImportJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	job, err := h.imports.StartJob(c.Request.Context(), c.Param("tenantId"), req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, job)
}

func (h *UserHandler) GetImportJob(c *gin.Context) {
	c.JSON(http.StatusOK, models.CSVImportJob{
		JobID:        c.Param("jobId"),
		TenantID:     c.Param("tenantId"),
		Status:       "completed",
		TotalRows:    0,
		SuccessCount: 0,
		FailedCount:  0,
	})
}

func (h *UserHandler) ListGroups(c *gin.Context) {
	list, err := h.groups.List(c.Request.Context(), c.Param("tenantId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *UserHandler) CreateGroup(c *gin.Context) {
	var req models.GroupCreateInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	req.TenantID = c.Param("tenantId")
	group, err := h.groups.Create(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, database.ErrGroupNameExists) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, group)
}

func (h *UserHandler) PatchGroup(c *gin.Context) {
	var patch models.GroupPatchInput
	if err := c.ShouldBindJSON(&patch); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	group, err := h.groups.Patch(c.Request.Context(), c.Param("tenantId"), c.Param("groupId"), patch)
	if err != nil {
		if errors.Is(err, database.ErrGroupNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, group)
}

func (h *UserHandler) PutGroupMembers(c *gin.Context) {
	var payload struct {
		UserIDs []string `json:"userIds"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	group, err := h.groups.ReplaceMembers(c.Request.Context(), models.GroupMembersUpsert{
		TenantID: c.Param("tenantId"),
		GroupID:  c.Param("groupId"),
		UserIDs:  payload.UserIDs,
	})
	if err != nil {
		if errors.Is(err, database.ErrGroupNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, group)
}
