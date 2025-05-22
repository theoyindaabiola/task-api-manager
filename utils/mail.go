package utils

import (
	"crypto/rand"
	// "crypto/tls"
	"encoding/hex"
	// "fmt"
	"log"
	"os"

	// "github.com/google/uuid"
	"gopkg.in/gomail.v2"
)

func SendMail(email string, body, messageType string) error {
	subject := "Email Verification"
	if messageType == "reset" {
        subject = "Password Reset"
    } else if messageType == "delegation" {
		subject = "Task Delegated"
	}

	// create an instance of gomail
	mail := gomail.NewMessage()
	mail.SetHeader("From", os.Getenv("EMAIL_SENDER"))
	mail.SetHeader("To", email)
	mail.SetHeader("Subject", subject)
	mail.SetBody("text/plain", body)
	// port := 587

	dialer := gomail.NewDialer(
		os.Getenv("SMTP_HOST"), 
		587,
		// 465, 
		os.Getenv("SMTP_USERNAME"), 
		os.Getenv("SMTP_PASSWORD"),
	)

	// for testing
	// if os.Getenv("SMTP_PORT") == "465" {
    //     port = 465
    // }
    // dialer := gomail.NewDialer(
    //     os.Getenv("SMTP_HOST"),
    //     port,
    //     os.Getenv("SMTP_USERNAME"),
    //     os.Getenv("SMTP_PASSWORD"),
    // )
    // if port == 465 {
    //     dialer.TLSConfig = &tls.Config{
    //         InsecureSkipVerify: true, // For testing only
    //         ServerName:        os.Getenv("SMTP_HOST"),
    //     }
    // }

	// fmt.Printf("SMTP: %s, User: %s, Sender: %s, To: %s\n", os.Getenv("SMTP_PORT"), os.Getenv("SMTP_HOST"), os.Getenv("SMTP_USERNAME"), os.Getenv("EMAIL_SENDER"), email)

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

// func GenerateVerificationCode() string {
//     token := uuid.New().String()
//     log.Printf("Generated verification token: %s", token)
//     return token
// }



// TESTING
// package utils

// import (
//     "crypto/rand"
//     "crypto/tls"
//     "encoding/hex"
//     "fmt"
//     "log"
//     "os"
//     "strconv"
//     "gopkg.in/gomail.v2"
// )

// func SendMail(email, body, messageType string) error {
//     subject := "Email Verification"
//     if messageType == "reset" {
//         subject = "Password Reset"
//     } else if messageType == "delegation" {
//         subject = "Task Delegated"
//     }

//     mail := gomail.NewMessage()
//     mail.SetHeader("From", os.Getenv("EMAIL_SENDER"))
//     mail.SetHeader("To", email)
//     mail.SetHeader("Subject", subject)
//     mail.SetBody("text/plain", body)

//     port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
//     if err != nil {
//         port = 587
//     }

//     dialer := gomail.NewDialer(
//         os.Getenv("SMTP_HOST"),
//         port,
//         os.Getenv("SMTP_USERNAME"),
//         os.Getenv("SMTP_PASSWORD"),
//     )

//     if port == 465 {
//         dialer.SSL = true
//         dialer.TLSConfig = &tls.Config{
//             InsecureSkipVerify: true, // For testing only
//             ServerName:        os.Getenv("SMTP_HOST"),
//             MinVersion:        tls.VersionTLS12, // Ensure TLS 1.2+
//         }
//     } else if port == 587 {
//         dialer.SSL = false
//         dialer.TLSConfig = &tls.Config{
//             ServerName: os.Getenv("SMTP_HOST"),
//             MinVersion: tls.VersionTLS12,
//         }
//     }

//     fmt.Printf("SMTP: %s:%d, User: %s, Sender: %s, To: %s\n", os.Getenv("SMTP_HOST"), port, os.Getenv("SMTP_USERNAME"), os.Getenv("EMAIL_SENDER"), email)
//     if err := dialer.DialAndSend(mail); err != nil {
//         log.Printf("Failed to send %s email to %s: %v", messageType, email, err)
//         return err
//     }
//     log.Printf("Sent %s email to %s", messageType, email)
//     return nil
// }

// func GenerateVerificationCode() string {
//     bytes := make([]byte, 16)
//     if _, err := rand.Read(bytes); err != nil {
//         log.Fatal("Failed to generate verification code: ", err)
//     }
//     return hex.EncodeToString(bytes)
// }