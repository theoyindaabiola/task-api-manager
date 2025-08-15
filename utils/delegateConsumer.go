package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

func DelegationConsumer() error{
    conn, err := amqp091.Dial(os.Getenv("RABBITMQ_URL"))
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("Failed to open a channel: %v", err)
    }
    defer ch.Close()

    queue := os.Getenv("TASK_DELEGATION_QUEUE")
    q, err := ch.QueueDeclare(
        queue,
        true,  // durable
        false, // autoDelete
        false, // exclusive
        false, // noWait
        nil,   // args
    )
    if err != nil {
        log.Fatalf("Failed to declare queue: %v", err)
    }

    msgs, err := ch.Consume(
        q.Name,
        "",    // consumer
        true,  // autoAck
        false, // exclusive
        false, // noLocal
        false, // noWait
        nil,   // args
    )
    if err != nil {
        log.Fatalf("Failed to register consumer: %v", err)
    }

	forever := make(chan bool)

    go func() {
        for msg := range msgs {
            var delegation struct {
                TaskID        	string	`json:"task_id"`
                DelegateeID   	string	`json:"delegatee_id"`
                DelegateeEmail 	string	`json:"delegatee_email"`
                TaskTitle     	string	`json:"task_title"`
            }
            if err := json.Unmarshal(msg.Body, &delegation); err != nil {
                log.Printf("Error unmarshaling delegation message: %v", err)
                continue
            }
            log.Printf("Task %s delegated to %s (%s): %s",
                delegation.TaskID, delegation.DelegateeEmail, delegation.DelegateeID, delegation.TaskTitle)
        }
    }()

    fmt.Printf("Started consumer for %s", queue)
	<-forever
	return nil
}