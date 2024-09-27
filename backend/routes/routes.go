package routes

import (
	"CryptWatchBE/controllers"

	gin "github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GeneralRouter(router *gin.Engine, input_db *gorm.DB) {
	// fmt.Println("GeneralRouter")
	gr := router.Group("/")
	gr.GET("/health", controllers.HealthCheck)
}

func ListRouter(router *gin.Engine, input_db *gorm.DB) {
	// fmt.Println("ListRouter")
	lr := router.Group("/")
	// lr.Use(IsAuthenticated)
	gr_db := controllers.GeneralDB{DB: input_db}
	{
		lr.GET("/list_cryptos", gr_db.GetCryptoList)
		lr.GET("/trending", gr_db.GetTrending)
		lr.GET("/supported_currencies", gr_db.GetSupportedCurrencies)
	}
}

func CryptoRouter(router *gin.Engine, input_db *gorm.DB) {
	// fmt.Println("CryptoRouter")
	cr := router.Group("/crypto")
	// cr.Use(IsAuthenticated)
	cr_db := controllers.CryptoDB{DB: input_db}
	{
		cr.GET("/price", cr_db.GetCoinPrice)
		cr.GET("/search", cr_db.SearchCoins)
		cr.GET("/ohlc", cr_db.GetCoinOHLC)
		cr.GET("/detailed", cr_db.GetCoinDetailedInfo)
		cr.GET("/conversion", cr_db.GetCoinConversions)
	}
}

func AccountRouter(router *gin.Engine, input_db *gorm.DB) {
	// fmt.Println("AccountRouter")
	ar := router.Group("/account")
	ar_db := controllers.AccountDB{DB: input_db}
	{
		ar.POST("/register", ar_db.Register)
		ar.POST("/login", ar_db.Login)
		// ar.GET("/logout", ar_db.Logout)
	}
}
