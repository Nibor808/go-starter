package utils

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// MailArgs is ...
type MailArgs struct {
	AdminEmail string
	APIKey     string
	Subject    string
	To         string
	HTML       string
}

// SendMail sends an email and returns a error
func SendMail(args MailArgs) error {
	message, err := MessageConfig(args)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	client := sendgrid.NewSendClient(args.APIKey)

	if res, err := client.Send(message); err != nil {
		fmt.Println("Failed to send email:", err)
		fmt.Println("CODE:", res.StatusCode)
		fmt.Println("BODY:", res.Body)
		fmt.Println("Headers:", res.Headers)

		return fmt.Errorf("Unable to send email: %w", err)
	}

	return nil
}

// MessageConfig sets up the message
func MessageConfig(args MailArgs) (*mail.SGMailV3, error) {
	to := mail.NewEmail("", args.To)
	from := mail.NewEmail(args.Subject, args.AdminEmail)
	message := mail.NewSingleEmail(from, args.Subject, to, args.HTML, args.HTML)

	if message.Subject == args.Subject && len(message.Content) == 2 {
		return message, nil
	}

	return nil, fmt.Errorf("Unable to configure message")
}
