package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"gopkg.in/gomail.v2"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	} else {
		fmt.Println(".env file loaded successfully.")
	}

	secret = os.Getenv("SECRET_KEY")

	if secret == "" {
		fmt.Println("SECRET_KEY is not set! Check your .env file.")
	} else {
		fmt.Println("SECRET_KEY loaded successfully:", secret)
	}
}

var sender string = os.Getenv("SENDER_EMAIL")
var receiver string = os.Getenv(""RECEIDERa)
func envite() {
	m := gomail.NewMessage()
	m.SetHeader("From", "your-email@gmail.com")
	m.SetHeader("To", "recipient@example.com")
	m.SetHeader("Subject", "Hello from Golang!")
	m.SetBody("text/plain", "This is a test email from Go.")

	d := gomail.NewDialer("smtp.gmail.com", 587, "your-email@gmail.com", "your-email-password")

	if err := d.DialAndSend(m); err != nil {
		log.Fatal(err)
	}

	log.Println("Email sent successfully!")
}
