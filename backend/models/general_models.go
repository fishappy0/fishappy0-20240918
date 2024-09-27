package model

type Cryptos struct {
	Crypt_id string `gorm:"primaryKey"`
	Name     string `gorm:"type:varchar(255)"`
	Symbol   string `gorm:"type:varchar(255)"`

	CryptosData CryptosData `gorm:"foreignKey:Crypt_id;references:Crypt_id"`
	Price       Price       `gorm:"foreignKey:Crypt_id;references:Crypt_id"`
	OHLC        OHLC        `gorm:"foreignKey:Crypt_id;references:Crypt_id"`
	Conversions Conversions `gorm:"foreignKey:Crypt_id;references:Crypt_id"`
}

type CryptosData struct {
	Crypt_id    string `gorm:"primaryKey;"`
	Volume      float32
	Price       float32
	Rank        int
	Supply      float32
	Market_cap  float32
	Update_time int
}

type Price struct {
	Crypt_id string `gorm:"primaryKey"`
	Stamp    int    `gorm:"primaryKey"`
	Price    float64
	Type     string `gorm:"type:varchar(255)"`
}

type OHLC struct {
	Crypt_id string `gorm:"primaryKey"`
	Stamp    int    `gorm:"primaryKey"`
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Type     string `gorm:"type:varchar(255)"`
}

type Conversions struct {
	Crypt_id    string `gorm:"primaryKey"`
	Symbol      string `gorm:"primaryKey;type:varchar(255)"`
	Rate        float64
	Update_time int
}

type Fiats struct {
	Symbol string `gorm:"primaryKey;type:varchar(255)"`
	Name   string `gorm:"type:varchar(255)"`

	Conversions Conversions `gorm:"foreignKey:Symbol;references:Symbol"`
}

type Accounts struct {
	Account_id int    `gorm:"primaryKey;auto_increment"`
	Username   string `gorm:"type:varchar(255)"`
	Password   string `gorm:"type:varchar(255)"`
	Email      string `gorm:"type:varchar(255)"`
}
