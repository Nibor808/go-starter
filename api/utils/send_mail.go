package utils

import (
	"fmt"
	"net/smtp"
)

// MailArgs is ...
type MailArgs struct {
	AdminEmail string
	Subject    string
	To         string
	Body       string
	MailPass   string
	MailHost   string
}

// SendMail sends an email and returns a error
func SendMail(args MailArgs) error {
	message := MessageConfig(args)
	mailHost := args.MailHost
	mailPass := args.MailPass
	adminEmail := args.AdminEmail

	auth := smtp.PlainAuth("", adminEmail, mailPass, mailHost)
	to := []string{args.To}

	if err := smtp.SendMail(mailHost+":587", auth, adminEmail, to, message); err != nil {
		return fmt.Errorf("unable to send email: %w", err)
	}

	return nil
}

// MessageConfig sets up the message
func MessageConfig(args MailArgs) []byte {
	return []byte("To: " + args.To + "\r\n" +
		"Subject: " + args.Subject + "\r\n" +
		"From: " + args.AdminEmail + "\r\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		args.Body +
		"\r\n")
}
