package server

import (
	"github.com/AniComix/server/api"
	"github.com/AniComix/server/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CORS中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许的源站（可以设置为具体的域名）
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func Run() {
	storage.InitStorage()
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.StaticFS("/data", http.Dir("data"))

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

	err := r.Run(":1919")
	if err != nil {
		panic(err)
	}
}
