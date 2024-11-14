package mailer

import (
	"os"

	"gopkg.in/gomail.v2"
)

// SendMailWithGomail sends an email using Gomail
func SendMailWithGomail(to, subject, body string) error {

	from := os.Getenv("EMAIL_SENDER")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := 587
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPassword)
	d.SSL = false

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
