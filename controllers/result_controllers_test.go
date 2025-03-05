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

func setupResultTestDB() {
	config.InitTestDB()
	config.DB = config.TestDB

	config.TestDB.Create(&models.User{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
	})

	config.TestDB.Create(&models.Table_student{
		ID:             20,
		User_id:        1,
		Name_Student:   "Test Student 1",
		Namber_lecture: 1,
	})
	config.TestDB.Create(&models.Table_lecture{
		ID:                20,
		User_id:           1,
		Lecture:           "Test lecture 1",
		Lecture_Person_id: 1,
	})
}

func TestResultPage(t *testing.T) {
	setupResultTestDB()
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.LoadHTMLGlob("../templates/*")

	r.GET("/result", func(c *gin.Context) {
		user := models.User{ID: 1, Email: "test@example.com"}
		c.Set("User", user)
		c.Set("email", user.Email)
		controllers.ResultPage(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/result", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Список")
	assert.Contains(t, w.Body.String(), "Test Student 1")
	assert.Contains(t, w.Body.String(), "Test lecture 1")
}
