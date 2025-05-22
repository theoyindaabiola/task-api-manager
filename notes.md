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


{
    "title": "John Tasks",
    "description": "John's first to-do"
}  

# FORGOT-PASSWORD
{
    "email": "johnnydaluv@yahoo.com"
}

# RESET-PASSWORD
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
verify
login

# CREATE-TASK
update the token in the Headers
    "delegator": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDgwMDYzOTgsImlzc3VlciI6InRhc2stYXBpLW1hbmFnZXIiLCJ1c2VyX2lkIjoiNDU4ZmIwMTMtMzUzYi00OTFmLTk5YjktNzZiNWVhYjMzYzg2In0.V1OQVKiC2Rrbi_PBtW2iHBzqcumXAPYqlkbds4z4eZ8"
# TASK
{
    "title": "Oyindamola Tasks",
    "description": "Oyin's first middleware",
    "completed": "false"
}
TaskID: 5203a290-1c9f-4ecc-b1c4-3fb3d2a8a359



# DELEGATEE 
{
    "username": "Dasola",
    "email": "aodasola@gmail.com",
    "password": "DasolaToTheWorld" || "dasolawelldone"
}

"delegatee": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDgwMDY4MTUsImlzc3VlciI6InRhc2stYXBpLW1hbmFnZXIiLCJ1c2VyX2lkIjoiNmViNTIyYjktZjEwYi00MDE1LWIzMTItNjA0N2M1YmJmODBmIn0.iMVaho1eai6r2Kje2bz8D9G0P8gVoUb_GYXsu9a_dm4"



# DELEGATE TASK (using the delegator's token in the header)
URL: POST http://localhost:8080/api/tasks/5203a290-1c9f-4ecc-b1c4-3fb3d2a8a359/delegate

{
  "delegatee_id": "6eb522b9-f10b-4015-b312-6047c5bbf80f",
  "permission": "R"  // or "U"
}


# URLS
POST: localhost:8080/api/users/register
GET: localhost:8080/api/users/verify-email?code=87d010e4101f9354a82561e0f797794f
POST: localhost:8080/api/users/login
POST: localhost:8080/api/tasks/
POST: localhost:8080/api/users/forgot-password
POST: localhost:8080/api/users/reset-password
GET: http://localhost:8080/api/tasks/:id // to access task


DB_USER=
DB_PASSWORD=
DB_NAME=neondb
DB_HOST=
DB_PORT=5432
SSL_MODE=require 
JWT_SECRET=
