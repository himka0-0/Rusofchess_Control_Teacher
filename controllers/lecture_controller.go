package controllers

import (
	"awesomeProject1/config"
	customLogger "awesomeProject1/logger"
	"awesomeProject1/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// lecture
func LecturePage(c *gin.Context) {
	userData, exists := c.Get("User")
	if !exists || userData == nil {
		customLogger.Logger.Warn("Проблема в контексте стр лекций", zap.String("error", "Пользователь не авторизован"))
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Пользователь не авторизован"})
		return
	}
	user, ok := userData.(models.User)
	if !ok {
		customLogger.Logger.Warn("ошибка в переводе из формата стр лекций", zap.String("error", "ошибка в переводе из формата"))
		return
	}
	UserID := user.ID
	var data []models.PrintLecture
	err := config.DB.Model(&models.Table_lecture{}).Select("Lecture_Person_id,Lecture").Where("User_id=?", UserID).Find(&data).Error
	if err != nil {
		customLogger.Logger.Error("Ошибка обращения к бд, стр лекции, инфа не вытащена", zap.Error(err))
	}
	c.HTML(http.StatusOK, "lecture.html", gin.H{
		"data": data,
		"User": user,
	})
} //стр управления лекциями
func LectureHandler(c *gin.Context) {
	var input []models.PostLecture
	err := c.ShouldBindJSON(&input)
	if err != nil {
		customLogger.Logger.Error("Ошибка парсинга входящих данных, стр лекции", zap.Error(err))
	}
	userData, exists := c.Get("User")
	if !exists || userData == nil {
		customLogger.Logger.Warn("Проблема в контексте стр лекций", zap.String("error", "Пользователь не авторизован"))
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Пользователь не авторизован"})
		return
	}
	user := userData.(models.User)
	var IDShnik []uint
	err = config.DB.Model(&models.Table_lecture{}).Select("id").Where("User_id=?", user.ID).Find(&IDShnik).Error
	if err != nil {
		customLogger.Logger.Error("Ошибка поиска в бд леций пользователя, стр лекции", zap.Error(err))
	}
	//только изменеие название или порядка
	if len(IDShnik) == len(input) {
		for idx, el := range IDShnik {
			Element_input := input[idx]
			Lecture_Element := Element_input.Lecture
			err = config.DB.Model(&models.Table_lecture{}).Where("id=?", el).Update("Lecture", Lecture_Element).Error
			if err != nil {
				customLogger.Logger.Error("Ошибка с обновлением в бд, стр лекции", zap.Error(err))
			}
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "massage": "изменения сохранены"})
	}
	//изменение порядка и(или) удаление
	if len(IDShnik) > len(input) {
		for idx, el := range IDShnik {
			if idx < len(input) {
				Element_input := input[idx]
				Lecture_Element := Element_input.Lecture
				err = config.DB.Model(&models.Table_lecture{}).Where("id=?", el).Updates(map[string]interface{}{
					"Lecture":           Lecture_Element,
					"Lecture_Person_id": idx + 1,
				}).Error
				if err != nil {
					customLogger.Logger.Error("Ошибка с обновлением в бд, стр лекции", zap.Error(err))
				}
			} else {
				err = config.DB.Where("id=?", el).Delete(&models.Table_lecture{}).Error
				if err != nil {
					customLogger.Logger.Error("Ошибка с удалением лекций в бд, стр лекции", zap.Error(err))
				}
			}
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "измененение сохранено"})
	}
	//изменение порядка и(или) добавление новой лекции
	if len(IDShnik) < len(input) {
		for idx, el := range input {
			if idx < len(IDShnik) {
				err = config.DB.Model(&models.Table_lecture{}).Where("id=?", IDShnik[idx]).Updates(map[string]interface{}{
					"Lecture":           el.Lecture,
					"Lecture_Person_id": idx + 1,
				}).Error
				if err != nil {
					customLogger.Logger.Error("Ошибка с обновлением в бд, стр лекции", zap.Error(err))
				}
			} else {
				err = config.DB.Create(&models.Table_lecture{Lecture: el.Lecture, User_id: user.ID, Lecture_Person_id: idx + 1}).Error
				if err != nil {
					customLogger.Logger.Error("Ошибка создания новых лекций в бд, стр лекции", zap.Error(err))
				}
			}
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "сохранено"})
	}
}
