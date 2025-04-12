package api

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"message": "success",
		"data":    data,
	})
}

func Fail(c *gin.Context, message string) {
	c.JSON(400, gin.H{
		"message": message,
	})
}
