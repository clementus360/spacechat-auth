package services

import (
	"fmt"
	"net/smtp"
	"os"
)

func InitializeSMTP() (string, smtp.Auth, string) {
	SMTP_EMAIL := os.Getenv("SMTP_EMAIL")
	SMTP_PASSWORD := os.Getenv("SMTP_PASSWORD")

	smtpHost := "smtp.mailgun.org"
	smtpPort := "587"

	fmt.Println(SMTP_EMAIL, " : ", SMTP_PASSWORD)
	auth := smtp.PlainAuth("", SMTP_EMAIL, SMTP_PASSWORD, smtpHost)

	connStr := smtpHost + ":" + smtpPort

	return connStr, auth, SMTP_EMAIL

}

func SendEmail(receiver string, message string) error {
	connStr, auth, from := InitializeSMTP()

	to := []string{
		receiver,
	}

	fmt.Println(receiver)

	// Sending email.
	err := smtp.SendMail(connStr, auth, from, to, []byte("Your spacechat otp is: "+message))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
