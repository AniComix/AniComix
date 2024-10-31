package api

import "github.com/gin-gonic/gin"

func badRequest(c *gin.Context, message string) {
	c.JSON(400, gin.H{"error": message})
}

func unauthorized(c *gin.Context) {
	c.JSON(401, gin.H{"error": "unauthorized"})
}
