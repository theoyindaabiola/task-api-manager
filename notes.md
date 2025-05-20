# USER
{
    "username": "Napoleon",
    "email": "johnnydaluv@yahoo.com",
    "password": "JohD@n" || 
}

// CREATE-TASK
// update the token in the Headers
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDcyMTc5MTAsImlzc3VlciI6InRhc2stYXBpLW1hbmFnZXIiLCJ1c2VyX2lkIjoiZTA1ZmM4NWItNjVkNS00NjgzLTg5MmItZWIwOGRkNGI0OGRlIn0.wzHkZv0lAMI_GiXh3Bxr7ON8hZs38toSTFFFVNz3eCU"
{
    "title": "John Tasks",
    "description": "John's first to-do"
}

// FORGOT-PASSWORD
{
    "email": "johnnydaluv@yahoo.com"
}

// RESET-PASSWORD
{
    "token": "99c4eff7-3220-489b-a0f1-2411870cdba5",
    "password": "JohDanToTheWorld"
}

# TASKOWNER
{
    "username": "Oyindamola",
    "email": "aodasola95@gmail.com",
    "password": "OyindaToTheWorld"
}

# DELEGATEE
{
    "username": "Dasola",
    "email": "aodasola@gmail.com",
    "password": "DasolaToTheWorld" || "dasolawelldone"
}

// TASK
{
    "title": "Oyindamola Tasks",
    "description": "Oyin's first middleware",
    "completed": "true"
}

# URLS
POST: localhost:8080/api/users/register
GET: localhost:8080/api/users/verify-email?code=bc8d24a254fa7975fa8525330b70e553
POST: localhost:8080/api/users/login
POST: localhost:8080/api/tasks/
POST: localhost:8080/api/users/forgot-password
POST: localhost:8080/api/users/reset-password


DB_USER=
DB_PASSWORD=
DB_NAME=neondb
DB_HOST=
DB_PORT=5432
SSL_MODE=require 
JWT_SECRET=


# PROJECT TASK
On the existing golang task api manager project complete the project by extending it based on the following features

1.⁠ ⁠Create a reset password for users who forgot their password
2.⁠ ⁠Add a feature that delegates a task to a user and give them a read and update permission
3.⁠ ⁠Add unit and integration test
Project is to be completed within a month. Also submission should be made via gitlab.







































TEST
package main

import (
	"bytes"
	"encoding/json"
	"strconv"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Connect to database using the sqlite
func setupTestDB() {
	var err error
	// uses the sqlite to use the local memorydisk machine to store data
	db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		// panic is for test
		panic("Failed to connect to the database")
	}
	db.AutoMigrate(&Task{})
}

// Set up the router
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default() // create a router for me
	
	// group route
	// taskRoute := r.Group("/")
	r.POST("/tasks", CreateTask)
	r.GET("/tasks", func(c *gin.Context)  {
		// store an array of multiple task in the variable task
		var tasks []Task
		db.Find(&tasks)
		c.JSON(http.StatusOK, tasks)
	})
	r.PUT("/tasks/:id", UpdateTask)
	r.DELETE("/tasks/:id", DeleteTask)
	return r
}

// functions unit testing
func TestTaskAPI(t *testing.T) {
	// call the db setup & router
	setupTestDB()
	router := setupRouter()

	// TEST CREATE TASK
	// create an instance of data to be used
	taskData := `{"title": "Oyindamola's first test", "completed": "true"}`
	// send post request in byte to be testable
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(([]byte(taskData))))
	// postman stuff!! content type
	req.Header.Set("Content-Type", "application/json")
	// test the endpoint
	w := httptest.NewRecorder()
	// route the request & test, and serve it a response
	router.ServeHTTP(w, req)
	// check that the response is statusCreated #201
	assert.Equal(t, http.StatusCreated, w.Code)

	// TEST GET TASKS
	req, _ = http.NewRequest("GET", "/tasks", nil)
	// same header
	req.Header.Set("Content-Type", "application/json")
	// same new endpoint test
	w = httptest.NewRecorder()
	// route the request and test "w"
	router.ServeHTTP(w, req)
	// check that the response is statusOK #200
	assert.Equal(t, http.StatusOK, w.Code)

	// confirm it's fetching task correctly
	var tasks []Task
	// serialize/convert to what GO understands and put it into the tasks var location 
	json.Unmarshal(w.Body.Bytes(), &tasks)
	// confirm the length of the response is correct
	assert.Len(t, tasks, 1)

	// TEST UPDATE TASK 
	// get the id of the tasked fetched above
	taskID := tasks[0].ID // at 1st index
	// create another instance of data to be used
	updatedData := `{"title": "Oyindamola's second test", "completed": "true"}`
	// send post request with the converted ID to string, instead of "/tasks/:id"
	req, _ = http.NewRequest("PUT", "/tasks/"+strconv.Itoa(int(taskID)), bytes.NewBuffer(([]byte(updatedData))))
	// postman stuff!! content type
	req.Header.Set("Content-Type", "application/json")
	// test the endpoint
	w = httptest.NewRecorder()
	// route the request & test, and serve it a response
	router.ServeHTTP(w, req)
	// check that the response is statusCreated #200
	assert.Equal(t, http.StatusOK, w.Code)

	// TEST DELETE TASK
	// here we can use the same id from the get and updated test
	// send post request with the converted ID to string, instead of "/tasks/:id"
	req, _ = http.NewRequest("DELETE", "/tasks/"+strconv.Itoa(int(taskID)), nil) // not passing any data
	// postman stuff!! content type !! NOT-NEED
	// req.Header.Set("Content-Type", "application/json")
	// test the endpoint
	w = httptest.NewRecorder()
	// route the request & test, and serve it a response
	router.ServeHTTP(w, req)
	// check that the response is StatusNoContent #204
	assert.Equal(t, http.StatusNoContent, w.Code)
}

