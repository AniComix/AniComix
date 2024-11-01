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
	"path/filepath"
	"strconv"
	"strings"
)

const (
	salt     = "d`DWA*D7=875dD+D988~7"
	pageSize = 20
)

type loginAndRegisterJson struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	json := loginAndRegisterJson{}
	if err := c.BindJSON(&json); err != nil {
		badRequest(c, "invalid json")
		return
	}
	username := json.Username
	password := json.Password
	if username == "" || password == "" {
		badRequest(c, "username and password are required")
		return
	}

	var userCount int64
	if err := storage.DB().Model(&models.User{}).Count(&userCount).Error; err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}

	user := models.User{
		Username:     username,
		PasswordHash: hashPassword(password),
	}
	if userCount == 0 {
		user.IsAdmin = true
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
	json := loginAndRegisterJson{}
	if err := c.BindJSON(&json); err != nil {
		badRequest(c, "invalid json")
		return
	}
	username := json.Username
	password := json.Password
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

type changePasswordJson struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func ChangePassword(c *gin.Context) {
	var json changePasswordJson
	if err := c.BindJSON(&json); err != nil {
		badRequest(c, "invalid json")
		return
	}
	oldPassword := json.OldPassword
	newPassword := json.NewPassword
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

func GetUserInfo(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		badRequest(c, "user not found")
		return
	}
	var user models.User
	if err := storage.DB().Where("username = ?", username).First(&user).Error; err != nil {
		badRequest(c, "user not found")
		return
	}
	var avatar string
	if user.AvatarPath != "" {
		avatar = "/api/avatar/" + user.Username
	}
	c.JSON(200, gin.H{
		"username": user.Username,
		"nickname": user.Nickname,
		"avatar":   avatar,
		"bio":      user.Bio,
		"is_admin": user.IsAdmin,
	})
}

func GetAvatar(c *gin.Context) {
	username := c.Param("username")
	var user models.User
	if err := storage.DB().Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}
	if user.AvatarPath == "" {
		c.JSON(404, gin.H{"error": "avatar not found"})
		return
	}
	c.File(filepath.Join(storage.DataDir(), "avatars", user.AvatarPath))
}

type setIsAdminJson struct {
	IsAdmin  bool   `json:"is_admin"`
	Username string `json:"username"`
}

func SetIsAdmin(c *gin.Context) {
	currentUser, err := getCurrentUser(c)
	if err != nil || !currentUser.IsAdmin {
		unauthorized(c)
		return
	}

	var json setIsAdminJson
	if err := c.BindJSON(&json); err != nil {
		badRequest(c, "invalid json")
		return
	}
	username := json.Username
	isAdmin := json.IsAdmin
	if username == "" {
		badRequest(c, "username is required")
		return
	}
	var user models.User
	if err := storage.DB().Where("username = ?", username).First(&user).Error; err != nil {
		badRequest(c, "user not found")
		return
	}
	user.IsAdmin = isAdmin
	if err := storage.DB().Save(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}

func ListUsers(c *gin.Context) {
	currentUser, err := getCurrentUser(c)
	if err != nil || !currentUser.IsAdmin {
		unauthorized(c)
		return
	}

	var users []models.User
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		badRequest(c, "invalid page")
		return
	}
	var count int64
	if err := storage.DB().Model(&models.User{}).Count(&count).Error; err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	if err := storage.DB().Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	result := make([]gin.H, len(users))
	for i, user := range users {
		var avatar string
		if user.AvatarPath != "" {
			avatar = "/api/avatar/" + user.Username
		}
		result[i] = gin.H{
			"username": user.Username,
			"nickname": user.Nickname,
			"bio":      user.Bio,
			"is_admin": user.IsAdmin,
			"avatar":   avatar,
		}
	}
	c.JSON(200, gin.H{
		"message":  "success",
		"users":    result,
		"max_page": (count + pageSize - 1) / pageSize,
	})
}
