package middlewares

import (
	"awesomeProject1/config"
	customLogger "awesomeProject1/logger"
	"awesomeProject1/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/authentication")
			c.Abort()
			return
		}

		claims := jwt.MapClaims{}
		jwtSecret := os.Getenv("JWT_SECRET") // Загружаем секретный ключ из .env
		if jwtSecret == "" {
			customLogger.Logger.Fatal("JWT_SECRET не найден в .env")
		}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.Redirect(http.StatusSeeOther, "/authentication")
			c.Abort()
			return
		}
		emailFromToken, ok := claims["email"].(string)
		if !ok {
			c.Redirect(http.StatusSeeOther, "/authentication")
			c.Abort()
			return
		}
		var user models.User
		err = config.DB.Where("email = ?", emailFromToken).First(&user).Error
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/authentication")
			c.Abort()
			return
		}
		c.Set("email", emailFromToken)
		c.Set("User", user)
		c.Next()
	}
}
