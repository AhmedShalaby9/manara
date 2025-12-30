package helpers

import (
	"net/http"

	"manara/models"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success    bool               `json:"success"`
	Data       interface{}        `json:"data,omitempty"`
	Message    string             `json:"message,omitempty"`
	Pagination *models.Pagination `json:"pagination,omitempty"`
}

func Respond(c *gin.Context, success bool, data interface{}, message string) {
	status := http.StatusOK
	if !success {
		status = http.StatusBadRequest
	}
	c.JSON(status, Response{
		Success: success,
		Data:    data,
		Message: message,
	})
}

func RespondWithPagin(c *gin.Context, success bool, data interface{}, message string, pagination models.Pagination) {
	status := http.StatusOK
	if !success {
		status = http.StatusBadRequest
	}
	c.JSON(status, Response{
		Success:    success,
		Data:       data,
		Message:    message,
		Pagination: &pagination,
	})
}

func RespondWithStatus(c *gin.Context, status int, success bool, data interface{}, message string) {
	c.JSON(status, Response{
		Success: success,
		Data:    data,
		Message: message,
	})
}

func RespondError(c *gin.Context, status int, message string) {
	c.JSON(status, Response{
		Success: false,
		Message: message,
	})
}

func RespondNotFound(c *gin.Context, message string) {
	RespondError(c, http.StatusNotFound, message)
}

func RespondBadRequest(c *gin.Context, message string) {
	RespondError(c, http.StatusBadRequest, message)
}

func RespondUnauthorized(c *gin.Context, message string) {
	RespondError(c, http.StatusUnauthorized, message)
}

func RespondInternalError(c *gin.Context, success bool, message string) {
	RespondError(c, http.StatusInternalServerError, message)
}

func RespondSuccess(c *gin.Context, data interface{}, message string) {
	Respond(c, true, data, message)
}

func RespondCreated(c *gin.Context, data interface{}, message string) {
	RespondWithStatus(c, http.StatusCreated, true, data, message)
}

func RespondUpdated(c *gin.Context, data interface{}, message string) {
	RespondWithStatus(c, http.StatusAccepted, true, data, message)
}

func RespondForbiden(c *gin.Context, message string) {
	RespondError(c, http.StatusForbidden, message)
}
func RespondTokenError(c *gin.Context, message string) {
	RespondError(c, http.StatusUnauthorized, message)
}
