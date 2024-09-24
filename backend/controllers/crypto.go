package controllers

import (
	models "CryptWatchBE/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CryptoDB struct {
	DB *gorm.DB
}

func (cdb *CryptoDB) GetCoinPrice(c *gin.Context) {
	coins_param := c.Query("name")
	response := []models.Price{}
	tx := cdb.DB.Table("prices").Where("name = ?", coins_param).Scan(&response)
	if tx.Error != nil {
		if tx.Error.Error() != "record not found" {
			c.JSON(500, gin.H{
				"message": "Internal server error",
			})
			return
		} else {
			c.JSON(200, gin.H{
				"message": "No data found",
			})
			return
		}
	}

	c.JSON(200, response)
}

func (cdb *CryptoDB) GetCoinOHLC(c *gin.Context) {
	coins_param := c.Query("name")
	response := []models.OHLC{}
	tx := cdb.DB.Table("ohlcs").Where("name = ?", coins_param).Scan(&response)
	if tx.Error != nil {
		if tx.Error.Error() != "record not found" {
			c.JSON(500, gin.H{
				"message": "Internal server error",
			})
			return
		} else {
			c.JSON(200, gin.H{
				"message": "No data found",
			})
			return
		}
	}

	c.JSON(200, response)
}

func (cdb *CryptoDB) GetCoinDetailedInfo(c *gin.Context) {
	coins_param := c.Query("name")
	response := struct {
		models.Cryptos
		models.CryptosData
	}{}

	// tx := cdb.DB.Table("cryptos").Where("name = ?", coins_param).Scan(&response)
	tx := cdb.DB.Table("cryptos_data, cryptos").Select("cryptos.*, cryptos_data.*").Where("cryptos.crypt_id = cryptos_data.crypt_id AND cryptos.name = ?", coins_param).Joins("JOIN cryptos ON cryptos.crypt_id = cryptos_data.crypt_id").Take(&response)
	if tx.Error != nil {
		if tx.Error.Error() != "record not found" {
			c.JSON(500, gin.H{
				"message": "Internal server error",
			})
			return
		} else {
			c.JSON(200, gin.H{
				"message": "No data found",
			})
			return
		}
	}

	c.JSON(200, response)
}

func (cdb *CryptoDB) SearchCoins(c *gin.Context) {
	coins_param := c.Query("name")
	response := []models.Cryptos{}
	tx := cdb.DB.Table("cryptos").Where("name LIKE ?", "%"+coins_param+"%").Scan(&response)
	if tx.Error != nil {
		if tx.Error.Error() != "record not found" {
			c.JSON(500, gin.H{
				"message": "Internal server error",
			})
			return
		} else {
			c.JSON(200, gin.H{
				"message": "No data found",
			})
			return
		}
	}
	c.JSON(200, response)
}
