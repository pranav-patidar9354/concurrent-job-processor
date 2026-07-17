package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, statusCode int, message string, data interface{}) {

	c.JSON(statusCode, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func Error(c *gin.Context, statusCode int, message string) {

	c.JSON(statusCode, gin.H{
		"success": false,
		"error":   message,
	})
}

func InternalServerError(c *gin.Context) {

	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"error":   "Internal server error",
	})
}