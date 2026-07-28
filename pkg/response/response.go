package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type APIResponse[T any] struct {
	Status    string    `json:"status"`
	Message   string    `json:"message,omitempty"`
	Data      T         `json:"data,omitempty"`
	Error     any       `json:"error,omitempty"`
	Path      string    `json:"path"`
	Timestamp time.Time `json:"timestamp"`
}

func send[T any](c *gin.Context, statusCode int, status string, message string, data T, errDetail any) {
	c.JSON(statusCode, APIResponse[T]{
		Status:    status,
		Message:   message,
		Data:      data,
		Error:     errDetail,
		Path:      c.Request.URL.RequestURI(),
		Timestamp: time.Now().UTC(),
	})
}

func Success[T any](c *gin.Context, message string, data T) {
	send(c, http.StatusOK, "success", message, data, nil)
}

func Created[T any](c *gin.Context, message string, data T) {
	send(c, http.StatusCreated, "success", message, data, nil)
}
