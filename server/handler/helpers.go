package handler

import (
	"github.com/gin-gonic/gin"
)

func JSONError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

func JSONSuccess(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

func GetUserID(c *gin.Context) (int, bool) {
	id, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	uid, ok := id.(int)
	return uid, ok
}
