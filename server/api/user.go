package api

import (
	"crypto/sha256"
	"github.com/AniComix/server"
	"github.com/AniComix/server/models"
	"github.com/gin-gonic/gin"
)

const (
	salt = "d`DWA*D7=875dD+D988~7"
)

func Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		badRequest(c, "username and password are required")
		return
	}

	user := models.User{
		Username:     username,
		PasswordHash: hashPassword(password),
	}
	if err := server.DB().Create(&user).Error; err == nil {
		badRequest(c, "username already exists")
		return
	}
	token, err := server.GenerateToken(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(200, gin.H{"token": token})
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		badRequest(c, "username and password are required")
		return
	}

	var user models.User
	if err := server.DB().Where("username = ?", username).First(&user).Error; err != nil {
		unauthorized(c)
		return
	}
	if !checkPassword(password, user.PasswordHash) {
		unauthorized(c)
		return
	}
	token, err := server.GenerateToken(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(200, gin.H{"token": token})
}

func hashPassword(password string) []byte {
	hash := sha256.New()
	password += salt
	hash.Write([]byte(password))
	return hash.Sum(nil)
}

func checkPassword(password string, hash []byte) bool {
	return string(hash) == string(hashPassword(password))
}
