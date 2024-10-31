package server

import (
	"github.com/AniComix/server/api"
	"github.com/gin-gonic/gin"
)

func Run() {
	initStorage()
	r := gin.Default()
	apiGroup := r.Group("/apiGroup", authMiddleware())
	{
		apiGroup.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "hello world",
			})
		})
		apiGroup.POST("/register", api.Register)
		apiGroup.POST("/login", api.Login)
	}
}
