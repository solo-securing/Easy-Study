package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"api/internal/database"
	"api/internal/models"
	"api/internal/services"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	courses     *services.CourseService
	enrollments *services.EnrollmentService
}

func NewCourseHandler(courses *services.CourseService, enrollments *services.EnrollmentService) *CourseHandler {
	return &CourseHandler{
		courses:     courses,
		enrollments: enrollments,
	}
}

func (h *CourseHandler) ListCourses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	status := models.CourseStatus(c.Query("status"))

	result, err := h.courses.List(c.Request.Context(), models.CourseListFilter{
		TenantID: c.Param("tenantId"),
		Status:   status,
		Page:     page,
		Size:     size,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var req models.CourseCreateInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	req.TenantID = c.Param("tenantId")

	course, err := h.courses.Create(c.Request.Context(), actorID(c), req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, course)
}

func (h *CourseHandler) GetCourse(c *gin.Context) {
	course, err := h.courses.Get(c.Request.Context(), c.Param("tenantId"), c.Param("courseId"))
	if err != nil {
		if errors.Is(err, database.ErrCourseNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, course)
}

func (h *CourseHandler) PatchCourse(c *gin.Context) {
	var patch models.CoursePatchInput
	if err := c.ShouldBindJSON(&patch); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	course, err := h.courses.Patch(c.Request.Context(), actorID(c), c.Param("tenantId"), c.Param("courseId"), patch)
	if err != nil {
		if errors.Is(err, database.ErrCourseNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, course)
}

func (h *CourseHandler) PutCourseStructure(c *gin.Context) {
	var req models.CourseStructureUpsertInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	course, err := h.courses.ReplaceStructure(c.Request.Context(), actorID(c), c.Param("tenantId"), c.Param("courseId"), req)
	if err != nil {
		if errors.Is(err, database.ErrCourseNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, course)
}

func (h *CourseHandler) PublishCourse(c *gin.Context) {
	course, issues, err := h.courses.Publish(c.Request.Context(), actorID(c), c.Param("tenantId"), c.Param("courseId"))
	if err != nil {
		if errors.Is(err, services.ErrCoursePublishValidation) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"code":         "COURSE_STRUCTURE_INVALID",
				"message":      "course structure must contain section/subsection/unit hierarchy",
				"missingNodes": issues,
			})
			return
		}
		if errors.Is(err, database.ErrCourseNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"courseId":    course.ID,
		"status":      course.Status,
		"publishedAt": time.Now().UTC(),
	})
}

func (h *CourseHandler) ArchiveCourse(c *gin.Context) {
	course, err := h.courses.Archive(c.Request.Context(), actorID(c), c.Param("tenantId"), c.Param("courseId"))
	if err != nil {
		if errors.Is(err, database.ErrCourseNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"courseId":   course.ID,
		"status":     course.Status,
		"archivedAt": time.Now().UTC(),
	})
}

func (h *CourseHandler) ListEnrollments(c *gin.Context) {
	list, err := h.enrollments.List(c.Request.Context(), c.Param("tenantId"), c.Param("courseId"))
	if err != nil {
		if errors.Is(err, database.ErrCourseNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *CourseHandler) AssignEnrollments(c *gin.Context) {
	var req models.EnrollmentAssignInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	req.TenantID = c.Param("tenantId")
	req.CourseID = c.Param("courseId")

	list, err := h.enrollments.Assign(c.Request.Context(), actorID(c), req)
	if err != nil {
		if errors.Is(err, database.ErrCourseNotFound) || errors.Is(err, database.ErrUserNotFound) || errors.Is(err, database.ErrGroupNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}
