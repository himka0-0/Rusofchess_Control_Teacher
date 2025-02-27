package controllers

import (
	"awesomeProject1/config"
	"awesomeProject1/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
)

func PaymentstudentPage(c *gin.Context) {
	userData, exists := c.Get("User") // Берем пользователя из контекста
	if !exists || userData == nil {
		log.Println("Пользователь не найден в middleware")
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Пользователь не авторизован"})
		return
	}
	user := userData.(models.User)
	tokenstr, _ := c.Cookie("token")
	claims := jwt.MapClaims{}
	jwtSecret := os.Getenv("JWT_SECRET") // Загружаем секретный ключ из .env
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET не найден в .env")
	}
	_, _ = jwt.ParseWithClaims(tokenstr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil // Возвращаем секретный ключ как []byte
	})
	emailFromToken, _ := claims["email"].(string)
	var userID uint
	err := config.DB.Model(&models.User{}).Select("id").Where("email = ?", emailFromToken).Scan(&userID).Error
	if err != nil {
		log.Println("не могу наути пользователя")
		return
	}
	var students []models.Table_student
	err = config.DB.Model(&models.Table_student{}).Select("ID,Name_Student").Where("User_id = ?", userID).Scan(&students).Error
	if err != nil {
		log.Println(err)
	}
	c.HTML(http.StatusOK, "paymentstudent.html", gin.H{
		"students": students,
		"User":     user,
	})
} //стр записи об оплате учеников
func PaymentstudentHandler(c *gin.Context) {
	var input models.Paymentstudent
	err := c.ShouldBindJSON(&input)
	if err != nil {
		log.Println(err)
	}
	var student models.Table_student
	if err = config.DB.Where("id = ?", input.ID).First(&student).Error; err != nil {
		log.Println("Ученик не найден:", err)
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Ученик не найден"})
		return
	}
	student.Payment += input.Payment
	if err = config.DB.Save(&student).Error; err != nil {
		log.Println("Ошибка при обновлении данных ученика:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Ошибка при сохранении"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Оплата сохранена"})
}
