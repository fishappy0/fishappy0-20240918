package model

type Account struct {
	Account_id int    `gorm:"primaryKey"`
	Username   string `gorm:"type:varchar(255)"`
	Password   string `gorm:"type:varchar(255)"`
}
