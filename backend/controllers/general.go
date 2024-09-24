package controllers

import (
	models "CryptWatchBE/models"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GeneralDB struct {
	DB *gorm.DB
}

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Feeling good!",
	})
}

// //////////////////////
// Function name: GetCryptoList
// Description: This function is used to get the list of cryptocurrencies
// input: gin context, db object
// output: None
// //////////////////////
func (gdb *GeneralDB) GetCryptoList(c *gin.Context) {
	num_of_coins, err := strconv.Atoi(c.Query("num"))
	if err != nil {
		num_of_coins = 10
		log.Println("Error converting string to int, trace: ", err)
	}

	sort_by := c.Query("sort_by")
	if sort_by == "" || (sort_by != "asc" && sort_by != "desc") {
		sort_by = "cryptos_data.rank asc"
	}
	return_data := []struct {
		Name string
		models.CryptosData
	}{}

	tx := gdb.DB.Table("cryptos_data, cryptos").Select("cryptos.name, cryptos_data.*").Where("cryptos.crypt_id = cryptos_data.crypt_id").Joins("JOIN cryptos ON cryptos.crypt_id = cryptos_data.crypt_id").Order(sort_by).Limit(num_of_coins).Scan(&return_data)
	if tx.Error != nil {
		if tx.Error.Error() != "record not found" {
			log.Println("Error fetching data from database, trace: ", tx.Error)
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

	c.JSON(200, return_data)
}

// //////////////////////
// Function name: GetTrending
// Description: This function is used to get the trending cryptocurrencies
// input: gin context, db object
// output: None
// //////////////////////
func (gdb *GeneralDB) GetTrending(c *gin.Context) {
	return_data := []models.CryptosData{}
	tx := gdb.DB.Table("cryptos_data").Select("cryptos.name, cryptos_data.*").Where("cryptos.crypt_id = cryptos_data.crypt_id").Joins("JOIN cryptos ON cryptos.crypt_id = cryptos_data.crypt_id").Where("rank < ?", 50).Order("rank asc").Scan(&return_data)
	if tx.Error != nil {
		if tx.Error.Error() != "record not found" {
			log.Println("Error fetching data from database, trace: ", tx.Error)
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
	c.JSON(200, return_data)
}
