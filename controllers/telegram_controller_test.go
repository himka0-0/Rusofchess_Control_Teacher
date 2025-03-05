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

func setupTelbotTestDB() {
	config.InitTestDB()
	config.DB = config.TestDB

	config.TestDB.Create(&models.User{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
	})

	config.TestDB.Create(&models.Table_telegram_bot{
		User_id:     1,
		Hash:        "admskakmsm",
		First_name:  "himkas",
		Telegram_id: 80908810990,
		Vhod:        true,
	})

	config.TestDB.Create(&models.Table_student{
		ID:             20,
		User_id:        1,
		Name_Student:   "Test Student 1",
		Namber_lecture: 1,
	})
}

func TestTelbotPage(t *testing.T) {
	setupTelbotTestDB()
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.LoadHTMLGlob("../templates/*")

	r.GET("/telbot", func(c *gin.Context) {
		user := models.User{ID: 1, Email: "test@example.com"}
		c.Set("User", user)
		controllers.TelbotPage(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/telbot", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Настройка")
}

func TestTelbotHandler_AllertTeacher(t *testing.T) {
	setupTelbotTestDB()
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.LoadHTMLGlob("../templates/*")

	r.POST("/telbot", func(c *gin.Context) {
		user := models.User{ID: 1, Email: "test@example.com"}
		c.Set("User", user)
		controllers.TelbotHandler(c)
	})

	input := []models.PostTelbot{
		{ModuleAllToggle: false},
	}
	jsonData, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/telbot", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Всё сохранено")

	var user models.Table_telegram_bot
	config.TestDB.First(&user, 1)
	assert.Equal(t, false, user.Vhod)
}

func TestTelbotHandler_AllertModulsStudent(t *testing.T) {
	setupTelbotTestDB()
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.LoadHTMLGlob("../templates/*")

	r.POST("/telbot", func(c *gin.Context) {
		user := models.User{ID: 1, Email: "test@example.com"}
		c.Set("User", user)
		controllers.TelbotHandler(c)
	})

	input := []models.PostTelbot{
		{Students: []models.Table_student{
			{ID: 20, Name_Student: "Test Student 1", Alert_payment: true, Alert_moduls: false},
		}},
	}
	jsonData, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/telbot", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Всё сохранено")

	var student models.Table_student
	config.TestDB.First(&student, 20)
	assert.Equal(t, false, student.Alert_moduls)
}

func TestTelbotHandler_AllertPaymentStudent(t *testing.T) {
	setupTelbotTestDB()
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.LoadHTMLGlob("../templates/*")

	r.POST("/telbot", func(c *gin.Context) {
		user := models.User{ID: 1, Email: "test@example.com"}
		c.Set("User", user)
		controllers.TelbotHandler(c)
	})

	input := []models.PostTelbot{
		{Students: []models.Table_student{
			{ID: 20, Name_Student: "Test Student 1", Alert_payment: false, Alert_moduls: true},
		}},
	}
	jsonData, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/telbot", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Всё сохранено")

	var student models.Table_student
	config.TestDB.First(&student, 20)
	assert.Equal(t, false, student.Alert_payment)
}
