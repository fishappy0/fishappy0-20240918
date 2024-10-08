package routines

import (
	model "CryptWatchBE/models"
	"CryptWatchBE/utils"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/gin-gonic/gin"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

func CacheDataAndConversion(db *gorm.DB, wg *sync.WaitGroup) {
	// Steps:
	// 1. Fetch infos like name, symbol to fill in the Cryptos table,
	// AND also cryptos data like volumne, market cap, etc to fill in the CryptosData table in USD
	// via the /coins/market?vs_currency=usd endpoint
	// 2. Fetch the exchange rates from the /simple/prices?ids=bitcoin,ethereum,doge,etc&vs_currencies=usd,eur,sgd,etc endpoint
	// and fill the exchange rates to the Conversions table
	// This should generate about (20 pages(each 250, totalling 5000 coins)) + 10 for the exchange rates = 30 requests per hour
	// and 1 day would be 30 * 24 = 720 requests per day
	log.Println("Caching data and conversion")
	for pages := 1; pages <= 20; pages++ {
		time.Sleep(5 * time.Second)
		url := "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=250&page=" + strconv.Itoa(pages) + "&sparkline=false"
		var response []struct {
			ID                string    `json:"id"`
			Symbol            string    `json:"symbol"`
			Name              string    `json:"name"`
			CurrentPrice      float32   `json:"current_price"`
			MarketCap         float32   `json:"market_cap"`
			MarketCapRank     int       `json:"market_cap_rank"`
			TotalVolume       float32   `json:"total_volume"`
			CirculatingSupply float32   `json:"circulating_supply"`
			LastUpdated       time.Time `json:"last_updated"`
		}

		var error_response struct {
			Status struct {
				ErrorCode    int    `json:"error_code"`
				ErrorMessage string `json:"error_message"`
			} `json:"status"`
		}

		coins_market_list_bytes := utils.FetchDataFromApiAsJson(url, os.Getenv("COINGECKO_API_KEY"))
		// TO DO: Deal with the "empty slice found" error aka the struct above somehow is empty somewhere
		err := json.Unmarshal(coins_market_list_bytes, &response)
		if err != nil {
			err_resp_err := json.Unmarshal(coins_market_list_bytes, &error_response)
			if err_resp_err != nil {
				log.Println("Error unmarshalling the response from the api, trace: ", err)
			}
			if error_response.Status.ErrorCode == 429 {
				log.Println("Rate limit exceeded, sleeping for 30 seconds")
				time.Sleep(30 * time.Second)
				continue
			} else {
				log.Println("Error unmarshalling the response from the api, trace: ", err)
			}

		}
		cryptos_structs := []model.Cryptos{}
		cryptos_data_structs := []model.CryptosData{}

		for _, coin := range response {
			cryptos_structs = append(cryptos_structs, model.Cryptos{
				Crypt_id: coin.ID,
				Name:     coin.Name,
				Symbol:   coin.Symbol,
			})
			cryptos_data_structs = append(cryptos_data_structs, model.CryptosData{
				Crypt_id:    coin.ID,
				Volume:      coin.TotalVolume,
				Price:       coin.CurrentPrice,
				Rank:        coin.MarketCapRank,
				Supply:      coin.CirculatingSupply,
				Market_cap:  coin.MarketCap,
				Update_time: int(coin.LastUpdated.Unix()),
			})
		}
		// create if not exists
		db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&cryptos_structs)
		db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&cryptos_data_structs)
	}

	// utils.CacheCurrencyData(db, 480, 0)
	println("Finished caching data and conversion")
	defer wg.Done()
	time.Sleep(30 * time.Minute)
}
