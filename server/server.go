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
		user := apiGroup.Group("/user")
		{
			user.POST("/register", api.Register)
			user.POST("/login", api.Login)
			user.POST("/update", api.UpdateUserInfo)
			user.POST("/changePassword", api.ChangePassword)
		}
	}
}
