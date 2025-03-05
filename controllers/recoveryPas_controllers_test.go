package controllers_test

import (
	"awesomeProject1/config"
	"awesomeProject1/controllers"
	"awesomeProject1/models"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupRecoveryPasTest() {
	config.InitTestDB()
	config.DB = config.TestDB

	config.TestDB.Create(&models.User{
		ID:             1,
		Name:           "Test User",
		Email:          "test@example.com",
		Password:       "123456789",
		Email_verified: true,
	})
	config.TestDB.Create(&models.User{
		ID:                 2,
		Name:               "Test User2",
		Email:              "test2@example.com",
		Password:           "123456789",
		Email_verified:     true,
		Verification_token: "valid-token",
	})
}

func TestRecoveryPasHandler(t *testing.T) {
	setupRecoveryPasTest()
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.LoadHTMLGlob("../templates/*")

	r.POST("/recoveryPassword", func(c *gin.Context) {
		user := models.User{ID: 1, Email: "test@example.com"}
		c.Set("User", user)
		controllers.RecoveryPasHandler(c)
	})

	input := models.PostRecovery{
		"test@example.com",
	}
	jsonData, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/recoveryPassword", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Инструкция")

	var user models.User
	config.TestDB.First(&user, 1)
	assert.NotEmpty(t, user.Verification_token, "Verification token should not be empty")
}

func TestRecMailHandler(t *testing.T) {
	setupRecoveryPasTest()
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.LoadHTMLGlob("../templates/*")

	r.POST("/recovery-password", func(c *gin.Context) {
		controllers.RecMailHandler(c)
	})

	input := models.PasswordRecovery{
		Password: "987654321",
		Token:    "valid-token",
	}
	jsonData, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/recovery-password", bytes.NewBuffer(jsonData))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Пароль успешно изменен")

	var user models.User
	config.TestDB.First(&user, 2)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("987654321"))
	assert.NoError(t, err, "Пароль должен быть успешно изменен и соответствовать хешу")
	assert.Equal(t, user.Verification_token, "")
}
