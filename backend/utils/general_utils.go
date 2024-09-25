package utils

import (
	model "CryptWatchBE/models"
	types "CryptWatchBE/types"
	io "io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// //////////////////////
// Function name: CreateGinLogFile
// Description: This function is used to create a gin log file
// input: path to the log file
// output: multiwriter object
// //////////////////////
func CreateGinLogFile(path string) io.Writer {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal("Unable to create gin log file, trace: ", err)
	}
	return io.MultiWriter(file, os.Stdout)
}

// //////////////////////
// Function name: ConnectToDB
// Description: This function connects to the database and returns a db object
// input: connection string
// output: db object
// //////////////////////
func ConnectToDB(connection_string string) *gorm.DB {
	db, db_conn_err := gorm.Open(postgres.Open(connection_string), &gorm.Config{})
	if db_conn_err != nil {
		if strings.Contains(db_conn_err.Error(), "SQLSTATE 3D000") {
			connection_string_without_db_name := strings.Split(connection_string, "dbname=")[0]
			db, db_conn_err = gorm.Open(postgres.Open(connection_string_without_db_name), &gorm.Config{})
			if db_conn_err != nil {
				log.Fatal("Unable to connnect to the database, trace: ", db_conn_err)
			}
			db.Exec("CREATE DATABASE " + strings.Split(connection_string, "dbname=")[1])
			db = ConnectToDB(connection_string)
		} else {
			log.Fatal("Unable to connnect to the database, trace: ", db_conn_err)
		}
	}
	return db
}

// //////////////////////
// Function name: readYamlConfig
// Description: This function is used to read config.yaml
//
//	and return a YamlConfig object
//
// input: path to config.yaml
// output: YamlConfig object
// //////////////////////
func ReadYamlConfig(path string) types.YamlConfig {
	yaml_file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	yaml_config := types.YamlConfig{}
	yaml_error := yaml.Unmarshal(yaml_file, &yaml_config)
	if yaml_error != nil {
		panic(yaml_error)
	}
	return yaml_config
}

// //////////////////////
// function name: parseDate
// Description: This function is used to parse date using
//
//	go's default time parse func, with the format YYYY-MM-DD
//
// input: input format (has to be the exact same date January 2nd, 2006 for format), refer to https://yourbasic.org/golang/format-parse-string-time-date-example/
//
//	and date string
//
// output: datatypes.Date object
// //////////////////////
func ParseDate(format string, input_date string) datatypes.Date {
	if input_date == "" {
		return datatypes.Date{}
	}
	if format == "dd-mm-yyyy" || format == "dd/mm/yyyy" {
		format = "02-01-2006"
	} else if format == "yyyy-mm-dd" || format == "yyyy/mm/dd" {
		format = "2006-01-02"
	}
	var date_value, error = time.Parse(format, input_date)
	if error != nil {
		log.Println("Error while parsing date: ", error)
		return datatypes.Date{}
	} else {
		return datatypes.Date(date_value)
	}
}

// //////////////////////
// Function name: parseConnectionString
// Description: This function is used to parse connection string
//
//	from the YamlConfig object and return a string
//
// input: YamlConfig object
// output: Connection string
// //////////////////////
func ParseConnectionString(yaml_config types.YamlConfig) string {
	db_ip := yaml_config.DbServer.Ip
	db_port := yaml_config.DbServer.Port
	db_username := yaml_config.DbServer.Username
	db_password := yaml_config.DbServer.Password
	db_name := yaml_config.DbServer.Db_name

	return "host=" + db_ip + " " +
		"port=" + db_port + " " +
		"user=" + db_username + " " +
		"password=" + db_password + " " +
		"dbname=" + db_name
}

func ParseConnectionStringFromEnv(envfile string) string {
	if envfile != "" {
		err := godotenv.Load(envfile)
		if err != nil {
			log.Fatal("Error loading .env file, error", err)
		}
	}

	db_ip := os.Getenv("POSTGRES_HOST")
	db_port := os.Getenv("POSTGRES_PORT")
	db_username := os.Getenv("POSTGRES_USER")
	db_password := os.Getenv("POSTGRES_PASSWORD")
	db_name := os.Getenv("POSTGRES_DB")

	return "host=" + db_ip + " " +
		"port=" + db_port + " " +
		"user=" + db_username + " " +
		"password=" + db_password + " " +
		"dbname=" + db_name
}

// //////////////////////
// Function name: parseTime
// Description: This function is used to parse time
//
//	with the format HH:MM:SS
//
// input: Time string
// output: datatypes.Time object
// //////////////////////
func ParseTime(time string) datatypes.Time {
	var split_string = strings.Split(time, ":")

	var hour, hour_err = strconv.Atoi(split_string[0])
	if hour_err != nil {
		log.Println("Error while parsing hour: ", hour_err)
		return datatypes.NewTime(0, 0, 0, 0)
	}

	var minute, min_err = strconv.Atoi(split_string[1])
	if min_err != nil {
		log.Println("Error while parsing hour: ", hour_err)
		return datatypes.NewTime(0, 0, 0, 0)
	}

	var second, sec_err = strconv.Atoi(split_string[2])
	if sec_err != nil {
		log.Println("Error while parsing hour: ", hour_err)
		return datatypes.NewTime(0, 0, 0, 0)
	}

	return datatypes.NewTime(hour, minute, second, 0)
}

// //////////////////////
// Function name: ParseUnixTime
// Description: This function is used to parse unix time
// Input: Unix time string, format: "1609459200"
// Output: time.Time object
// //////////////////////
func ParseUnixTime(input_time string) time.Time {
	if input_time == "" {
		return time.Time{}
	}
	unix_time, err := strconv.ParseInt(input_time, 10, 64)
	if err != nil {
		log.Panicln("Error while parsing unix time: ", err)
		return time.Time{}
	}
	return time.Unix(unix_time, 0)
}

// //////////////////////
// Function name: AutoMigrate
// Description: This function is used to migrate tables
// Input: db object
// Output: None
// //////////////////////
func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.Cryptos{},
		&model.Accounts{},
		&model.CryptosData{},
		&model.Fiats{},
		&model.Conversions{},
		&model.OHLC{},
		&model.Price{},
	)
	if err != nil {
		log.Fatal("Unable to migrate tables, trace: ", err)
	} else {
		log.Println("Tables migrated successfully")
	}
}

// //////////////////////
// Function name: FetchDataFromApiAsJson
// Description: This function is used to fetch data from the coin gecko api
// Input: URL, API key, Struct type
// Output: JSON string
// //////////////////////
func FetchDataFromApiAsJson(url string, api_key string) []byte {
	req, req_err := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("x-cg-demo-api-key", api_key)

	res, res_err := http.DefaultClient.Do(req)

	if res_err != nil {
		log.Panicln("Error while making request: ", req_err)
		return nil
	}

	defer res.Body.Close()
	body, body_err := io.ReadAll(res.Body)

	if body_err != nil {
		log.Panicln("Error while reading response body: ", body_err)
		return nil
	}

	return body
}
