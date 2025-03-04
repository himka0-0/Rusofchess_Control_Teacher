package controllers

import (
	"awesomeProject1/config"
	customLogger "awesomeProject1/logger"
	"awesomeProject1/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func ResultPage(c *gin.Context) {
	userData, exists := c.Get("User") // Берем пользователя из контекста
	if !exists || userData == nil {
		customLogger.Logger.Warn("Проблема в контексте", zap.String("error", "Пользователь не авторизован"))
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Пользователь не авторизован"})
		return
	}
	user := userData.(models.User)
	emailData, exists := c.Get("email") // Берем пользователя из контекста
	if !exists || emailData == nil {
		customLogger.Logger.Warn("Проблема в контексте", zap.String("error", "Пользователь не авторизован"))
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Пользователь не авторизован"})
		return
	}
	email := emailData.(string)
	var userID uint
	err := config.DB.Model(&models.User{}).Select("ID").Where("email=?", email).Scan(&userID).Error
	if err != nil {
		customLogger.Logger.Error("Не смог определить id, стр вывода всего", zap.Error(err))
	}
	var output []models.Table_student
	err = config.DB.Model(&models.Table_student{}).Select("ID,Name_Student,Payment,Theory,Practice,Tasks,Namber_lecture").Where("User_id = ?", userID).Find(&output).Error
	if err != nil {
		customLogger.Logger.Error("Ошибка определения в бд, стр вывести всё", zap.Error(err))
	}
	var lecture []models.Table_lecture
	err = config.DB.Model(&models.Table_lecture{}).Select("Lecture_Person_id,Lecture").Where("User_id=?", userID).Find(&lecture).Error
	lectureMap := make(map[int]string)
	for _, lec := range lecture {
		lectureMap[lec.Lecture_Person_id] = lec.Lecture
	}
	var outputtofront []models.Resultstruct
	for _, el := range output {
		if el.Namber_lecture == 0 {
			outputtofront = append(outputtofront, models.Resultstruct{ID: el.ID, Name_Student: el.Name_Student, Payment: el.Payment, Lecture: "Лекция не выбрана", Theory: el.Theory, Practice: el.Practice, Tasks: el.Tasks})
		} else {
			value, _ := lectureMap[el.Namber_lecture]
			outputtofront = append(outputtofront, models.Resultstruct{ID: el.ID, Name_Student: el.Name_Student, Payment: el.Payment, Lecture: value, Theory: el.Theory, Practice: el.Practice, Tasks: el.Tasks})
		}
	}

	c.HTML(http.StatusOK, "result.html", gin.H{
		"outputtofront": outputtofront,
		"User":          user,
	})
} //вывод таблицы со всеми значениями
