package utils

import (
	customLogger "awesomeProject1/logger"
	"fmt"
	"go.uber.org/zap"
	"net/smtp"
	"os"
)

func SendVerificationEmail(email string, token string) {
	from := os.Getenv("MAIL_LOGIN")
	password := os.Getenv("MAIL_PASSWORD")
	to := []string{email}
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	verificationLink := fmt.Sprintf("http://localhost:8080/verify-email?token=%s", token)
	message := []byte("Subject: Подтверждение регистрации\n\n" +
		"Перейдите по ссылке для подтверждения email: " + verificationLink)
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		customLogger.Logger.Error("Ошибка при отправке письмо", zap.Error(err))
	} else {
		customLogger.Logger.Info("Письмо с регистрацией отправлено")
	}
}
