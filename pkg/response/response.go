package response

import (
	"github.com/gin-gonic/gin"
)

// Success - 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"data": data,
	})
}

// Error - 错误响应
func Error(c *gin.Context, status int, code string, message string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}
