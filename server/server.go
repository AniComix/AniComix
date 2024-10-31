package server

import (
	"github.com/AniComix/server/api"
	"github.com/AniComix/server/storage"
	"github.com/gin-gonic/gin"
)

func Run() {
	storage.InitStorage()
	r := gin.Default()
	apiGroup := r.Group("/apiGroup", api.AuthMiddleware())
	{
		apiGroup.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "hello world",
			})
		})
		apiGroup.POST("/register", api.Register)
		apiGroup.POST("/login", api.Login)
		apiGroup.POST("/user/update", api.UpdateUserInfo)
	}
}
