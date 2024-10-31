package server

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"time"
)

var (
	secret = ""
)

const (
	maxTokenAge = 60 * 60 * 24 * 7 // 1 week
)

// / authMiddleware is a middleware that checks for a valid JWT token in the Authorization header
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "" {
			t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err == nil && t.Valid {
				claims, ok := t.Claims.(jwt.MapClaims)
				if ok {
					timestamp := int64(claims["time"].(float64))
					if time.Now().Unix()-timestamp > maxTokenAge {
						c.JSON(401, gin.H{"error": "token expired"})
						c.Abort()
						return
					}
					c.Set("id", int(claims["id"].(float64)))
					c.Next()
					return
				}
			}
		}
		c.Next()
	}
}

func generateSecret() {
	if secret == "" {
		chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		for i := 0; i < 64; i++ {
			secret += string(chars[rand.Intn(len(chars))])
		}
	}
}

// / generateToken generates a JWT token with the given id
func generateToken(id int) (string, error) {
	generateSecret()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   id,
		"time": time.Now().Unix(),
	})
	token, err := t.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}
