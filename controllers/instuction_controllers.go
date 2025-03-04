package controllers

import (
	customLogger "awesomeProject1/logger"
	"awesomeProject1/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func InstructionPage(c *gin.Context) {
	userData, exists := c.Get("User")
	if !exists || userData == nil {
		customLogger.Logger.Warn("Проблема в контексте", zap.String("error", "Пользователь не авторизован"))
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Пользователь не авторизован"})
		return
	}
	user := userData.(models.User)
	c.HTML(http.StatusOK, "instruction.html", gin.H{"User": user})
}
