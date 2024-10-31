package api

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/AniComix/server/models"
	"github.com/AniComix/server/storage"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strings"
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
	if err := storage.DB().Create(&user).Error; err == nil {
		badRequest(c, "username already exists")
		return
	}
	token, err := GenerateToken(user.ID)
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
	if err := storage.DB().Where("username = ?", username).First(&user).Error; err != nil {
		unauthorized(c)
		return
	}
	if !checkPassword(password, user.PasswordHash) {
		unauthorized(c)
		return
	}
	token, err := GenerateToken(user.ID)
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
	hash1 := hashPassword(password)
	if len(hash1) != len(hash) {
		return false
	}
	for i := 0; i < len(hash); i++ {
		if hash1[i] != hash[i] {
			return false
		}
	}
	return true
}

type updateUserInfoJson struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
}

func UpdateUserInfo(c *gin.Context) {
	var json updateUserInfoJson
	if err := c.BindJSON(&json); err != nil {
		badRequest(c, "invalid json")
		return
	}

	uid, exists := c.Get("uid")
	if !exists {
		unauthorized(c)
		return
	}
	var user models.User
	if err := storage.DB().Where("id = ?", uid).First(&user).Error; err != nil {
		badRequest(c, "user not found")
		return
	}

	if json.Nickname != "" {
		user.Nickname = json.Nickname
	}
	if json.Avatar != "" {
		reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(json.Avatar))
		bytes, err := io.ReadAll(reader)
		if err != nil {
			badRequest(c, "invalid base64")
			return
		}
		contentType := http.DetectContentType(bytes)
		var ext string
		switch contentType {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		case "image/webp":
			ext = ".webp"
		default:
			badRequest(c, "invalid image format")
		}
		avatarPath := storage.DataDir() + "/avatars/" + user.Username + ext
		file, err := os.Create(avatarPath)
		if err != nil {
			c.JSON(500, gin.H{"error": "internal server error"})
			return
		}
		if _, err := file.Write(bytes); err != nil {
			c.JSON(500, gin.H{"error": "internal server error"})
			return
		}
		if err := file.Close(); err != nil {
			c.JSON(500, gin.H{"error": "internal server error"})
			return
		}
		user.AvatarPath = avatarPath
	}
	if json.Bio != "" {
		user.Bio = json.Bio
	}

	if err := storage.DB().Save(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}

func ChangePassword(c *gin.Context) {
	oldPassword := c.PostForm("old_password")
	newPassword := c.PostForm("new_password")
	if oldPassword == "" || newPassword == "" {
		badRequest(c, "old_password and new_password are required")
		return
	}

	uid, exists := c.Get("uid")
	if !exists {
		unauthorized(c)
		return
	}
	var user models.User
	if err := storage.DB().Where("id = ?", uid).First(&user).Error; err != nil {
		badRequest(c, "user not found")
		return
	}
	if !checkPassword(oldPassword, user.PasswordHash) {
		unauthorized(c)
		return
	}
	user.PasswordHash = hashPassword(newPassword)
	if err := storage.DB().Save(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}
