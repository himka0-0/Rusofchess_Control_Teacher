package models

type Table_lecture struct {
	ID                int    `gorm:"primarykey"`
	User_id           uint   `gorm:"index"`
	Lecture           string `gorm:"size 100;not null"`
	Lecture_Person_id int    `gorm:"size 100;not null"`
}
