package config

import (
	"awesomeProject1/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var TestDB *gorm.DB

// Инициализация тестовой базы SQLite (в памяти)
func InitTestDB() {
	var err error
	TestDB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к тестовой базе:", err)
	}

	// Миграция моделей (замени на свои)
	err = TestDB.AutoMigrate(
		&models.User{},
		&models.Table_telegram_bot{}, // ✅ Добавляем сюда
		&models.Table_student{},      // ✅ Добавляем, если используется
		&models.Table_lecture{},      // ✅ Добавляем, если используется
	)
	if err != nil {
		log.Fatal("Ошибка миграции таблиц:", err)
	}
}
