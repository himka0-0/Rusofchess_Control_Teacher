package controllers

import (
	"awesomeProject1/config"
	"awesomeProject1/models"
	"awesomeProject1/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"time"
)

// aut
func FirstPage(c *gin.Context) {
	c.HTML(http.StatusOK, "nachalo.html", nil)
} //первая страница

func RegPage(c *gin.Context) {
	c.HTML(http.StatusOK, "registration.html", nil)
} //стр регистрации

func RegHandler(c *gin.Context) {
	var user models.User
	if er := c.ShouldBindJSON(&user); er != nil {
		fmt.Println("ошибка парсинга при реге", er)
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("проблемы хеширования", err)
	}
	user.Password = string(hashPass)

	verifyToken := utils.GenerationToken()
	user.Verification_token = verifyToken

	err = config.DB.Create(&user).Error
	if err != nil {
		log.Println("Ошибка создания пользователя в бд", err)
	}
	go func() {
		hash := utils.HashIDAndEmail(user.ID, user.Email)
		err = config.DB.Create(&models.Table_telegram_bot{User_id: user.ID, Hash: hash, Vhod: false}).Error
		if err != nil {
			log.Println("Ошибка создания пользователя в таблице телеграм бота", err)
		}
	}()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Регистрация успешна!Нобходимо подтвердить почту",
	})
	go utils.SendVerificationEmail(user.Email, verifyToken)
}

func AutPage(c *gin.Context) {
	c.HTML(http.StatusOK, "authentication.html", nil)
} //стр аутентификации

func AutHandler(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("Ошибка парсинга", err)
	}
	var user models.User
	if er := config.DB.Where("email=?", input.Email).First(&user).Error; er != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
		return
	}

	if !user.Email_verified {
		c.JSON(http.StatusForbidden, gin.H{"error": "Подтвердите email перед входом!"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
		return
	}
	token := GenerateJwt(user.Email)
	c.SetCookie("token", token, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"redirect": "/kabinet",
	})
}

// expectationverify
func VerifyEmailPage(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Токен отсутствует"})
		return
	}
	var user models.User
	if err := config.DB.Where("verification_token=?", token).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недействительный токен"})
		return
	}

	user.Email_verified = true
	user.Verification_token = ""
	config.DB.Save(&user)
	c.HTML(http.StatusOK, "verifyEmail.html", gin.H{"message": "Email успешно подтвержден! Теперь вы можете войти."})
} //стр подтверждение почты
func VerifyPage(c *gin.Context) {
	c.HTML(http.StatusOK, "verify.html", nil)
} //

func GetProfile(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Ошибка авторизации"})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	c.HTML(http.StatusOK, "kabinet.html", gin.H{"User": user})
} //вывод кабинета

func LogoutHandler(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Вы успешно вышли из системы", "redirect": "/"})
} //удаление куки

func GenerateJwt(email string) string {
	expirationTime := time.Now().Add(24 * time.Hour) // Токен на 24 часа
	claims := &models.Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Println("ошибка в GenerateJwt", err)
	}
	return tokenString
}
