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
Napoleon Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTQzMTMxNDQsImlzc3VlciI6InRhc2stYXBpLW1hbmFnZXIiLCJ1c2VyX2lkIjoiYWZkMTliNGEtZWViOS00YzVmLWI2MDMtMGMzYzg5ZGE5NzYyIn0.p8lyvdOVC4LYSr9oNmFk9uRCisb41NXZoTThGkbL3ts"

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
    "token": "bbe79748bf19edb790db69264bdd2ba1", // reset password token
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
token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbmFibGVkXzJmYSI6dHJ1ZSwiZXhwIjoxNzU4NTA4MjAwLCJpc190b3RwX3ZlcmlmaWVkIjpmYWxzZSwiaXNzdWVyIjoidGFzay1hcGktbWFuYWdlciIsInVzZXJfaWQiOiI3MjhjYmY1ZS0yN2NkLTQwNGItYTQwYS1hYzBlNjE4NWY4MjQifQ.JZWKxaVW2bTWTyXZXoZYzZPQ1Mj8KVpxi6dy-q59vuY"

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
  "delegatee_id": "afd19b4a-eeb9-4c5f-b603-0c3c89da9762",
  "permission": "R"  // or "U"
}

### URLS
# USER
POST: localhost:8080/api/users/register
GET: localhost:8080/api/users/verify-email?code=064bc7d09f5daf558c553fe2e53b8985
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

# 2FA - TOTP 
POST: localhost:8080/api/users/enable-totp // enable2fa-totp
POST: localhost:8080/api/users/verify-totp // verify-totp-code