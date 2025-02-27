package models

import "github.com/golang-jwt/jwt/v5"

type User struct {
	ID                  uint   `gorm:"primarykey"`
	Name                string `gorm:"size 100;not null"`
	Email               string `gorm:"size 100;unique;not null"`
	Password            string `gorm:"not null"`
	Lectures_introduced int    `gorm:"size 100"`
	Email_verified      bool   `gorm:"default:false"`
	Verification_token  string `gorm:"size 100;unique"`
	Table_student       []Table_student
	Table_lecture       []Table_lecture
	Table_telegram_bot  []Table_telegram_bot
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}
