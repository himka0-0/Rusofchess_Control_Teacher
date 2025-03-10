package controllers

import (
	"awesomeProject1/config"
	customLogger "awesomeProject1/logger"
	"awesomeProject1/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func TelbotPage(c *gin.Context) {
	userData, exists := c.Get("User")
	if !exists || userData == nil {
		customLogger.Logger.Warn("Проблема в контексте", zap.String("error", "Пользователь не авторизован"))
	}
	user := userData.(models.User)
	var students []models.Table_student
	err := config.DB.Model(&models.Table_student{}).Select("ID,Name_Student,Alert_payment,Alert_moduls").Where("User_id=?", user.ID).Find(&students).Error
	if err != nil {
		customLogger.Logger.Error("Ошибка обращения к бд, поиск студентов, стр телеграм бота", zap.Error(err))
	}
	var vhod bool
	err = config.DB.Model(&models.Table_telegram_bot{}).Select("Vhod").Where("User_id=?", user.ID).Scan(&vhod).Error
	if err != nil {
		customLogger.Logger.Error("Ошибка поиска информации о наличии метки входа,стр телеграм бота", zap.Error(err))
	}
	c.HTML(http.StatusOK, "telgrambot.html", gin.H{
		"students": students,
		"User":     user,
		"vhod":     vhod,
	})
}
func TelbotHandler(c *gin.Context) {
	userData, exists := c.Get("User")
	if !exists || userData == nil {
		customLogger.Logger.Warn("Проблема в контексте", zap.String("error", "Пользователь не авторизован"))
	}
	user := userData.(models.User)

	var input models.PostTelbot
	err := c.ShouldBindJSON(&input)
	if err != nil {
		customLogger.Logger.Error("Ошибка парсинга входящих данных, стр телеграм бота", zap.Error(err))
	}
	var conrolID models.Table_telegram_bot
	err = config.DB.Model(&models.Table_telegram_bot{}).Select("Telegram_id").Where("User_id=?", user.ID).Find(&conrolID).Error
	if err != nil {
		customLogger.Logger.Error("Ошибка обращения к бд, налчиие telegram id, стр телеграм бота", zap.Error(err))
	}

	if conrolID.Telegram_id == 0 {
		c.JSON(http.StatusOK, gin.H{"suc": true, "message": "Вы не настроили бота"})
	} else {
		err = config.DB.Model(&models.Table_telegram_bot{}).Where("User_id=?", user.ID).Update("Vhod", input.ModuleAllToggle).Error
		if err != nil {
			customLogger.Logger.Error("ошибка в обнавлении разрешения на отправку уведомлений, стр телеграм бота", zap.Error(err))
		}

		var outStudent_Alertpay []models.Table_student
		var outStudent_Alertmod []models.Table_student
		err = config.DB.Model(&models.Table_student{}).Select("ID,Alert_payment").Where("User_id=?", user.ID).Find(&outStudent_Alertpay).Error
		if err != nil {
			customLogger.Logger.Error("Ошибка обращения к бд, поиск студентов, стр телеграм бота", zap.Error(err))
		}
		err = config.DB.Model(&models.Table_student{}).Select("ID,Alert_moduls").Where("User_id=?", user.ID).Find(&outStudent_Alertmod).Error
		if err != nil {
			customLogger.Logger.Error("Ошибка обращения к бд, поиск студентов, стр телеграм бота", zap.Error(err))
		}
		ModulMap := make(map[int]bool)
		for _, modul := range outStudent_Alertmod {
			ModulMap[modul.ID] = modul.Alert_moduls
		}
		PaymentMap := make(map[int]bool)
		for _, Pay := range outStudent_Alertpay {
			PaymentMap[Pay.ID] = Pay.Alert_payment
		}
		for _, el := range input.Students {
			if ModulMap[el.ID] != el.Alert_moduls {
				err = config.DB.Model(&models.Table_student{}).Where("ID=?", el.ID).Update("Alert_moduls", el.Alert_moduls).Error
				if err != nil {
					customLogger.Logger.Error("Ошибка обращения к бд, ошибка в записи уведомлений модулей, стр телеграм бота", zap.Error(err))
				}
			}
		}
		for _, el := range input.Students {
			if PaymentMap[el.ID] != el.Alert_payment {
				err = config.DB.Model(&models.Table_student{}).Where("ID=?", el.ID).Update("Alert_payment", el.Alert_payment).Error
				if err != nil {
					customLogger.Logger.Error("Ошибка обращения к бд, ошибка в записи уведомлений оплаты, стр телеграм бота", zap.Error(err))
				}
			}
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Всё сохранено"})
	}
}
