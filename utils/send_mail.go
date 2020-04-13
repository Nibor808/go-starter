package utils

import (
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
	"os"
)

func SendMail(subject string, toEmail string, content string) bool {
	supportEmail, sEExists := os.LookupEnv("SUPPORT_EMAIL")
	if !sEExists {
		log.Fatal("Cannot get Support email")
	}
	apiKey, kExists := os.LookupEnv("SENDGRID_API_KEY")
	if !kExists {
		log.Fatal("No api key for sendgrid")
	}

	to := mail.NewEmail("", toEmail)
	from := mail.NewEmail("ODF Support", supportEmail)

	message := mail.NewSingleEmail(from, subject, to, content, content)
	client := sendgrid.NewSendClient(apiKey)

	res, err := client.Send(message)
	if err != nil {
		fmt.Println("Failed to send email:", err)
		fmt.Println("CODE:", res.StatusCode)
		fmt.Println("BODY:", res.Body)
		fmt.Println("Headers:", res.Headers)

		return false
	}

	return true
}

