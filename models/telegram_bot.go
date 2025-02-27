package models

type Table_telegram_bot struct {
	User_id     uint   `gorm:"index"`
	Hash        string `gorm:"size 100;unique;not null"`
	First_name  string `gorm:"size 100"`
	Telegram_id int64  `gorm:"size 100"`
	Vhod        bool   `gorm:"size 100;not null"`
}
