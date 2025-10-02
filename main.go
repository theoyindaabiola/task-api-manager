package main

import (
	"log"
	"os"
	"taskapi/controllers"
	"taskapi/dao"
	"taskapi/routes"
	"taskapi/services"
	"taskapi/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env file: %w", err)
	}

	db := ConnectDB()

	// instance of gin from the routes handles routes request
	router := gin.Default()

	// main.go as the entry point, need the routes to use

	// run the goroutine, this runs forever to send out emails
	go utils.ProcessQueueMessages((os.Getenv("EMAIL_VERIFICATION_QUEUE"))) 
	go utils.ProcessQueueMessages((os.Getenv("EMAIL_RESET_QUEUE")))
	go utils.ProcessQueueMessages((os.Getenv("TASK_DELEGATION_QUEUE")))
	go utils.ProcessQueueMessages((os.Getenv("EMAIL_OTP_QUEUE")))
	go utils.ProcessQueueMessages((os.Getenv("SMS_OTP_QUEUE")))

	// Initialize dao, services and controller USER
	userDao := dao.NewUserDAO(db)
	userService := services.NewUserService(userDao)
	userController := controllers.NewUserController(userService)

	// Initialize dao, services and controller TASK
	taskDao := dao.NewTaskDAO(db)
	taskPermissionDAO := dao.NewTaskPermissionDAO(db)
	taskService := services.NewTaskService(taskDao, taskPermissionDAO, userDao)
	taskController := controllers.NewTaskController(taskService, userService)

	// Routessss
	routes.RegisterUserRoutes(router, userController, userService) 
	routes.RegisterRoutes(router, taskController, taskDao, taskPermissionDAO, userService) 

	log.Println("Server started at: 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
