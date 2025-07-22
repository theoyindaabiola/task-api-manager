# PROJECT TASK
On the existing golang task api manager project complete the project by extending it based on the following features

1.⁠ ⁠Create a reset password for users who forgot their password
2.⁠ ⁠Add a feature that delegates a task to a user and give them a read and update permission
3.⁠ ⁠Add unit and integration test
Project is to be completed within a month. Also submission should be made via gitlab.


# USER
{
    "username": "Napoleon",
    "email": "johnnydaluv@yahoo.com",
    "password": "JohD@n" || 
}
verify
login
Napoleon Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTMzMjMzMzgsImlzc3VlciI6InRhc2stYXBpLW1hbmFnZXIiLCJ1c2VyX2lkIjoiMDE3ZDZjOGEtNDBjYy00NTNiLWEwNDItN2ZmYjE3ZGUxNjY5In0.z47aDbDU0MzMTat0bPAKd9AUhOL7ppklWdc7Ze1Js18"

{
  "title": "First Task",
  "description": "This is John's first delegated task.",
  "completed": false
}

TaskID: d3b7fbff-b87e-40bf-b0e8-5586f298de70

# FORGOT-PASSWORD
{
    "email": "johnnydaluv@yahoo.com"
}

# RESET-PASSWORD
{
    "token": "30d4f2784b104bd22ed3250e601e4bb6", // reset password token
    "password": "JohD@n123"
}

# TASKOWNER
{
    "username": "Oyindamola",
    "email": "aodasola95@gmail.com",
    "password": "OyindaToTheWorld"
}
verify
login 
token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTMzMTAyNzksImlzc3VlciI6InRhc2stYXBpLW1hbmFnZXIiLCJ1c2VyX2lkIjoiNDRmMjNhNTUtMjIxOC00YjVlLTlkODAtMjYzYzhkZTc5N2MwIn0.WmeHaCzIhnSfxOzhMoWZ9kbzzyidgpuG6QEyr9nWYe4"

# CREATE-TASK
update the token in the Headers
    "delegator": ""
# TASK
TaskID: d3b7fbff-b87e-40bf-b0e8-5586f298de70



# DELEGATEE 
{
    "username": "Dasola",
    "email": "aodasola@gmail.com",
    "password": "DasolaToTheWorld" || "dasolawelldone"
}

"delegatee": ""


# DELEGATE 
{
  "delegatee_id": "44f23a55-2218-4b5e-9d80-263c8de797c0",
  "permission": "R"  // or "U"
}

### URLS
# USER
POST: localhost:8080/api/users/register
GET: localhost:8080/api/users/verify-email?code=dea01a277d5e485f4aa3699e20b7a1c6
POST: localhost:8080/api/users/login

# PASSWORD
POST: localhost:8080/api/users/forgot-password
POST: localhost:8080/api/users/reset-password

# TASK & DELEGATION
POST: localhost:8080/api/tasks/     // create task
GET: localhost:8080/api/tasks/<task-id>     // to access task
DELETE: localhost:8080/api/tasks/<task-id>/delete      // delete task by owner
PUT: localhost:8080/api/tasks/<task-id>       // update task

# DELEGATE TASK (using the delegator's token in the header)
POST: localhost:8080/api/tasks/d3b7fbff-b87e-40bf-b0e8-5586f298de70/delegate       // delegate task by owner
PATCH: localhost:8080/api/tasks/<task-id>/permission    // update the permission by owner
DELETE localhost:8080/api/tasks/<task-id>/permission    // revoke permission by owner


DB_USER=
DB_PASSWORD=
DB_NAME=neondb
DB_HOST=
DB_PORT=5432
SSL_MODE=require 
JWT_SECRET=
