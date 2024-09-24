package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func IsAuthenticated(c *gin.Context) {
	jwt_phrase := os.Getenv("JWT_SECRET")
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{
			"message": "No token provided",
		})
		c.Abort()
		return
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwt_phrase), nil
	})
	if err != nil {
		c.JSON(401, gin.H{
			"message": "Invalid token",
		})
		c.Abort()
		return
	}
	if !token.Valid {
		c.JSON(401, gin.H{
			"message": "Invalid token",
		})
		c.Abort()
		return
	}
	c.Next()
}
