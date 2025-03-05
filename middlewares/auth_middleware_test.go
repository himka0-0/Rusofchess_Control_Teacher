package middlewares_test

import (
	"awesomeProject1/config"
	"awesomeProject1/middlewares"
	"awesomeProject1/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func setupTestEnvironment() {
	// Устанавливаем JWT_SECRET для тестов
	os.Setenv("JWT_SECRET", "testsecret")
	// Инициализация тестовой базы данных
	config.InitTestDB()
	config.DB = config.TestDB
}

func createTestToken(email string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return tokenString
}

func TestAuthMiddleware(t *testing.T) {
	setupTestEnvironment()

	// Создаем тестового пользователя
	user := models.User{Email: "test@example.com"}
	config.TestDB.Create(&user)

	tests := []struct {
		name        string
		token       string
		expected    int
		description string
	}{
		{
			name:        "Valid token",
			token:       createTestToken("test@example.com"),
			expected:    http.StatusOK,
			description: "Должен пропустить запрос с валидным токеном",
		},
		{
			name:        "No token",
			token:       "",
			expected:    http.StatusSeeOther,
			description: "Должен перенаправить, если токен отсутствует",
		},
		{
			name:        "Invalid token",
			token:       "invalid.token.here",
			expected:    http.StatusSeeOther,
			description: "Должен перенаправить, если токен невалидный",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем тестовый контекст Gin
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req, _ := http.NewRequest("GET", "/", nil)
			c.Request = req

			// Устанавливаем токен в куки, если он есть
			if tt.token != "" {
				c.Request.AddCookie(&http.Cookie{Name: "token", Value: tt.token})
			}

			// Вызываем middleware
			middlewares.AuthMiddleware()(c)

			// Проверяем статус код
			if w.Code != tt.expected {
				t.Errorf("%s: ожидался статус код %d, получили %d", tt.description, tt.expected, w.Code)
			}
		})
	}

	// Очистка базы данных после тестов
	config.TestDB.Migrator().DropTable(&models.User{})
}
