package main

import (
	models "CryptWatchBE/models"
	"CryptWatchBE/routes"
	"CryptWatchBE/routines"
	cw_utils "CryptWatchBE/utils"
	"log"
	"os"
	"strings"
	"sync"

	cors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func main() {
	// yaml_config := cw_utils.ReadYamlConfig("config.yaml")
	app_mode := os.Getenv("APP_MODE")
	env_file_path := ""
	if app_mode == "" || app_mode == "development" {
		env_file_path = "../stack.env"
	}

	connection_string := cw_utils.ParseConnectionStringFromEnv(env_file_path)
	dbo := cw_utils.ConnectToDB(connection_string)
	cw_utils.AutoMigrate(dbo)
	log.Println("Connected to database", connection_string)

	var supported_currencies []models.Fiats
	config_yaml := cw_utils.ReadYamlConfig("config.yaml")
	for _, config := range config_yaml.SupportedCurrencies {
		supported_currencies = append(supported_currencies, models.Fiats{
			Name:   config.Name,
			Symbol: strings.ToLower(config.Symbol),
		})
	}
	dbo.Clauses(clause.OnConflict{DoNothing: true}).Create(&supported_currencies)

	var wg sync.WaitGroup
	wg.Add(1)
	go routines.CacheDataAndConversion(dbo, &wg)

	router := gin.Default()
	// AccountRouter(router, dbo)
	// router.Use(IsAuthenticated)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))
	routes.ListRouter(router, dbo)
	routes.GeneralRouter(router, dbo)
	routes.CryptoRouter(router, dbo)
	// router.RunTLS(":"+config_yaml.LocalServer.Port, "./certs/cert.pem", "./certs/key.pem")
	router.Run(":" + config_yaml.LocalServer.Port)
	wg.Wait()
}
