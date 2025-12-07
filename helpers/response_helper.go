package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Respond(c *gin.Context, success bool, data interface{}, message string) {
	status := http.StatusOK
	if !success {
		status = http.StatusBadRequest
	}
	c.JSON(status, gin.H{
		"success": success,
		"data":    data,
		"message": message,
	})
}
