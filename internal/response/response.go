package response

import "github.com/gin-gonic/gin"

func JSON(c *gin.Context, status int, success bool, message string, data interface{}, err string) {
	c.JSON(status, gin.H{
		"success": success,
		"message": message,
		"data":    data,
		"error":   err,
	})
}
