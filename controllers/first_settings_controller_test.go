package controllers

import (
	"awesomeProject1/config"
	"awesomeProject1/models"
	"bytes"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Настройка тестового окружения
func setupFirstSettingsTestServer() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/firstsetting", FirstSettinPage)
	r.POST("/firstsetting", FirstSettingHandler)
	return r
}

// Подготовка тестовой базы
func setupFirstSettingsTestDB() {
	config.InitTestDB()
	config.DB = config.TestDB // Используем тестовую БД

	// Создаем все нужные таблицы
	config.DB.AutoMigrate(
		&models.User{},
		&models.Table_student{},
		&models.Table_lecture{},
		&models.Table_telegram_bot{},
	)
}

func TestFirstSettingHandler_Student(t *testing.T) {
	setupFirstSettingsTestDB() // Инициализируем тестовую БД
	r := setupFirstSettingsTestServer()

	// Создаем тестового пользователя
	user := models.User{
		Email:              "test@example.com",
		Verification_token: "test_token_123", // Добавляем уникальный токен
	}
	config.DB.Create(&user)

	token := generateTestJWT("test@example.com")
	// Устанавливаем тестовый JWT_SECRET
	os.Setenv("JWT_SECRET", "testsecret")

	// Создаем тестовый запрос для добавления ученика
	requestBody := models.PostSettings{
		Meaning: "Тестовый Ученик",
		Marking: "1", // 1 – это ученик
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/firstsetting", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "token", Value: token}) // Симуляция JWT-токена

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем, что сервер вернул HTTP 200 OK
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)

	// Проверяем, что ученик был добавлен в БД
	var student models.Table_student
	config.DB.First(&student)
	assert.Equal(t, "Тестовый Ученик", student.Name_Student)
}

func TestFirstSettingHandler_Lecture(t *testing.T) {
	setupFirstSettingsTestDB() // Инициализируем тестовую БД
	r := setupFirstSettingsTestServer()

	// Создаем тестового пользователя с уникальным verification_token
	user := models.User{
		Email:              "test@example.com",
		Verification_token: "unique_token_123", // ✅ Исправляем ошибку
	}
	config.DB.Create(&user)

	// Устанавливаем тестовый JWT_SECRET
	os.Setenv("JWT_SECRET", "testsecret")

	// Генерируем JWT-токен с email
	token := generateTestJWT("test@example.com")

	// Создаем тестовый запрос для добавления лекции
	requestBody := models.PostSettings{
		Meaning: "Тестовая Лекция",
		Marking: "0", // 0 – это лекция
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/firstsetting", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "token", Value: token}) // ✅ Передаем настоящий JWT-токен

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем, что сервер вернул HTTP 200 OK
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)

	// Проверяем, что лекция добавилась
	var lecture models.Table_lecture
	config.DB.First(&lecture)
	assert.Equal(t, "Тестовая Лекция", lecture.Lecture)
}
func generateTestJWT(email string) string {
	secret := "testsecret"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}
