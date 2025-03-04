package controllers

import (
	"awesomeProject1/config"
	customLogger "awesomeProject1/logger"
	"awesomeProject1/models"
	"awesomeProject1/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func RecoveryPasPage(c *gin.Context) {
	c.HTML(http.StatusOK, "recoveryPassword.html", gin.H{})
}
func RecoveryPasHandler(c *gin.Context) {
	var input models.PostRecovery
	err := c.ShouldBindJSON(&input)
	if err != nil {
		customLogger.Logger.Error("Ошибка парсинга почты при востановлении, стр востановления пароля", zap.Error(err))
	}
	var user models.User
	if err = config.DB.Model(&models.User{}).Where("email=?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Инструкция отправлена на почту"})
		return
	}
	RecoveryToken := utils.GenerationToken()
	err = config.DB.Model(&models.User{}).Where("email=?", input.Email).Update("Verification_token", RecoveryToken).Error
	if err != nil {
		customLogger.Logger.Error("Ошибка обнавления токена, стр востановления пароля", zap.Error(err))
	}
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
		customLogger.Logger.Error("Ошибка парсинга данных при востановлении пароля после почты, стр востановления пароля", zap.Error(err))
	}
	var user models.User
	err = config.DB.Model(&models.User{}).Where("Verification_token=?", input.Token).First(&user).Error
	if err != nil {
		customLogger.Logger.Error("Ошибка определения пользователя по токену при востановлении пароля, стр востановления пароля", zap.Error(err))
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
		customLogger.Logger.Error("ошибка при изменеии пароля, стр востановления пароля", zap.Error(err))
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Пароль успешно изменен"})
}
