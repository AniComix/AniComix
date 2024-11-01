package api

import (
	"github.com/AniComix/server/models"
	"github.com/AniComix/server/storage"
	"github.com/gin-gonic/gin"
)

func badRequest(c *gin.Context, message string) {
	c.JSON(400, gin.H{"error": message})
}

func unauthorized(c *gin.Context) {
	c.JSON(401, gin.H{"error": "unauthorized"})
}

func getCurrentUser(c *gin.Context) (models.User, error) {
	uid, _ := c.Get("uid")
	var user models.User
	if err := storage.DB().Where("id = ?", uid).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func internalServerError(c *gin.Context) {
	c.JSON(500, gin.H{"error": "internal server error"})
}
