package utils

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

// sendOTPEmail sends an OTP email
func SendOTPEmail(to string, otp uint64) error {
	from := "knowledgemartv01@gmail.com"
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	appPassword := os.Getenv("SMTPAPP")

	auth := smtp.PlainAuth("", from, appPassword, "smtp.gmail.com")

	msg := []byte("Subject: Verify your email\n\n" +
		fmt.Sprintf("Your OTP is %d", otp))
	err = smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, msg)
	if err != nil {
		fmt.Printf("Error in sending email: %v\n", err)
		return errors.New("failed to send email")
	}
	return nil
}
