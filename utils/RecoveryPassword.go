package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

func RecoveryPassword(email string, RecoveryToken string) {
	from := os.Getenv("MAIL_LOGIN")
	password := os.Getenv("MAIL_PASSWORD")
	to := []string{email}
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	RecoveryLink := fmt.Sprintf("http://localhost:8080/recovery-password?token=%s", RecoveryToken)
	message := []byte("Subject: Востановление пароля\n\n" + "Перейдите по ссылке для изменения пароля: " + RecoveryLink)
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Println("Ошибка при отправке письма:", err)
	} else {
		log.Println("Письмо отправлено на:", email)
	}
}
