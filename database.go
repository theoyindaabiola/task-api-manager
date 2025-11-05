package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"taskapi/models"
)

// gorm instance
// var db *gorm.DB

// loads the env and connect to gorm then migrate
func ConnectDB() *gorm.DB {
	var err error

	// create and read the env variables strings
	connectStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := gorm.Open(postgres.Open(connectStr), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: %w", err)
	}

	err = db.AutoMigrate(
		&models.Task{}, 
		&models.User{}, 
		&models.TaskDelegation{}, 
		&models.OtpVerification{},
	)
	if err != nil {
		log.Fatal("Failed to run Migration", err)
	}

	return db
}
