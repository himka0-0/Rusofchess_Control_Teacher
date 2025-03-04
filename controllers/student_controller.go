package controllers

import (
	"awesomeProject1/config"
	customLogger "awesomeProject1/logger"
	"awesomeProject1/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func StudentPage(c *gin.Context) {
	userData, exists := c.Get("User")
	if !exists || userData == nil {
		customLogger.Logger.Warn("Проблема в контексте", zap.String("error", "Пользователь не авторизован"))
	}
	user := userData.(models.User)
	var lecture []models.Table_lecture
	err := config.DB.Model(&models.Table_lecture{}).Select("Lecture_Person_id,Lecture").Where("User_id=?", user.ID).Find(&lecture).Error
	if err != nil {
		customLogger.Logger.Error("Ошибка обращения к бд, поиск лекций, стр телеграм бота", zap.Error(err))
	}
	var students []models.Table_student
	err = config.DB.Model(&models.Table_student{}).Select("ID,Name_Student,Namber_lecture").Where("User_id=?", user.ID).Find(&students).Error
	if err != nil {
		customLogger.Logger.Error("Ошибка обращения к бд, Ошибка при вытаскивании студента, стр телеграм бота", zap.Error(err))
	}

	c.HTML(http.StatusOK, "student.html", gin.H{
		"students": students,
		"User":     user,
		"lecture":  lecture,
	})
} //стр управления учениками

func StudentHandler(c *gin.Context) {
	userData, exists := c.Get("User")
	if !exists || userData == nil {
		customLogger.Logger.Warn("Проблема в контексте", zap.String("error", "Пользователь не авторизован"))
	}
	user, ok := userData.(models.User)
	if !ok {
		customLogger.Logger.Warn("Ошибка приведения userData к User,стр студентов", zap.String("error", "Пользователь не авторизован"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка авторизации"})
		return
	}
	var students []models.Table_student
	err := config.DB.Model(&models.Table_student{}).Select("ID,Name_Student,Namber_lecture").Where("User_id=?", user.ID).Find(&students).Error
	if err != nil {
		customLogger.Logger.Error("Ошибка обращения к бд, поиск студентов, стр студенты", zap.Error(err))
	}
	var input []models.PostStudent
	if err = c.ShouldBindJSON(&input); err != nil {
		customLogger.Logger.Error("Ошибка парсинга входящих данных, стр студенты", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные входные данные"})
		return
	}
	studentMap := make(map[int]models.Table_student)
	for _, student := range students {
		studentMap[student.ID] = student
	}
	var updatedInput []models.PostStudent
	var newStudentIDs []int

	for _, el := range input {
		if el.ID == 0 {
			newStudent := models.Table_student{User_id: user.ID, Name_Student: el.Name, Namber_lecture: el.Lecture, Alert_payment: true, Alert_moduls: true}
			err = config.DB.Create(&newStudent).Error
			if err != nil {
				customLogger.Logger.Error("Ошибка сохранения студента, стр студенты", zap.Error(err))
			} else {
				updatedInput = append(updatedInput, models.PostStudent{ID: newStudent.ID, Name: newStudent.Name_Student, Lecture: newStudent.Namber_lecture})
				newStudentIDs = append(newStudentIDs, newStudent.ID)
			}
		} else {
			updatedInput = append(updatedInput, el)
		}
	}

	var existingStudentIDs []int
	for _, student := range students {
		existingStudentIDs = append(existingStudentIDs, student.ID)
	}

	var inputStudentIDs []int
	for _, el := range updatedInput {
		inputStudentIDs = append(inputStudentIDs, el.ID)
	}

	var studentsToDelete []int
	for _, id := range existingStudentIDs {
		if !contains(inputStudentIDs, id) {
			studentsToDelete = append(studentsToDelete, id)
		}
	}

	if len(studentsToDelete) > 0 {
		err = config.DB.Where("user_id = ? AND id IN (?)", user.ID, studentsToDelete).Delete(&models.Table_student{}).Error
		if err != nil {
			customLogger.Logger.Error("Ошибка удаления ученика, стр студенты", zap.Error(err))
		}
	}

	for _, el := range updatedInput {
		if existing, found := studentMap[el.ID]; found {
			if el.Lecture != existing.Namber_lecture {
				if err = config.DB.Model(&models.Table_student{}).Where("ID = ?", el.ID).UpdateColumn("Namber_lecture", el.Lecture).Error; err != nil {
					customLogger.Logger.Error("Ошибка при обновлении лекций ученика, стр студенты", zap.Error(err))
				}
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Всё сохранено"})
}

func contains(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
