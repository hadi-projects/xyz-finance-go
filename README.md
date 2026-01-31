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
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJlbWFpbCI6ImJ1ZGlAbWFpbC5jb20iLCJleHAiOjE3Njk5MjA4MzUsImlhdCI6MTc2OTgzNDQzNX0.Rt6_rhXrOcYcxSAKNO31NbWkrs8Sd28ROEaK0oKCr30"
```

user id, email, consumer { nik, fullname, legal name, place of birth, date of birth, salary, ktp image, self image}

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
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJlbWFpbCI6ImJ1ZGlAbWFpbC5jb20iLCJleHAiOjE3Njk5MjA4MzUsImlhdCI6MTc2OTgzNDQzNX0.Rt6_rhXrOcYcxSAKNO31NbWkrs8Sd28ROEaK0oKCr30"
```

Response:
```json
{
    "data":[
        {
            "id":1,
            "tenor_month":1,
            "limit_amount":100000,
            "created_at":"2026-01-31T06:20:46.255+07:00",
            "updated_at":"2026-01-31T06:20:46.255+07:00"
        },{
            "id":2,
            "tenor_month":2,
            "limit_amount":200000,
            "created_at":"2026-01-31T06:20:46.264+07:00",
            "updated_at":"2026-01-31T06:20:46.264+07:00"
        },{
            "id":3,
            "tenor_month":3,
            "limit_amount":500000,
            "created_at":"2026-01-31T06:20:46.269+07:00",
            "updated_at":"2026-01-31T06:20:46.269+07:00"
        },{
            "id":4,
            "tenor_month":6,
            "limit_amount":700000,
            "created_at":"2026-01-31T06:20:46.271+07:00",
            "updated_at":"2026-01-31T06:20:46.271+07:00"
        },{
            "id":17,
            "tenor_month":1,
            "limit_amount":100000,
            "created_at":"2026-01-31T06:27:02.508+07:00",
            "updated_at":"2026-01-31T06:27:02.508+07:00"
        },{
            "id":18,
            "tenor_month":2,
            "limit_amount":200000,
            "created_at":"2026-01-31T06:27:02.512+07:00",
            "updated_at":"2026-01-31T06:27:02.512+07:00"
        },{
            "id":19,
            "tenor_month":3,
            "limit_amount":500000,
            "created_at":"2026-01-31T06:27:02.515+07:00",
            "updated_at":"2026-01-31T06:27:02.515+07:00"
        },{
            "id":20,
            "tenor_month":6,
            "limit_amount":700000,
            "created_at":"2026-01-31T06:27:02.518+07:00",
            "updated_at":"2026-01-31T06:27:02.518+07:00"
        },{
            "id":25,"tenor_month":1,"limit_amount":100000,"created_at":"2026-01-31T06:46:20.886+07:00","updated_at":"2026-01-31T06:46:20.886+07:00"},{"id":26,"tenor_month":2,"limit_amount":200000,"created_at":"2026-01-31T06:46:20.891+07:00","updated_at":"2026-01-31T06:46:20.891+07:00"},{"id":27,"tenor_month":3,"limit_amount":500000,"created_at":"2026-01-31T06:46:20.893+07:00","updated_at":"2026-01-31T06:46:20.893+07:00"},{"id":28,"tenor_month":6,"limit_amount":700000,"created_at":"2026-01-31T06:46:20.896+07:00","updated_at":"2026-01-31T06:46:20.896+07:00"},{"id":33,"tenor_month":1,"limit_amount":100000,"created_at":"2026-01-31T11:50:14.105+07:00","updated_at":"2026-01-31T11:50:14.105+07:00"},{"id":34,"tenor_month":2,"limit_amount":200000,"created_at":"2026-01-31T11:50:14.11+07:00","updated_at":"2026-01-31T11:50:14.11+07:00"},{"id":35,"tenor_month":3,"limit_amount":500000,"created_at":"2026-01-31T11:50:14.114+07:00","updated_at":"2026-01-31T11:50:14.114+07:00"},{"id":36,"tenor_month":6,"limit_amount":700000,"created_at":"2026-01-31T11:50:14.117+07:00","updated_at":"2026-01-31T11:50:14.117+07:00"}],"message":"Get User Limits"}
```

