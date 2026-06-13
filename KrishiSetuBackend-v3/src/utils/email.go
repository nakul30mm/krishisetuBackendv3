package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendOTPEmail(toEmail, otp string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	if from == "" || password == "" {
		return fmt.Errorf("SMTP credentials not set")
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	subject := "Subject: KrishiSetu Password Reset OTP\n"
	body := fmt.Sprintf(
		"Your OTP for resetting password is: %s\n\nThis OTP is valid for a limited time.\n\nIf you did not request this, ignore this email.",
		otp,
	)

	message := []byte(subject + "\n" + body)

	return smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from,
		[]string{toEmail},
		message,
	)
}
