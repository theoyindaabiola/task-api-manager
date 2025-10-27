package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ProcessQueueMessages(queueName string) error {
	log.Printf("Consumer started for queue: %s", queueName)
	// connect to the rabbitmq server: consumer goes straight for pick up
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_SERVER"))
	if err != nil {
		log.Println("Failed to connect to RabbiteMQ", err)
		return err
	}
	defer conn.Close()

	// connect to a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Println("Failed to connect to a channel", err)
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		os.Getenv("SMS_OTP_QUEUE"),
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// now comsume/pick up the message from the queue
	messages, err := ch.Consume(
		queueName, // pickup point
		"", // consumer tag
		true, // auto acknowledge
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to declare a queue: %v", err)
		return err
	}

	forever := make(chan bool)
	
	// use goroutine to loop the messages we have
	go func () {
		for msg := range messages {
			var task MessageTask
			if err := json.Unmarshal(msg.Body, &task); err != nil {
                log.Printf("Failed to unmarshal task from %s: %v", queueName, err)
                continue
            }
			var url, body string
			switch task.Type {
				case "verification":
					url = fmt.Sprintf("http://localhost:8080/api/users/verify-email?code=%s", task.Code)
					body = fmt.Sprintf("Click to verify your account: %s", url)
				case "reset":
					url = fmt.Sprintf("http://localhost:8080/api/users/reset-password?token=%s", task.Code)
					body = fmt.Sprintf("Click to reset your password: %s", url)
				case "delegation":
					url = fmt.Sprintf("http://localhost:8080/api/tasks/%s", task.TaskID)
					body = fmt.Sprintf("Task '%s' has been delegated to you. View details: %s", task.TaskTitle, url)
				case "sms_otp":
					err := SendSMSOTP(task.Recipient, task.Code)
					if err != nil {
						log.Printf("Failed to send SMS OTP to %s: %v", task.Recipient, err)
						continue
					}
					log.Printf("Sent SMS OTP to %s", task.Recipient)
					continue
				default:
					log.Printf("Unknown task type %s from %s", task.Type, queueName)
				continue
			}
	
			if err := SendMail(task.Recipient, body, task.Type); err != nil {
				log.Printf("Failed to send %s email to %s: %v", task.Type, task.Recipient, err)
				continue
			}
			log.Printf("Sent %s email to %s", task.Type, task.Recipient)
		}
	}()

	fmt.Println("[*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
}

// SendSMSOTP sends OTP via Termii
func SendSMSOTP(phoneNumber, otp string) error {
	apiKey := os.Getenv("TERMII_API_KEY")
	senderID := os.Getenv("TERMII_SENDER_ID")
	baseURL := os.Getenv("TERMII_BASE_URL")

	if apiKey == "" || senderID == "" || baseURL == "" {
		return fmt.Errorf("missing Termii configuration")
	}

	message := fmt.Sprintf("Your TaskAPI verification code is: %s", otp)

	payload := map[string]interface{}{
		"to":      phoneNumber,
		"from":    senderID,
		"sms":     message,
		"type":    "plain",
		"channel": "generic",
		"api_key": apiKey,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal Termii payload: %v", err)
	}

	resp, err := http.Post(baseURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send Termii request: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("Sent SMS OTP to %s via Termii", phoneNumber)
		return nil
	}

	return fmt.Errorf("termii API error: %s", string(respBody))
}
