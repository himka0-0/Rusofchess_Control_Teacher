package utils

import (
	"fmt"
	"log"
	"net/smtp"
)

func SendVerificationEmail(email string, token string) {
	from := "maksim.grig02@mail.ru"
	password := "fFM8jjTUrE9uzgEM7G40"
	to := []string{email}
	smtpHost := "smtp.mail.ru"
	smtpPort := "587"

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
