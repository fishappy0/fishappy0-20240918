package controllers

import (
	models "CryptWatchBE/models"
	"CryptWatchBE/utils"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CryptoDB struct {
	DB *gorm.DB
}

// ////////////////////
// Function name: GetCoinPrice
// Description: This function is used to get the price of a coin to fill a chart
// input: gin context, db object
// input query: id, duration
// output: None
// ////////////////////
func (cdb *CryptoDB) GetCoinPrice(c *gin.Context) {
	coins_id := c.Query("id")
	duration := c.Query("duration")
	if duration == "" {
		duration = "1"
	}

	url := "https://api.coingecko.com/api/v3/coins/" + coins_id + "/market_chart?vs_currency=usd&days=" + duration
	bytes := utils.FetchDataFromApiAsJson(url, os.Getenv("API_KEY"))
	if bytes == nil {
		c.JSON(500, gin.H{
			"message": "Internal server error",
		})
		log.Panicln("Error fetching data from api")
		return
	}
	var api_resp struct {
		Price  [][]float64 `json:"prices"`
		Market [][]float64 `json:"market_caps"`
		Volume [][]float64 `json:"total_volumes"`
	}

	err := json.Unmarshal(bytes, &api_resp)
	if api_resp.Price == nil {
		c.JSON(404, gin.H{
			"message": "Coin not found",
		})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal server error",
		})
		log.Panicln("Error unmarshalling data from api")
		return
	}
	c.JSON(200, gin.H{
		"prices": api_resp.Price,
	})
}

// ////////////////////
// Function name: GetCoinOHLC
// Description: This function is used to get the OHLC data of a coin to fill a chart
// input: gin context, db object
// input query: id, duration
// output: None
// ////////////////////
func (cdb *CryptoDB) GetCoinOHLC(c *gin.Context) {
	coins_id := c.Query("id")
	duration := c.Query("duration")
	if duration == "" {
		duration = "1"
	}

	url := "https://api.coingecko.com/api/v3/coins/" + coins_id + "/ohlc?vs_currency=usd&days=" + duration
	bytes := utils.FetchDataFromApiAsJson(url, os.Getenv("API_KEY"))
	if bytes == nil {
		c.JSON(500, gin.H{
			"message": "Internal server error",
		})
		log.Panicln("Error fetching data from api")
		return
	}
	var api_resp [][]float64
	err := json.Unmarshal(bytes, &api_resp)
	if err != nil {
		if strings.Contains(string(bytes), "coin not found") {
			c.JSON(404, gin.H{
				"message": "Coin not found",
			})
			return
		}
		c.JSON(500, gin.H{
			"message": "Internal server error",
		})
		log.Panicln("Error unmarshalling data from api")
		return
	}
	c.JSON(200, gin.H{
		"prices": api_resp,
	})
}

// ////////////////////
// Function name: GetCoinDetailedInfo
// Description: This function is used to get the detailed information of a coin, market data, name, etc
// input: gin context, db object
// input query: name
// output: None
// ////////////////////
func (cdb *CryptoDB) GetCoinDetailedInfo(c *gin.Context) {
	coins_param := c.Query("name")
	db_resp := struct {
		models.Cryptos
		models.CryptosData
	}{}

	tx := cdb.DB.Table("cryptos_data").
		Select("cryptos.*, cryptos_data.*").
		Where("cryptos.crypt_id = cryptos_data.crypt_id AND cryptos.name = ?", coins_param).
		Joins("JOIN cryptos ON cryptos.crypt_id = cryptos_data.crypt_id").
		First(&db_resp)

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
	var response struct {
		Name   string
		Symbol string
		Rank   int
		Price  float64
		Volume float64
		Supply float64
		Market float64
	}
	response.Name = db_resp.Name
	response.Symbol = db_resp.Symbol
	response.Rank = db_resp.Rank
	response.Price = float64(db_resp.CryptosData.Price)
	response.Volume = float64(db_resp.Volume)
	response.Supply = float64(db_resp.Supply)
	response.Market = float64(db_resp.Market_cap)

	c.JSON(200, response)
}

// ////////////////////
// Function name: SearchCoins
// Description: This function is used to search for coins
// input: gin context, db object
// input query: name
// output: None
// ////////////////////
func (cdb *CryptoDB) SearchCoins(c *gin.Context) {
	coins_param := c.Query("name")
	db_resp := []models.Cryptos{}
	tx := cdb.DB.Table("cryptos").Where("cryptos.name LIKE ?", "%"+coins_param+"%").Scan(&db_resp)
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
	var result []struct {
		ID     string
		Name   string
		Symbol string
	}
	for _, coin := range db_resp {
		result = append(result, struct {
			ID     string
			Name   string
			Symbol string
		}{
			ID:     coin.Crypt_id,
			Name:   coin.Name,
			Symbol: coin.Symbol,
		})
	}

	c.JSON(200, result)
}
