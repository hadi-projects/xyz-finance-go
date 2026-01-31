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

## Get Profile
```shell
curl -X GET http://localhost:8080/api/user/profile \
  -H "Content-Type: application/json" \
  -H "X-API-KEY: biytf7rciyubyt6r7g89py" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJlbWFpbCI6ImJ1ZGlAbWFpbC5jb20iLCJleHAiOjE3Njk5MjQ5NDAsImlhdCI6MTc2OTgzODU0MH0._SA7QNtqzBF6PrLMbun8MoRPSHkHyWkjbmemwTK6iKA"
```

Response:
```json
{
    "data":{
        "user_id":2,
        "email":"budi@mail.com",
        "consumer":{
            "nik":"1234567890123456","full_name":"Budi Santoso","legal_name":"Budi Santoso","place_of_birth":"Jakarta","date_of_birth":"1990-01-01T00:00:00+07:00",
            "salary":10000000,"ktp_image":"ktp_placeholder.jpg","selfie_image":"selfie_placeholder.jpg"
        }
    },
    "message":"Get User Profile"
}
```

## Get Limit
```shell
curl -X GET http://localhost:8080/api/limit/ \
  -H "Content-Type: application/json" \
  -H "X-API-KEY: biytf7rciyubyt6r7g89py" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJlbWFpbCI6ImJ1ZGlAbWFpbC5jb20iLCJleHAiOjE3Njk5MjQ5NDAsImlhdCI6MTc2OTgzODU0MH0._SA7QNtqzBF6PrLMbun8MoRPSHkHyWkjbmemwTK6iKA"
```

Response:
```json
{
    "data": [
        {
            "user_id": 2,
            "tenor_month": 1,
            "limit_amount": 100000
        },
        {
            "user_id": 2,
            "tenor_month": 2,
            "limit_amount": 200000
        },
        {
            "user_id": 2,
            "tenor_month": 3,
            "limit_amount": 500000
        },
        {
            "user_id": 2,
            "tenor_month": 6,
            "limit_amount": 700000
        }
    ],
    "message": "Get User Limits"
}
```

