package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendMail(email string, body, messageType string) error {
	subject := "Email Verification"
	switch messageType {
		case "reset":
        	subject = "Password Reset"
    	case "delegation":
			subject = "Task Delegated"
		case "otp":
			subject = "2F Authentication"
	}

	// create an instance of gomail
	mail := gomail.NewMessage()
	mail.SetHeader("From", os.Getenv("EMAIL_SENDER"))
	mail.SetHeader("To", email)
	mail.SetHeader("Subject", subject)
	mail.SetBody("text/plain", body)
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Println("Invalid SMTP_PORT:", err)
		port = 2525 // fallback
	}

	// dialer is Elastic Email
	dialer := gomail.NewDialer(
		os.Getenv("SMTP_HOST"), 
		port,
		os.Getenv("SMTP_USERNAME"), 
		os.Getenv("SMTP_PASSWORD"),
	)

	if err := dialer.DialAndSend(mail); err != nil {
		log.Printf("Failed to send %s email to %s: %v ", messageType, email, err)
		return err
	}
	log.Printf("Sent %s email to %s", messageType, email)
	return nil
}

func GenerateVerificationCode() string {
	// create byte
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal("Failed to generate verification code: ", err)
	}
	return hex.EncodeToString(bytes)
}

// generate 2fa OTP
func Generate2FACode() string {
	// random 6 digit code 0 - 999999
	num, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		log.Fatal("Failed to generate OTP code: ", err)
	}
	return fmt.Sprintf("%06d", num.Int64())
}
