package controllers_test

import (
	"awesomeProject1/config"
	"awesomeProject1/controllers"
	"awesomeProject1/models"
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
		User_id: 1,
		Vhod:    true,
	})

	config.TestDB.Create(&models.Table_student{
		ID:             20,
		User_id:        1,
		Name_Student:   "Test Student 1",
		Namber_lecture: 1,
	})
	config.TestDB.Create(&models.Table_student{
		ID:             21,
		User_id:        1,
		Name_Student:   "Test Student 2",
		Namber_lecture: 1,
	})
	config.TestDB.Create(&models.Table_student{
		ID:             22,
		User_id:        1,
		Name_Student:   "Test Student 3",
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
