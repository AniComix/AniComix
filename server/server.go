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
		user := apiGroup.Group("/user")
		{
			user.POST("/register", api.Register)
			user.POST("/login", api.Login)
			user.POST("/update", api.UpdateUserInfo)
			user.POST("/changePassword", api.ChangePassword)
			user.GET("/:username", api.GetUserInfo)
			user.GET("/avatar/:username", api.GetAvatar)
			user.POST("/setAdmin", api.SetIsAdmin)
			user.GET("/list", api.ListUsers)
		}
	}

	err := r.Run(":4320")
	if err != nil {
		panic(err)
	}
}
