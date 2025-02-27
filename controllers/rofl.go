package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// rofl
func WePage(c *gin.Context) {
	c.HTML(http.StatusOK, "We.html", nil)
} //моя страничка
