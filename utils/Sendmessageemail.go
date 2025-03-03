package utils

import (
	"fmt"
	"log"
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
		log.Println("Ошибка при отправке письма:", err)
	} else {
		log.Println("Письмо отправлено на:", email)
	}
}
