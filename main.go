package main

import (
	"awesomeProject1/config"
	"awesomeProject1/models"
	"awesomeProject1/router"
	"awesomeProject1/telegram"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
}

func main() {
	config.InitDB()
	config.DB.AutoMigrate(&models.User{}, &models.Table_student{}, &models.Table_lecture{}, &models.Table_telegram_bot{})
	go telegram.RunBot()
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	router.SetupRoutes(r)

	r.Run(":8080")
}
