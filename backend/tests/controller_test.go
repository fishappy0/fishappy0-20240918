package main

import (
	"CryptWatchBE/routes"
	"CryptWatchBE/utils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	app_mode := os.Getenv("APP_MODE")
	env_file_path := ""
	if app_mode == "" || app_mode == "development" {
		env_file_path = "../../stack.env"
	}

	connection_string := utils.ParseConnectionStringFromEnv(env_file_path)
	dbo := utils.ConnectToDB(connection_string)
	router := gin.Default()
	routes.GeneralRouter(router, dbo)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse the response")
	}
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Feeling good!", response["message"])
}

func TestGetTrending(t *testing.T) {
	app_mode := os.Getenv("APP_MODE")
	env_file_path := ""
	if app_mode == "" || app_mode == "development" {
		env_file_path = "../../stack.env"
	}

	connection_string := utils.ParseConnectionStringFromEnv(env_file_path)
	dbo := utils.ConnectToDB(connection_string)
	router := gin.Default()
	routes.ListRouter(router, dbo)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/trending", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestGetListCryptos(t *testing.T) {
	app_mode := os.Getenv("APP_MODE")
	env_file_path := ""
	if app_mode == "" || app_mode == "development" {
		env_file_path = "../../stack.env"
	}

	connection_string := utils.ParseConnectionStringFromEnv(env_file_path)
	dbo := utils.ConnectToDB(connection_string)
	router := gin.Default()
	routes.ListRouter(router, dbo)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/list_cryptos", nil)
	router.ServeHTTP(w, req)
	t.Log(w.Body.String())
	assert.Equal(t, 200, w.Code)
}

func TestGetCoinPrice(t *testing.T) {
	app_mode := os.Getenv("APP_MODE")
	env_file_path := ""
	if app_mode == "" || app_mode == "development" {
		env_file_path = "../../stack.env"
	}

	connection_string := utils.ParseConnectionStringFromEnv(env_file_path)
	dbo := utils.ConnectToDB(connection_string)
	router := gin.Default()
	routes.CryptoRouter(router, dbo)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/crypto/price?id=bitcoin", nil)
	router.ServeHTTP(w, req)
	t.Log(w.Body.String())
	assert.Equal(t, 200, w.Code)
	time.Sleep(5 * time.Second)
}

func TestGetCoinOHLC(t *testing.T) {
	app_mode := os.Getenv("APP_MODE")
	env_file_path := ""
	if app_mode == "" || app_mode == "development" {
		env_file_path = "../../stack.env"
	}

	connection_string := utils.ParseConnectionStringFromEnv(env_file_path)
	dbo := utils.ConnectToDB(connection_string)
	router := gin.Default()
	routes.CryptoRouter(router, dbo)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/crypto/ohlc?id=bitcoin", nil)
	router.ServeHTTP(w, req)
	t.Log(w.Body.String())
	assert.Equal(t, 200, w.Code)
	time.Sleep(5 * time.Second)
}

func TestGetCoinDetailedInfo(t *testing.T) {
	app_mode := os.Getenv("APP_MODE")
	env_file_path := ""
	if app_mode == "" || app_mode == "development" {
		env_file_path = "../../stack.env"
	}

	connection_string := utils.ParseConnectionStringFromEnv(env_file_path)
	dbo := utils.ConnectToDB(connection_string)
	router := gin.Default()
	routes.CryptoRouter(router, dbo)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/crypto/detailed?id=bitcoin", nil)
	router.ServeHTTP(w, req)
	t.Log(w.Body.String())
	assert.Equal(t, 200, w.Code)
}
