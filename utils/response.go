package utils

import (
	"github.com/gin-gonic/gin"
)

func ErrorResponse(c *gin.Context, request int, message string, error string) {
	c.JSON(request, gin.H{
		"success": false,
		"message": message,
		"error":   error,
	})
}

func SuccessResponse(c *gin.Context, request int, message string, data interface{}) {
	c.JSON(request, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}
