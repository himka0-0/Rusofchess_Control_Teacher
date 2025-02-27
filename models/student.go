package models

type Table_student struct {
	ID             int    `gorm:"primarykey"`
	User_id        uint   `gorm:"index"`
	Name_Student   string `gorm:"size 100;not null"`
	Payment        int    `gorm:"size 100"`
	Theory         int    `gorm:"size 1000"`
	Practice       int    `gorm:"size 1000"`
	Tasks          int    `gorm:"size 1000"`
	Namber_lecture int    `gorm:"size 1000"`
	Alert_payment  bool   `json:"Alertpayment" gorm:"size:1000"`
	Alert_moduls   bool   `json:"Alertmodules" gorm:"size:1000"`
}
