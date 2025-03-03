package controllers

import (
	"awesomeProject1/config"
	"awesomeProject1/models"
	"awesomeProject1/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func RecoveryPasPage(c *gin.Context) {
	c.HTML(http.StatusOK, "recoveryPassword.html", gin.H{})
}
func RecoveryPasHandler(c *gin.Context) {
	var input models.PostRecovery
	err := c.ShouldBindJSON(&input)
	if err != nil {
		log.Println("ошибка парсинга почты при востановлении", err)
	}
	var user models.User
	if err = config.DB.Model(&models.User{}).Where("email=?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Инструкция отправлена на почту"})
		return
	}
	RecoveryToken := utils.GenerationToken()
	err = config.DB.Model(&models.User{}).Where("email=?", input.Email).Update("Verification_token", RecoveryToken).Error

	go utils.RecoveryPassword(user.Email, RecoveryToken)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Инструкция отправлена на почту"})
}
func RecMailPage(c *gin.Context) {
	c.HTML(http.StatusOK, "recoveryPasMail.html", nil)
}
func RecMailHandler(c *gin.Context) {
	var input models.PasswordRecovery
	err := c.ShouldBindJSON(&input)
	if err != nil {
		log.Println("Ошибка парсинга данных при востановлении пароля после почты", err)
	}
	var user models.User
	err = config.DB.Model(&models.User{}).Where("Verification_token=?", input.Token).First(&user).Error
	if err != nil {
		log.Println("Ошибка определения пользователя по токену при востановлении пароля", err)
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("проблемы хеширования", err)
	}

	err = config.DB.Model(&models.User{}).Where(" Verification_token=?", input.Token).Updates(map[string]interface{}{
		"Password":           string(hashPass),
		"Verification_token": "",
	}).Error
	if err != nil {
		log.Println("ошибка при изменеии пароля", err)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Пароль успешно изменен"})
}
