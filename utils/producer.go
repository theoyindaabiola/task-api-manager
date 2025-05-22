package utils

import (
	"encoding/json"
	"log"
	"os"
	amqp "github.com/rabbitmq/amqp091-go"
)

// understand that we are replacing the server with the rabbitmq to help send the 
// verification messaage to the users at the point of registration.
type EmailTask struct {
    Email         	string  `json:"email"`
    Code          	string  `json:"code"`
    Type          	string	`json:"type"`
    TaskID        	string	`json:"task_id"`
    DelegateeID   	string 	`json:"delegatee_id"`
    DelegateeEmail 	string	`json:"delegatee_email"`
    TaskTitle     	string  `json:"task_title"`
}

func PublishMessage (queue, email, code, messageType, delegateeID, delegateeEmail, taskTitle string) error {
	task := EmailTask{
		Email:          email,
        Code:           code,
        Type:           messageType,
        TaskID:         code,
        DelegateeID:    delegateeID,
        DelegateeEmail: delegateeEmail,
        TaskTitle:      taskTitle,
	}

	message, err := json.Marshal(task)
    if err != nil {
		log.Printf("Failed to marshal message: %v", err)
        return err
    }

	// connect to the rabbitmq server
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_SERVER"))
	if err != nil {
		log.Println("Failed to connect to RabbitMQ:", err)
		return err
	}
	// if anything pannicks, close the connection
	defer conn.Close()

	// now connect to a channel/routes
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to open channel: %v", err)
		return err
	}
	defer ch.Close()

	// now create a queue
	q, err := ch.QueueDeclare(
		// os.Getenv("EMAIL_VERIFICATION_QUEUE"), // routing key
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to declare a queue: %v", err)
		return err
	}

	// if all works well, then publish the message
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body: message, // this contain the verification code
		},
	)
	if err != nil {
		log.Printf("Failed to publish a message: %v", err)
		return err
	}
	log.Printf("Published %s task to %s for %s", messageType, queue, email)
	return nil
}