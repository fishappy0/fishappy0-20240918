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

	tx := gdb.DB.
		Table("cryptos_data, cryptos").
		Select("cryptos.name, cryptos_data.*").
		Where("cryptos.crypt_id = cryptos_data.crypt_id").
		Limit(num_of_coins).
		Scan(&return_data)

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
// Function name: GetSupportedCurrencies
// Description: This function is used to get the list of supported currencies
// input: gin context, db object
// output: None
// //////////////////////
func (gdb *GeneralDB) GetSupportedCurrencies(c *gin.Context) {
	db_resp := []models.Fiats{}
	tx := gdb.DB.Find(&db_resp)
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
	return_data := []struct {
		Name   string
		Symbol string
	}{}
	for _, data := range db_resp {
		return_data = append(return_data, struct {
			Name   string
			Symbol string
		}{data.Name, data.Symbol})
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
	db_resp := []struct {
		models.Cryptos
	}{}
	tx := gdb.DB.
		Table("cryptos_data, cryptos").
		Select("cryptos.name, cryptos.symbol, cryptos.crypt_id").
		Where("cryptos.crypt_id = cryptos_data.crypt_id").
		Where("cryptos_data.rank < ?", 15).
		Where("NOT cryptos_data.rank = ?", 0).
		Order("cryptos_data.rank asc").
		Scan(&db_resp)

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
	return_data := []struct {
		ID     string
		Name   string
		Symbol string
	}{}
	for _, data := range db_resp {
		return_data = append(return_data, struct {
			ID     string
			Name   string
			Symbol string
		}{data.Cryptos.Crypt_id, data.Name, data.Symbol})
	}
	c.JSON(200, return_data)
}
