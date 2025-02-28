package controllers

import (
	"awesomeProject1/config"
	"awesomeProject1/models"
	"bytes"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Настройка тестового окружения
func setupTestServer() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/registration", RegHandler)
	r.POST("/authentication", AutHandler)
	return r
}

// Подготовка тестовой базы
func setupTestDB() {
	config.InitTestDB()
	config.DB = config.TestDB // Используем тестовую БД

	// Создаем все нужные таблицы
	config.DB.AutoMigrate(
		&models.User{},
		&models.Table_telegram_bot{}, // ✅ Добавляем таблицу Telegram
		&models.Table_student{},      // ✅ Добавляем, если используется
		&models.Table_lecture{},      // ✅ Добавляем, если используется
	)
}

func TestRegisterUser(t *testing.T) {
	setupTestDB() // Инициализируем тестовую БД
	r := setupTestServer()

	userData := models.User{
		Email:    "test@example.com",
		Password: "securepassword",
	}
	userPas, _ := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	UserDataHash := models.User{
		Email:    "test@example.com",
		Password: string(userPas),
	}
	body, _ := json.Marshal(UserDataHash)

	req, _ := http.NewRequest("POST", "/registration", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Регистрация успешна")
}

func TestLoginUser(t *testing.T) {
	setupTestDB() // Инициализируем тестовую БД
	r := setupTestServer()

	// Хешируем пароль перед сохранением в базу
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("securepassword"), bcrypt.DefaultCost)

	// Создаем пользователя в тестовой базе с хешированным паролем
	config.DB.Create(&models.User{
		Email:          "test@example.com",
		Password:       string(hashedPassword), // 👈 Используем хеш
		Email_verified: true,
	})

	loginData := models.User{
		Email:    "test@example.com",
		Password: "securepassword", // 👈 Обычный пароль, без хеша
	}

	body, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/authentication", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"redirect":"/kabinet"`)
}
