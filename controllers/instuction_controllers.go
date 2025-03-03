package controllers

import (
	"awesomeProject1/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func InstructionPage(c *gin.Context) {
	userData, exists := c.Get("User")
	if !exists || userData == nil {
		log.Println("Пользователь не найден в middleware")
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Пользователь не авторизован"})
		return
	}
	user := userData.(models.User)
	c.HTML(http.StatusOK, "instruction.html", gin.H{"User": user})
}
