package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProblemDetails struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

func WriteProblem(c *gin.Context, status int, problemType, title, detail string) {
	c.AbortWithStatusJSON(status, ProblemDetails{
		Type:   problemType,
		Title:  title,
		Status: status,
		Detail: detail,
	})
}

func InternalServerProblem(c *gin.Context, detail string) {
	WriteProblem(c, http.StatusInternalServerError, "about:blank", "Internal Server Error", detail)
}
