package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// GORM
type User struct {
	// json is in the go application while the db i s mapped to it in the database
	ID					uuid.UUID 		`gorm:"primarykey" json:"id"`
	Username 			string 			`gorm:"unique;not null" json:"username"`
	Password			string 			`gorm:"not null" json:"-"`
	Email				string 			`gorm:"unique;not null" json:"email"`
	Verified			bool   			`gorm:"default:false" json:"verified"`
	VerificationToken 	string 			`gorm:"unique" json:"-"`
	ResetToken			*string			`gorm:"unique" json:"-"`
	ResetTokenExpiresAt *time.Time      `gorm:"default:null" json:"-"`
	CreatedAt			time.Time 		`json:"created_at"`
	UpdatedAt   		time.Time 		`json:"updated_at"`
	DeletedAt   		gorm.DeletedAt 	`gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return nil
}






// func (dao *UserDAO) CreateUserDB(user *models.User) error {
// 	// gorm needs the instance of User{} not the user struct
// 	if err := dao.DB.Create(user).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (dao *UserDAO) GetUserVerification(VerificationToken string) (*models.User, error) {
// 	var user models.User
// 	if err := dao.DB.Where("verification_token = ?", VerificationToken).First(&user).Error; err != nil {
// 		return nil, err
// 	}
// 	return &user, nil
// }

// package utils

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"os"

// 	amqp "github.com/rabbitmq/amqp091-go"
// )

// func EmailConsumer(queueName string) error {
// 	// connect to the rabbitmq server: consumer goes straight for pick up
// 	conn, err := amqp.Dial(os.Getenv("RABBITMQ_SERVER"))
// 	if err != nil {
// 		log.Println("Failed to connect to RabbiteMQ", err)
// 		return err
// 	}
// 	defer conn.Close()

// 	// connect to a channel
// 	ch, err := conn.Channel()
// 	if err != nil {
// 		log.Println("Failed to connect to a channel", err)
// 		return err
// 	}
// 	defer ch.Close()

// 	// redeclare the quueue here ?? testing 
// 	// now create a queue
// 	q, err := ch.QueueDeclare (
// 		queueName,
// 		true,
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		log.Printf("Failed to declare a queue: %v", err)
// 		return err
// 	}

// 	// now comsume/pick up the message from the queue
// 	messages, err := ch.Consume(
// 		q.Name, // pickup point
// 		"", // consumer tag
// 		true, // auto acknowledge
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		log.Printf("Failed to declare a queue: %v", err)
// 		return err
// 	}

// 	forever := make(chan bool)
	
// 	// use goroutine to loop the messages we have
// 	go func () {
// 		for msg := range messages {
// 			var task EmailTask
// 			if err := json.Unmarshal(msg.Body, &task); err != nil {
//                 log.Printf("Failed to unmarshal task from %s: %v", queueName, err)
//                 continue
//             }
// 			var url, body string
// 			if task.Type == "verification" {
//                 url = fmt.Sprintf("http://localhost:8080/api/users/verify-email?code=%s", task.Code)
// 				body = fmt.Sprintf("Click to verify your account: %s", url)
//             } else if task.Type == "reset" {
//                 url = fmt.Sprintf("http://localhost:8080/api/users/reset-password?token=%s", task.Code)
// 				body = fmt.Sprintf("Click to reset your password: %s", url)
//             } else if task.Type == "delegation" {
// 				url = fmt.Sprintf("http://localhost:8080/api/tasks/%s", task.TaskID)
// 				body = fmt.Sprintf("Task '%s' has been delegated to you. View details: %s", task.TaskTitle, url)
// 			} else {
// 				log.Printf("Unknown task type %s from %s", task.Type, queueName)
//                 continue
// 			}

// 			if err := SendMail(task.Email, body, task.Type); err != nil {
//                 log.Printf("Failed to send %s email to %s: %v", task.Type, task.Email, err)
//                 continue
//             }
// 			log.Printf("Sent %s email to %s", task.Type, task.Email)
// 		}
// 	}()

// 	fmt.Println("[*] Waiting for messages. To exit press CTRL+C")
// 	<-forever
// 	return nil
// }

// package utils

// import (
// 	"encoding/json"
// 	"log"
// 	"os"
// 	amqp "github.com/rabbitmq/amqp091-go"
// )

// // understand that we are replacing the server with the rabbitmq to help send the 
// // verification messaage to the users at the point of registration.
// type EmailTask struct {
//     Email         	string  `json:"email"`
//     Code          	string  `json:"code"`
//     Type          	string	`json:"type"`
//     TaskID        	string	`json:"task_id"`
//     DelegateeID   	string 	`json:"delegatee_id"`
//     DelegateeEmail 	string	`json:"delegatee_email"`
//     TaskTitle     	string  `json:"task_title"`
// }

// func PublishMessage (queue, email, code, messageType, delegateeID, delegateeEmail, taskTitle string) error {
// 	task := EmailTask{
// 		Email:          email,
//         Code:           code,
//         Type:           messageType,
//         TaskID:         code, // For delegation, code is TaskID
//         DelegateeID:    delegateeID,
//         DelegateeEmail: delegateeEmail,
//         TaskTitle:      taskTitle,
// 	}

// 	message, err := json.Marshal(task)
//     if err != nil {
// 		log.Printf("Failed to marshal message: %v", err)
//         return err
//     }

// 	// connect to the rabbitmq server
// 	conn, err := amqp.Dial(os.Getenv("RABBITMQ_SERVER"))
// 	if err != nil {
// 		log.Println("Failed to connect to RabbitMQ:", err)
// 		return err
// 	}
// 	// if anything pannicks, close the connection
// 	defer conn.Close()

// 	// now connect to a channel/routes
// 	ch, err := conn.Channel()
// 	if err != nil {
// 		log.Printf("Failed to open channel: %v", err)
// 		return err
// 	}
// 	defer ch.Close()

// 	// now create a queue
// 	q, err := ch.QueueDeclare(
// 		// os.Getenv("EMAIL_VERIFICATION_QUEUE"), // routing key
// 		queue,
// 		true,
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		log.Printf("Failed to declare a queue: %v", err)
// 		return err
// 	}

// 	// if all works well, then publish the message
// 	err = ch.Publish(
// 		"",
// 		q.Name,
// 		false,
// 		false,
// 		amqp.Publishing{
// 			ContentType: "application/json",
// 			Body: message, // this contain the verification code
// 		},
// 	)
// 	if err != nil {
// 		log.Printf("Failed to publish a message: %v", err)
// 		return err
// 	}
// 	log.Printf("Published %s task to %s for %s", messageType, queue, email)
// 	return nil
// }

// // GORM
// type User struct {
// 	// json is in the go application while the db i s mapped to it in the database
// 	ID					uuid.UUID 		`gorm:"primarykey" json:"id"`
// 	Username 			string 			`gorm:"unique;not null" json:"username"`
// 	Password			string 			`gorm:"not null" json:"-"`
// 	Email				string 			`gorm:"unique;not null" json:"email"`
// 	Verified			bool   			`gorm:"default:false" json:"verified"`
// 	VerificationToken 	string 			`gorm:"unique" json:"-"`
// 	ResetToken			*string			`gorm:"unique" json:"-"`
// 	ResetTokenExpiresAt *time.Time      `gorm:"default:null" json:"-"`
// 	CreatedAt			time.Time 		`json:"created_at"`
// 	UpdatedAt   		time.Time 		`json:"updated_at"`
// 	DeletedAt   		gorm.DeletedAt 	`gorm:"index" json:"-"`
// }
