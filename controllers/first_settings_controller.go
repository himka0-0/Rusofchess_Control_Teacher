package controllers

import (
	"awesomeProject1/config"
	"awesomeProject1/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
)

func FirstSettinPage(c *gin.Context) {
	userData, exists := c.Get("User") // Берем пользователя из контекста
	if !exists || userData == nil {
		log.Println("Пользователь не найден в middleware")
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Пользователь не авторизован"})
		return
	}
	user := userData.(models.User)
	var significationStruct models.Table_telegram_bot
	err := config.DB.Model(&models.Table_telegram_bot{}).Select("Hash").Where("User_id=?", user.ID).Scan(&significationStruct).Error
	if err != nil {
		log.Println("Не получается вытащить хеш из таблицы", err)
	}
	var signification string
	signification = significationStruct.Hash
	c.HTML(http.StatusOK, "firstSetting.html", gin.H{
		"User":          user,
		"signification": signification,
	})
} //стр первой настройки
func FirstSettingHandler(c *gin.Context) {
	var input models.PostSettings
	tokenstr, _ := c.Cookie("token")
	claims := jwt.MapClaims{}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET не найден в .env")
	}
	_, _ = jwt.ParseWithClaims(tokenstr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	emailFromToken, _ := claims["email"].(string)
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("ошибка парсинга начальной настройки")
	}
	//запись учеников
	if input.Marking == "1" {
		var userID uint
		err := config.DB.Model(&models.User{}).Select("id").Where("email = ?", emailFromToken).Scan(&userID).Error
		if err != nil {
			fmt.Println("Ошибка:", err)
		}
		config.DB.Create(&models.Table_student{User_id: userID, Name_Student: input.Meaning, Alert_payment: true, Alert_moduls: true})
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Ученик сохранен"})
	}
	//запись лекций
	if input.Marking == "0" {
		var lection models.User
		var num_lectures_int int
		err := config.DB.Where("email =?", emailFromToken).First(&lection).Error
		if err != nil {
			log.Println("ошибка при определении количества лекций")
		}
		lection.Lectures_introduced += 1
		num_lectures_int = lection.Lectures_introduced
		if err = config.DB.Save(&lection).Error; err != nil {
			log.Println("Ошибка при обновлении данных ученика:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Ошибка при сохранении"})
			return
		}
		var UserID uint
		err = config.DB.Model(&models.User{}).Select("id").Where("email =?", emailFromToken).Scan(&UserID).Error
		if err != nil {
			log.Println("ошибка при определениии UserID", err)
		}
		config.DB.Create(&models.Table_lecture{User_id: UserID, Lecture: input.Meaning, Lecture_Person_id: num_lectures_int})
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Лекция сохранена"})
	}
}
