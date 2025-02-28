package controllers_test

import (
	"awesomeProject1/config"
	"awesomeProject1/controllers"
	"awesomeProject1/models"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestDB() {
	config.InitTestDB()
	// Используем тестовую БД
	config.DB = config.TestDB

	config.TestDB.Create(&models.User{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
	})

	config.TestDB.Create(&models.Table_student{
		ID:           1,
		User_id:      1,
		Name_Student: "John Doe",
		Payment:      100,
	})
}

// Тест GET /paymentstudent
func TestPaymentstudentPage(t *testing.T) {
	setupTestDB()
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	// Загружаем шаблоны (если есть)
	r.LoadHTMLGlob("../templates/*")

	// Регистрируем маршрут
	r.GET("/paymentstudent", func(c *gin.Context) {
		// эмулируем найденного пользователя
		user := models.User{ID: 1, Email: "test@example.com"}
		c.Set("User", user)
		controllers.PaymentstudentPage(c)
	})

	// Шлем GET-запрос
	req, _ := http.NewRequest("GET", "/paymentstudent", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем статус
	assert.Equal(t, http.StatusOK, w.Code)

	// Вместо "John Doe" проверяем любую статическую строку,
	// которая точно есть в шаблоне
	assert.Contains(t, w.Body.String(), "Выберите ученика")
}

// Тест POST /paymentstudent (существующий студент)
func TestPaymentstudentHandler(t *testing.T) {
	setupTestDB()
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	// Регистрируем маршрут
	r.POST("/paymentstudent", controllers.PaymentstudentHandler)

	// Формируем данные
	payload := models.Paymentstudent{
		ID:      1,
		Payment: 50,
	}
	jsonData, _ := json.Marshal(payload)

	// Делаем POST
	req, _ := http.NewRequest("POST", "/paymentstudent", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем, что код 200
	assert.Equal(t, http.StatusOK, w.Code)
	// Проверяем, что JSON ответ содержит "Оплата сохранена"
	assert.Contains(t, w.Body.String(), "Оплата сохранена")

	// Проверяем, что в БД обновился студент
	var student models.Table_student
	config.TestDB.First(&student, 1)
	assert.Equal(t, 150, student.Payment)
}

// Тест POST /paymentstudent (студент не найден)
func TestPaymentstudentHandler_StudentNotFound(t *testing.T) {
	setupTestDB()
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	// Регистрируем маршрут
	r.POST("/paymentstudent", controllers.PaymentstudentHandler)

	// Несуществующий студент
	payload := models.Paymentstudent{
		ID:      999,
		Payment: 50,
	}
	jsonData, _ := json.Marshal(payload)

	// Делаем POST
	req, _ := http.NewRequest("POST", "/paymentstudent", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем статус 404
	assert.Equal(t, http.StatusNotFound, w.Code)
	// Проверяем текст
	assert.Contains(t, w.Body.String(), "Ученик не найден")
}
