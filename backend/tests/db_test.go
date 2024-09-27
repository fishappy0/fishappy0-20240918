package main

import (
	cw_utils "CryptWatchBE/utils"
	"os"
	"testing"
)

func TestDBConnect(t *testing.T) {
	app_mode := os.Getenv("APP_MODE")
	env_file_path := ""
	if app_mode == "" || app_mode == "development" {
		env_file_path = "../../stack.env"
	}

	connection_string := cw_utils.ParseConnectionStringFromEnv(env_file_path)
	dbo := cw_utils.ConnectToDB(connection_string)
	if dbo == nil {
		t.Errorf("Failed to connect to database")
	}
}
