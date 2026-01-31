# XYZ Finance API

Backend API untuk aplikasi pembiayaan XYZ Finance yang dibangun menggunakan Go dengan Gin Framework.

## Tech Stack

- **Language**: Go 1.25+
- **Framework**: [Gin](https://github.com/gin-gonic/gin)
- **Database**: MySQL with [GORM](https://gorm.io/)
- **Authentication**: JWT (JSON Web Tokens)
- **Logging**: [Zerolog](https://github.com/rs/zerolog) with file rotation
- **Testing**: Go testing + [Testify](https://github.com/stretchr/testify) + [GoMock](https://github.com/uber-go/mock)

## Project Structure

```
xyz-finance/
├── cmd/api/              # Application entry point
├── config/               # Configuration management
├── internal/
│   ├── dto/              # Data Transfer Objects
│   ├── entity/           # Database models
│   ├── handler/          # HTTP request handlers
│   ├── middleware/       # HTTP middlewares (JWT, CORS, Rate Limit, etc.)
│   ├── repository/       # Data access layer
│   ├── router/           # Route definitions
│   └── service/          # Business logic layer
├── pkg/
│   ├── database/         # Database connection & seeding
│   ├── logger/           # Structured logging
│   └── validator/        # Custom validators
├── storage/
│   ├── logs/             # Log files (audit, auth, system)
│   └── uploads/          # File uploads (KTP, selfie images)
└── migrations/           # Database migrations
```

## Prerequisites

- Go 1.25 atau lebih baru
- MySQL 8.0+
- Make (optional, untuk menjalankan Makefile)

## Installation

1. **Clone repository**
   ```bash
   git clone https://github.com/hadi-projects/xyz-finance-go.git
   cd xyz-finance-go
   ```

2. **Salin file environment dan sesuaikan konfigurasi**
   ```bash
   cp .env-example .env
   ```

3. **Edit file `.env`**
   ```env
   # App
   APP_PORT=8080
   APP_ENV=development

   # Database
   DB_USER=root
   DB_PASSWORD=your_password
   DB_HOST=localhost
   DB_PORT=3306
   DB_NAME=xyz_finance

   # CORS Configuration
   CORS_ALLOWED_ORIGINS=http://localhost:3000
   CORS_ALLOW_CREDENTIALS=true

   # Rate Limiter Configuration
   RATE_LIMIT_RPS=10
   RATE_LIMIT_BURST=20

   # Request Timeout (seconds)
   REQUEST_TIMEOUT=30

   # JWT Configuration
   JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
   JWT_EXPIRY_HOURS=24

   # API Key
   API_KEY=your-api-key
   ```

4. **Install dependencies**
   ```bash
   go mod tidy
   ```

5. **Buat database MySQL**
   ```sql
   CREATE DATABASE xyz_finance;
   ```

## Running the Application

### Development Mode

```bash
# Menggunakan Makefile
make run

# Atau langsung dengan Go
go run cmd/api/main.go
```

### Build Production Binary

```bash
make build
./bin/api
```

### Running with Docker

```bash
# Build and start all services
make docker-up

# View logs
make docker-logs

# Stop all services
make docker-down

# Rebuild and restart
make docker-build
make docker-restart

# Clean up (remove volumes and images)
make docker-clean
```

**Docker Services:**
- **API**: `http://localhost:8080`
- **MySQL**: `localhost:3307` (mapped to container port 3306)

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/handler/...
go test ./internal/service/...
```

## Default Seeded Users

Aplikasi akan otomatis seed data berikut saat pertama kali dijalankan:

| Email              | Password       | Role  |
|--------------------|----------------|-------|
| admin@mail.com     | pAsswj@123     | Admin |
| budi@mail.com      | pAsswj@1873    | User  |
| annisa@mail.com    | pAsswj@1763    | User  |

## Authentication

API menggunakan kombinasi **API Key** dan **JWT Token**:

- **API Key**: Dikirim via header `X-API-KEY`
- **JWT Token**: Dikirim via header `Authorization: Bearer <token>`

## API Endpoints

### Public Routes
| Method | Endpoint             | Description         |
|--------|----------------------|---------------------|
| GET    | `/health`            | Health check        |
| POST   | `/api/auth/register` | Register user       |
| POST   | `/api/auth/login`    | Login user          |
| GET    | `/uploads/*filepath` | Static files        |

### Protected Routes (Requires API Key + JWT)
| Method | Endpoint              | Permission           | Description            |
|--------|-----------------------|----------------------|------------------------|
| GET    | `/api/user/profile`   | -                    | Get user profile       |
| GET    | `/api/limit/`         | `get-limit`          | Get user limits        |
| POST   | `/api/limit/`         | `create-limit`       | Create limit (Admin)   |
| PUT    | `/api/limit/:id`      | `edit-limit`         | Update limit (Admin)   |
| DELETE | `/api/limit/:id`      | `delete-limit`       | Delete limit (Admin)   |
| POST   | `/api/transaction/`   | `create-transaction` | Create transaction     |
| GET    | `/api/transaction/`   | `get-transactions`   | Get transactions       |
| GET    | `/api/logs/audit`     | `get-audit-log`      | Get audit logs (Admin) |
| GET    | `/api/logs/auth`      | `get-auth-log`       | Get auth logs (Admin)  |

## API Examples

### Register
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@mail.com",
    "password": "pAssword@123"
}'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "budi@mail.com",
    "password": "pAsswj@1873"
}'
```

### Get Profile (Protected)
```bash
curl -X GET http://localhost:8080/api/user/profile \
  -H "Content-Type: application/json" \
  -H "X-API-KEY: your-api-key" \
  -H "Authorization: Bearer <your-jwt-token>"
```

### Get Limits (Protected)
```bash
curl -X GET http://localhost:8080/api/limit/ \
  -H "Content-Type: application/json" \
  -H "X-API-KEY: your-api-key" \
  -H "Authorization: Bearer <your-jwt-token>"
```

### Create Transaction (Protected)
```bash
curl -X POST http://localhost:8080/api/transaction/ \
  -H "Content-Type: application/json" \
  -H "X-API-KEY: your-api-key" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "contract_number": "CTR-2024-001",
    "otr": 600000,
    "admin_fee": 10000,
    "installment_amount": 105000,
    "interest_amount": 10000,
    "asset_name": "Samsung Galaxy A05",
    "tenor": 6
}'
```

## Performance Testing

Tersedia script untuk performance testing menggunakan k6:

```bash
# Install k6 (macOS)
brew install k6

# Run performance test
k6 run performance-test.js
```

## Postman Collection

Import file `XYZ Finance API.postman_collection.json` ke Postman untuk testing API secara interaktif.

## Logging

Aplikasi menggunakan structured logging dengan 3 jenis log file:

- `storage/logs/audit.log` - Log aktivitas pengguna (CRUD operations)
- `storage/logs/auth.log` - Log autentikasi (login, register, etc.)
- `storage/logs/system.log` - Log sistem (startup, errors, etc.)

Log otomatis di-rotate ketika mencapai 10MB dengan retensi 7 file.

## Security Features

- **JWT Authentication** dengan refresh token
- **API Key** validation
- **Rate Limiting** untuk mencegah abuse
- **CORS** configuration
- **XSS Protection** middleware
- **Security Headers** middleware
- **Request Timeout** untuk mencegah slow loris attacks
- **RBAC** (Role-Based Access Control)

## License

MIT License