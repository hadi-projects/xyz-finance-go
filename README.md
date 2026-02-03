# XYZ Finance API

Backend API untuk aplikasi pembiayaan XYZ Finance yang dibangun menggunakan Go dengan Gin Framework.

**[Lihat Dokumentasi Arsitektur](docs/architecture.md)** - Layer structure, security flow, database schema

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.25+ |
| Framework | [Gin](https://github.com/gin-gonic/gin) |
| Database | MySQL 8.0 + [GORM](https://gorm.io/) |
| Auth | JWT + API Key |
| Logging | [Zerolog](https://github.com/rs/zerolog) |
| Testing | [Testify](https://github.com/stretchr/testify) + [GoMock](https://github.com/uber-go/mock) |
| Container | Docker + Docker Compose |

## Quick Start

### Prerequisites
- Go 1.25+
- MySQL 8.0+
- Docker (optional)

### Installation

```bash
# Clone repository
git clone https://github.com/hadi-projects/xyz-finance-go.git
cd xyz-finance-go

# Setup environment
cp .env-example .env
# Edit .env dengan konfigurasi database

# Install dependencies
go mod tidy

# Run
make run
```

### Running with Docker

```bash
make docker-up      # Start all services
make docker-logs    # View logs
make docker-down    # Stop services
```

**Services:** API (`localhost:8080`) | MySQL (`localhost:3307`)

### Running Tests

```bash
make test           # Run all tests
make test-cover     # Run with coverage
```

## Default Users

| Email              | Password       | Role  |
|--------------------|----------------|-------|
| admin@mail.com     | pAsswj@123     | Admin |
| budi@mail.com      | pAsswj@1873    | User  |
| annisa@mail.com    | pAsswj@1763    | User  |

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

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "budi@mail.com", "password": "pAsswj@1873"}'
```

### Protected Request
```bash
curl -X GET http://localhost:8080/api/user/profile \
  -H "X-API-KEY: your-api-key" \
  -H "Authorization: Bearer <jwt-token>"
```

## Additional Resources

- **Postman Collection**: `XYZ Finance API.postman_collection.json`
- **Performance Testing**: `k6 run performance-test.js`

## License

MIT License
