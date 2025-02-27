package controllers

import (
	"awesomeProject1/config"
	"awesomeProject1/models"
	"awesomeProject1/telegram"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"os"
)

func NotelessonPage(c *gin.Context) {
	userData, exists := c.Get("User") // Берем пользователя из контекста
	if !exists || userData == nil {
		log.Println("Пользователь не найден в middleware")
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Пользователь не авторизован"})
		return
	}
	user := userData.(models.User)
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
	c.HTML(http.StatusOK, "notelesson.html", gin.H{
		"students": students,
		"User":     user,
	})

} //стр отметки урока
func NotelessonHandler(c *gin.Context) {
	var input models.PostModuls
	err := c.ShouldBindJSON(&input)
	if err != nil {
		log.Println(err)
	}
	var student models.Table_student
	if err = config.DB.Where("id=?", input.Student_id).First(&student).Error; err != nil {
		log.Println("Ученик не найден:", err)
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Ученик не найден"})
		return
	}
	switch input.Module {
	case "Теория":
		student.Theory += 1
		student.Practice = 0
		student.Tasks = 0
		if (student.Namber_lecture != 0) && (input.Lock_lecture == false) {
			student.Namber_lecture += 1
		}
		if student.Theory > 2 && student.Alert_moduls == true {
			message := "Вы провели 3 лекции подряд,не забывайте, что практика и задачи тоже важны!Ученик:"
			telegram.MessageBot(message, student.Name_Student, student.User_id)
		}
	case "Практика":
		student.Theory = 0
		student.Practice += 1
		student.Tasks = 0
		if student.Practice > 2 && student.Alert_moduls == true {
			message := "Вы провели 3 практики подряд,не забывайте, что теория и задачи тоже важны!Ученик:"
			telegram.MessageBot(message, student.Name_Student, student.User_id)
		}
	case "Задачи":
		student.Theory = 0
		student.Practice = 0
		student.Tasks += 1
		if student.Tasks > 2 && student.Alert_moduls == true {
			message := "Ученик решает задача 3 урока подряд,не забывайте, что теория и практика тоже важны!Ученик:"
			telegram.MessageBot(message, student.Name_Student, student.User_id)
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Неверный модуль"})
		return
	}
	student.Payment -= 1
	if student.Payment < 1 && student.Alert_payment == true {
		message := "На балансе ученика не достаточно средств, напомни об оплате за уроки!Ученик:"
		telegram.MessageBot(message, student.Name_Student, student.User_id)
	}
	if err = config.DB.Save(&student).Error; err != nil {
		log.Println("Ошибка при обновлении данных ученика:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Ошибка при сохранении"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "massage": "Сохранено"})
}
