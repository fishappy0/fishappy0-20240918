package types

type YamlConfig struct {
	LogLevel string `yaml:"log-level"`
	DbServer struct {
		Ip       string `yaml:"ip"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Db_name  string `yaml:"db_name"`
	} `yaml:"db-server"`
	LocalServer struct {
		Port string `yaml:"port"`
	} `yaml:"local-server"`
	APIURL              string `yaml:"api-url"`
	SupportedCurrencies []struct {
		Name   string `yaml:"name"`
		Symbol string `yaml:"symbol"`
	} `yaml:"supported-currencies"`
}

type ArrayStampVal struct {
	Timestamp int     `json:"timestamp"`
	Price     float64 `json:"price"`
}

type ArrayOHLC struct {
	Timestamp int     `json:"timestamp"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
}

type JsonPricesResponse struct {
	Prices       []ArrayStampVal `json:"prices"`
	MarketCaps   []ArrayStampVal `json:"market_caps"`
	TotalVolumne []ArrayStampVal `json:"total_volume"`
}
