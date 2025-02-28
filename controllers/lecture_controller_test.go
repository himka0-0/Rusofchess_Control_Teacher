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

func setupLectureTestDB() {
	config.InitTestDB()
	config.DB = config.TestDB

	// Создаём тестового пользователя
	config.TestDB.Create(&models.User{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
	})

	// Добавим три лекции
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

// Тест LecturePage (GET /lecture)
func TestLecturePage(t *testing.T) {
	setupLectureTestDB()
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	// Если нужно, загружаем шаблоны (иначе при рендере тоже будет panic)
	r.LoadHTMLGlob("../templates/*")

	// Эмулируем route: ставим пользователя в контекст и вызываем LecturePage
	r.GET("/lecture", func(c *gin.Context) {
		user := models.User{ID: 1, Email: "test@example.com"}
		c.Set("User", user)
		controllers.LecturePage(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/lecture", nil)
	r.ServeHTTP(w, req)

	// Ожидаем 200 OK
	assert.Equal(t, http.StatusOK, w.Code)
	// Если в lecture.html есть, например, слово "лекциями", проверяем его наличие:
	assert.Contains(t, w.Body.String(), "лекциями")
}

// Тест LectureHandler (POST /lecture) - случай, когда len(IDShnik) == len(input)
func TestLectureHandler_Update(t *testing.T) {
	setupLectureTestDB()
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	// Важно: эмулируем пользователя вручную
	r.POST("/lecture", func(c *gin.Context) {
		user := models.User{ID: 1, Email: "test@example.com"}
		c.Set("User", user)
		controllers.LectureHandler(c)
	})

	// Отправляем 3 лекции (столько же, сколько в БД), меняем их названия
	input := []models.PostLecture{
		{Number: 1, Lecture: "New Lecture A"},
		{Number: 2, Lecture: "New Lecture B"},
		{Number: 3, Lecture: "New Lecture C"},
	}
	jsonData, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/lecture", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Ожидаем 200 и "изменения сохранены"
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "изменения сохранены")

	// Проверяем, что названия в БД поменялись
	var lectures []models.Table_lecture
	config.TestDB.Order("lecture_person_id asc").Find(&lectures)
	assert.Len(t, lectures, 3)
	assert.Equal(t, "New Lecture A", lectures[0].Lecture)
	assert.Equal(t, "New Lecture B", lectures[1].Lecture)
	assert.Equal(t, "New Lecture C", lectures[2].Lecture)
}

// Тест LectureHandler (POST /lecture) - когда len(IDShnik) > len(input) (удаление)
func TestLectureHandler_Delete(t *testing.T) {
	setupLectureTestDB()
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.POST("/lecture", func(c *gin.Context) {
		user := models.User{ID: 1, Email: "test@example.com"}
		c.Set("User", user)
		controllers.LectureHandler(c)
	})

	// Отправим только 2 лекции, значит третья (Lecture C) должна удалиться
	input := []models.PostLecture{
		{Number: 1, Lecture: "Keep Lecture A"},
		{Number: 2, Lecture: "Keep Lecture B"},
	}
	jsonData, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/lecture", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Ожидаем 200
	assert.Equal(t, http.StatusOK, w.Code)
	// Твой код возвращает "измененение сохранено" (проверь точное написание)
	// но пусть будет так:
	assert.Contains(t, w.Body.String(), "измененение сохранено")

	// Теперь в БД должно остаться 2 лекции
	var lectures []models.Table_lecture
	config.TestDB.Order("lecture_person_id asc").Find(&lectures)
	assert.Len(t, lectures, 2)
	assert.Equal(t, "Keep Lecture A", lectures[0].Lecture)
	assert.Equal(t, "Keep Lecture B", lectures[1].Lecture)
}

// Тест LectureHandler (POST /lecture) - когда len(IDShnik) < len(input) (добавление)
func TestLectureHandler_Add(t *testing.T) {
	setupLectureTestDB()
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.POST("/lecture", func(c *gin.Context) {
		user := models.User{ID: 1, Email: "test@example.com"}
		c.Set("User", user)
		controllers.LectureHandler(c)
	})

	// Отправим 4 лекции, значит 4-ю нужно добавить
	input := []models.PostLecture{
		{Number: 1, Lecture: "Lecture A"},
		{Number: 2, Lecture: "Lecture B"},
		{Number: 3, Lecture: "Lecture C"},
		{Number: 4, Lecture: "New Lecture D"},
	}
	jsonData, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/lecture", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Ожидаем 200
	assert.Equal(t, http.StatusOK, w.Code)
	// Код возвращает "сохранено"
	assert.Contains(t, w.Body.String(), "сохранено")

	// Проверяем, что теперь 4 лекции
	var lectures []models.Table_lecture
	config.TestDB.Order("lecture_person_id asc").Find(&lectures)
	assert.Len(t, lectures, 4)
	assert.Equal(t, "New Lecture D", lectures[3].Lecture)
	assert.Equal(t, 4, lectures[3].Lecture_Person_id)
}
