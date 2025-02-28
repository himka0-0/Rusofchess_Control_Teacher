package controllers

import (
	"awesomeProject1/config"
	"awesomeProject1/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Настройка тестового окружения
func setupNotesTestServer() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/notelesson", NotelessonHandler)
	return r
}

// Подготовка тестовой базы
func setupNotesTestDB() {
	config.InitTestDB()
	config.DB = config.TestDB // Используем тестовую БД

	// Создаем все нужные таблицы
	config.DB.AutoMigrate(
		&models.User{},
		&models.Table_student{},
	)
}

func TestNotelessonHandler(t *testing.T) {
	setupNotesTestDB() // Инициализируем тестовую БД
	r := setupNotesTestServer()

	// Создаем тестового пользователя
	user := models.User{Email: "test@example.com"}
	config.DB.Create(&user)

	// Создаем тестового ученика
	student := models.Table_student{
		User_id:       user.ID,
		Name_Student:  "Тестовый Ученик",
		Theory:        0,
		Practice:      0,
		Tasks:         0,
		Alert_moduls:  true,
		Alert_payment: true,
		Payment:       5,
	}
	config.DB.Create(&student)

	// Отправляем POST-запрос на добавление отметки
	requestBody := models.PostModuls{
		Student_id:   uint(student.ID),
		Module:       "Теория",
		Lock_lecture: false,
	}
	body, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/notelesson", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем, что сервер вернул HTTP 200 OK
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)

	// Проверяем, что у студента обновились данные
	var updatedStudent models.Table_student
	config.DB.First(&updatedStudent, student.ID)
	assert.Equal(t, 1, updatedStudent.Theory) // Теория должна увеличиться
}
