package server

import (
	"github.com/gin-gonic/gin"
)

func Run() {
	initStorage()
	r := gin.Default()
	api := r.Group("/api", authMiddleware())
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "hello world",
			})
		})
	}
}
