package utils

import (
	"crypto/tls"
	"fmt"
	"os"
	"regexp"

	"gopkg.in/gomail.v2"
)

// type EmailConfig struct {
// 	SMTPHost     string
// 	SMTPPort     int
// 	SMTPUsername string
// 	SMTPPassword string
// }

// var EmailSettings = EmailConfig{
// 	SMTPHost:     os.Getenv("EmailHost"),
// 	SMTPPort:     587,
// 	SMTPUsername: os.Getenv("EmailSender"),
// 	SMTPPassword: os.Getenv("EmailPassword"),
// }

func IsValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

func SendVerificationEmail(email, name, verificationLink string) error {
	SMTPHost := os.Getenv("EmailHost")
	SMTPPort := 587
	SMTPUsername := os.Getenv("EmailSender")
	SMTPPassword := os.Getenv("EmailPassword")
	m := gomail.NewMessage()
	m.SetHeader("From", SMTPUsername)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Please verify your email")

	emailBody := fmt.Sprintf(`
		Hello %s,

		Thank you for signing up! Please verify your email by clicking the link below:

		%s

		This link will expire in 24 hours.

		If you didn't create an account, please ignore this email.

		Best regards,
		%s Team
	`, name, verificationLink, SMTPUsername)

	m.SetBody("text/html", emailBody)

	d := gomail.NewDialer(
		SMTPHost,
		SMTPPort,
		SMTPUsername,
		SMTPPassword,
	)
	fmt.Println(SMTPHost, SMTPPort, SMTPUsername, SMTPPassword, email)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Error sending email: ", err)
		return err
	}
	return nil
}
