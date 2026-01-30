# XYZ Finance API

## Register
```shell
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "pAssword@123"
}'
```

Response:
```json
{
    "message":"User registered successfully.",
    "user":{
        "email":"user@example.com",
        "id":0
    }
}
```