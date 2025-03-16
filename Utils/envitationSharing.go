package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

var (
	sender    string
	formLink  string
	emailPass string
	secretKey string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	} else {
		// log.Println(".env file loaded successfully.")
	}

	secretKey = os.Getenv("SECRET_KEY")
	sender = os.Getenv("ADMIN_EMAIL")
	formLink = os.Getenv("LINK")
	emailPass = os.Getenv("ADMIN_EMAIL_PASS")

	if secretKey == "" {
		log.Println("Warning: SECRET_KEY is not set! Check your .env file.")
	}
	if formLink == "" || sender == "" || emailPass == "" {
		log.Fatal("Error: Missing required environment variables. Please check your .env file.")
	}
}

func Invite(receiver string, vendorFirstName string, vendorOtherName string, token string) bool {
	link := fmt.Sprintf("%s?token=%s", formLink, token)
	emailBody := fmt.Sprintf(
		"Dear %s %s,\n\n"+
			"Thank you for reaching out regarding access to Miti Art. To proceed, please fill out the required details in the form linked below:\n\n"+
			"%s\n\n"+
			"This form will help us process your request efficiently. If you have any questions while completing it, feel free to reply to this email.\n\n"+
			"We look forward to reviewing your submission.\n\n"+
			"Best regards,\n"+
			"Poli Ninshuti Ndiramiye\n"+
			"Miti Art Super-user\n"+
			"+250 791 287 640",
		vendorOtherName, vendorFirstName, link)

	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", receiver)
	m.SetHeader("Subject", "Hello!! You are invited to join Miti Art")
	m.SetBody("text/plain", emailBody)

	d := gomail.NewDialer("smtp.gmail.com", 587, sender, emailPass)

	if err := d.DialAndSend(m); err != nil {
		log.Println("Failed to send email:", err)
		return false
	}

	log.Println("Email sent successfully!")
	// log.Println(token)
	return true
}
