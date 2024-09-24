package controllers

import (
	models "CryptWatchBE/models"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AccountDB struct {
	DB *gorm.DB
}

func (adb *AccountDB) Register(c *gin.Context) {
	user := models.Account{}
	c.BindJSON(&user)
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal server error",
		})
		return
	}
	user.Password = string(hashed_password)
	tx := adb.DB.Create(&user)
	if tx.Error != nil {
		c.JSON(500, gin.H{
			"message": "Internal server error",
		})
		return
	}
	c.JSON(200, gin.H{
		"message":  "User created successfully",
		"username": user.Username,
		"password": user.Password,
	})
}

func (adb *AccountDB) Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	jwt_phrase := os.Getenv("JWT_SECRET")

	user := models.Account{
		Username: username,
		Password: password,
	}
	c.BindJSON(&user)
	tx := adb.DB.Where("username = ?", user.Username).First(&user)
	if tx.Error != nil {
		c.JSON(500, gin.H{
			"message": "Internal server error",
		})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(user.Password))
	if err != nil {
		c.JSON(401, gin.H{
			"message": "Invalid credentials",
		})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
	})
	tokenString, err := token.SignedString([]byte(jwt_phrase))
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal server error",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Login successful",
		"token":   tokenString,
	})
}
