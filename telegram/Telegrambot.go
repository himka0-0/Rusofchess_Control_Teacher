package telegram

import (
	"awesomeProject1/config"
	customLogger "awesomeProject1/logger"
	"awesomeProject1/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"os"
)

var bot *tgbotapi.BotAPI

func RunBot() {
	var err error
	botenv := os.Getenv("TELEGRAM_TOKEN")
	bot, err = tgbotapi.NewBotAPI(botenv)
	if err != nil {
		customLogger.Logger.Fatal("телеграм бот лег")
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		username := update.Message.From.UserName
		if username == "" {
			username = update.Message.From.FirstName
		}

		messageText := update.Message.Text
		chatID := update.Message.Chat.ID

		msg := tgbotapi.NewMessage(chatID, "")

		switch messageText {
		case "/start":
			msg.Text = "Добрый день!Введите значение с сайта"
		default:
			if Validation_hash(messageText) == true {
				SaveUsers(username, messageText, chatID)
				msg.Text = "Добро пожаловать!" + username
			} else {
				msg.Text = "Ваше значение не подходит,сначала зарегистрируйтесь на сайте"
			}
		}
		bot.Send(msg)
	}
}

func Validation_hash(text string) bool {
	var hashs []models.Table_telegram_bot
	err := config.DB.Model(&models.Table_telegram_bot{}).Select("User_id,hash").Scan(&hashs).Error
	if err != nil {
		customLogger.Logger.Error("ошибка обращения к бд, вытаскивание хеша, функция телеграм бота", zap.Error(err))
	}
	var presence_hash int
	var presence bool
	for _, el := range hashs {
		if el.Hash == text {
			presence_hash += 1
		}
	}
	if presence_hash == 1 {
		presence = true
	} else {
		presence = false
	}
	return presence
}

func SaveUsers(username string, messageText string, chatID int64) {
	err := config.DB.Model(&models.Table_telegram_bot{}).Where("hash=?", messageText).Updates(&models.Table_telegram_bot{First_name: username, Telegram_id: chatID, Vhod: true}).Error
	if err != nil {
		customLogger.Logger.Error("ошибка обращения к бд, не обнавился препод, функция телеграм бота", zap.Error(err))
	}
}

func MessageBot(message string, nameStudent string, IdTeacher uint) {
	var onoff bool
	err := config.DB.Model(&models.Table_telegram_bot{}).Select("vhod").Where("User_id=?", IdTeacher).Find(&onoff).Error
	if err != nil {
		customLogger.Logger.Error("проблема обращения к бд, определение выключателя, функция телеграмм бота", zap.Error(err))
	}
	if onoff == true {
		var dialogue int64
		err = config.DB.Model(&models.Table_telegram_bot{}).Select("telegram_id").Where("User_id=?", IdTeacher).Find(&dialogue).Error
		if err != nil {
			customLogger.Logger.Error("проблема обращения к бд, определение выключателя, функция телеграмм бота", zap.Error(err))
		}
		message = message + " " + nameStudent
		msg := tgbotapi.NewMessage(dialogue, message)
		bot.Send(msg)
	}
}
