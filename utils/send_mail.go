package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendMail(subject string, toEmail string, html string) bool {
	adminEmail, sEExists := os.LookupEnv("ADMIN_EMAIL")
	if !sEExists {
		log.Fatal("Cannot get ADMIN_EMAIL")
	}

	apiKey, kExists := os.LookupEnv("SENDGRID_API_KEY")
	if !kExists {
		log.Fatal("Cannot get SENDGRID_API_KEY")
	}

	to := mail.NewEmail("", toEmail)
	from := mail.NewEmail("Go Starter", adminEmail)

	message := mail.NewSingleEmail(from, subject, to, html, html)
	client := sendgrid.NewSendClient(apiKey)

	if res, err := client.Send(message); err != nil {
		fmt.Println("Failed to send email:", err)
		fmt.Println("CODE:", res.StatusCode)
		fmt.Println("BODY:", res.Body)
		fmt.Println("Headers:", res.Headers)

		return false
	}

	return true
}
