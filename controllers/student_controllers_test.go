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

func setupStudentTestDB() {
	config.InitTestDB()
	config.DB = config.TestDB
	//тестовый юзер
	config.TestDB.Create(&models.User{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
	})
	//тестовые ученики
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
	//тестовые лекции
	config.TestDB.Create(&models.Table_lecture{
		ID:                11,
		User_id:           1,
		Lecture:           "Lecture A",
		Lecture_Person_id: 1,
	})
	config.TestDB.Create(&models.Table_lecture{
		ID:                12,
		User_id:           1,
		Lecture:           "Lecture B",
		Lecture_Person_id: 2,
	})
	config.TestDB.Create(&models.Table_lecture{
		ID:                13,
		User_id:           1,
		Lecture:           "Lecture C",
		Lecture_Person_id: 3,
	})
}

// тест Studentpage
func TestStudentPage(t *testing.T) {
	setupStudentTestDB()
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.LoadHTMLGlob("../templates/*")

	r.GET("/student", func(c *gin.Context) {
		user := models.User{ID: 1, Email: "test@example.com"}
		c.Set("User", user)
		controllers.StudentPage(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/student", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "учениками")
}
func TestStudentHandler_Updatenamelecture(t *testing.T) {
	setupStudentTestDB()
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.LoadHTMLGlob("../templates/*")

	r.POST("/student", func(c *gin.Context) {
		user := models.User{ID: 1, Email: "test@example.com"}
		c.Set("User", user)
		controllers.StudentHandler(c)
	})

	input := []models.PostStudent{
		{ID: 20, Name: "Test Student 1", Lecture: 2},
		{ID: 21, Name: "Test Student 2", Lecture: 1},
		{ID: 22, Name: "Test Student 3", Lecture: 1},
	}
	jsonData, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/student", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Всё сохранено")

	var student models.Table_student
	config.TestDB.First(&student, 20)
	assert.Equal(t, 2, student.Namber_lecture)
}
func TestStudentHandler_Delete(t *testing.T) {
	setupStudentTestDB()
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.LoadHTMLGlob("../templates/*")

	r.POST("/student", func(c *gin.Context) {
		user := models.User{ID: 1, Email: "test@example.com"}
		c.Set("User", user)
		controllers.StudentHandler(c)
	})

	input := []models.Table_student{
		{ID: 20, Name_Student: "Test Student 1", Namber_lecture: 1},
		{ID: 21, Name_Student: "Test Student 2", Namber_lecture: 1},
	}

	jsonData, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/student", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Всё сохранено")

	var students []models.Table_student
	config.TestDB.Order("ID asc").Find(&students)
	assert.Len(t, students, 2)
	assert.Equal(t, students[0].Name_Student, "Test Student 1")
	assert.Equal(t, students[1].Name_Student, "Test Student 2")

}
func TestStudentHandler_Add(t *testing.T) {
	setupStudentTestDB()
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.LoadHTMLGlob("../templates/*")

	r.POST("/student", func(c *gin.Context) {
		user := models.User{ID: 1, Email: "test@example.com"}
		c.Set("User", user)
		controllers.StudentHandler(c)
	})

	input := []models.PostStudent{
		{ID: 20, Name: "Test Student 1", Lecture: 1},
		{ID: 21, Name: "Test Student 2", Lecture: 1},
		{ID: 22, Name: "Test Student 3", Lecture: 1},
		{ID: 0, Name: "Test Student 4", Lecture: 1},
	}

	jsonData, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/student", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Всё сохранено")

	var students []models.Table_student
	config.TestDB.Order("ID asc").Find(&students)
	assert.Len(t, students, 4)
	assert.Equal(t, students[0].Name_Student, "Test Student 1")
	assert.Equal(t, students[1].Name_Student, "Test Student 2")
	assert.Equal(t, students[2].Name_Student, "Test Student 3")
	assert.Equal(t, students[3].Name_Student, "Test Student 4")
}
