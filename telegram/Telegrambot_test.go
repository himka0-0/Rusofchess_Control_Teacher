package telegram_test

import (
	"awesomeProject1/config"
	"awesomeProject1/models"
	"awesomeProject1/telegram"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupTelegrambotTestDB() {
	config.InitTestDB()
	config.DB = config.TestDB

	config.TestDB.Create(&models.User{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
	})

	config.TestDB.Create(&models.Table_telegram_bot{
		User_id:     1,
		Hash:        "test_hash",
		First_name:  "ahimkas",
		Telegram_id: 123456789,
		Vhod:        true,
	})

	config.TestDB.Create(&models.Table_telegram_bot{
		User_id: 2,
		Hash:    "test_hashtwo",
	})

}
func TestValidation_hash(t *testing.T) {
	setupTelegrambotTestDB()
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name string
		hash string
		want bool
	}{
		{name: "Valid hash",
			hash: "test_hash",
			want: true,
		},
		{name: "Invalid hash",
			hash: "",
			want: false},
		{name: "Invalid hash",
			hash: "false_hash",
			want: false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := telegram.Validation_hash(tc.hash)
			if got != tc.want {
				t.Errorf("Validation_hash(%q) = %v; want %v", tc.hash, got, tc.want)
			}
		})
	}
}

func TestSaveUsers(t *testing.T) {
	setupTelegrambotTestDB()
	gin.SetMode(gin.TestMode)

	tests := []struct {
		username    string
		messageText string
		chatID      int64
	}{
		{
			username:    "testusername",
			messageText: "test_hashtwo",
			chatID:      123456782,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.username, func(t *testing.T) {
			telegram.SaveUsers(tc.username, tc.messageText, tc.chatID)
		})
	}
	var user models.Table_telegram_bot
	config.TestDB.First(&user, 2)
	assert.Equal(t, int64(123456782), user.Telegram_id)
	assert.Equal(t, "testusername", user.First_name)
	assert.Equal(t, true, user.Vhod)
}
