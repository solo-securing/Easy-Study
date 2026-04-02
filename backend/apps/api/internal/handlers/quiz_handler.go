package handlers

import (
	"errors"
	"net/http"
	"time"

	"api/internal/database"
	"api/internal/middleware"
	"api/internal/models"
	"api/internal/services"

	"github.com/gin-gonic/gin"
)

type QuizHandler struct {
	quizzes *services.QuizService
}

func NewQuizHandler(quizzes *services.QuizService) *QuizHandler {
	return &QuizHandler{quizzes: quizzes}
}

func (h *QuizHandler) GetQuiz(c *gin.Context) {
	quiz, err := h.quizzes.GetQuiz(c.Request.Context(), c.Param("tenantId"), c.Param("quizId"))
	if err != nil {
		if errors.Is(err, database.ErrQuizNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, quiz)
}

func (h *QuizHandler) StartAttempt(c *gin.Context) {
	var req struct {
		ResumeAttemptID string `json:"resumeAttemptId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil && c.Request.ContentLength > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	attempt, err := h.quizzes.StartAttempt(c.Request.Context(), actorID(c), models.QuizAttemptStartInput{
		TenantID:        c.Param("tenantId"),
		QuizID:          c.Param("quizId"),
		UserID:          userIDFromClaims(c),
		ResumeAttemptID: req.ResumeAttemptID,
	})
	if err != nil {
		switch {
		case errors.Is(err, database.ErrQuizNotFound):
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, database.ErrQuizAttemptConflict):
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, attempt)
}

func (h *QuizHandler) GetAttempt(c *gin.Context) {
	attempt, err := h.quizzes.GetAttempt(c.Request.Context(), c.Param("tenantId"), c.Param("quizId"), c.Param("attemptId"))
	if err != nil {
		if errors.Is(err, database.ErrQuizAttemptNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, attempt)
}

func (h *QuizHandler) Autosave(c *gin.Context) {
	var req struct {
		Answers       []models.QuizAnswer `json:"answers"`
		ClientSavedAt time.Time           `json:"clientSavedAt"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	attempt, err := h.quizzes.Autosave(c.Request.Context(), actorID(c), models.QuizAutosaveInput{
		TenantID:      c.Param("tenantId"),
		QuizID:        c.Param("quizId"),
		AttemptID:     c.Param("attemptId"),
		UserID:        userIDFromClaims(c),
		Answers:       req.Answers,
		ClientSavedAt: req.ClientSavedAt,
	})
	if err != nil {
		switch {
		case errors.Is(err, database.ErrQuizAttemptNotFound):
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"attemptId": attempt.AttemptID,
		"status":    attempt.Status,
		"savedAt":   time.Now().UTC(),
	})
}

func (h *QuizHandler) Submit(c *gin.Context) {
	result, err := h.quizzes.Submit(c.Request.Context(), actorID(c), c.Param("tenantId"), c.Param("quizId"), c.Param("attemptId"), false)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrQuizAttemptNotFound):
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, result)
}

func userIDFromClaims(c *gin.Context) string {
	claims, ok := middleware.JWTClaims(c)
	if !ok {
		return "anonymous"
	}
	if raw, ok := claims["sub"].(string); ok && raw != "" {
		return raw
	}
	if raw, ok := claims["user_id"].(string); ok && raw != "" {
		return raw
	}
	return "anonymous"
}
