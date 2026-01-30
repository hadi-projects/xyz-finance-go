# XYZ Finance API

## Register
```shell
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "budi@mail.com",
    "password": "pAssword@123"
}'
```

Response:
```json
{
    "message":"User registered successfully.",
    "user":{
        "email":"budi@mail.com",
        "id":1
    }
}
```

## Login
```shell
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "budi@mail.com",
    "password": "pAsswj@1873"
}'
```

Response:
```json
{
    "message":"Login successful",
    "user":{
        "email":"user@example.com",
        "id":1
    }
}
```

## Get Limit
```shell
curl -X GET http://localhost:8080/api/user/limit \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJlbWFpbCI6ImJ1ZGlAbWFpbC5jb20iLCJleHAiOjE3Njk4OTk5MjUsImlhdCI6MTc2OTgxMzUyNX0.Jb4oK8NEwTXIzpZPZbozYtbOj0P1GRTiz9HWaVoqK6c"
```

Response:
```json
{
    "message":"Get Consumer Limit",
    "data":{
        "tenor_1": 100000,
        "tenor_2": 200000,
        "tenor_3": 500000,
        "tenor_6": 700000
    }
}
```