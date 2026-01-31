# XYZ Finance - Improvement TODO

## Database Optimization

- [ ] Add database indexes for frequently queried columns
  - [ ] `CREATE INDEX idx_transactions_user_id ON transactions(user_id)`
  - [ ] `CREATE INDEX idx_limit_mutations_user_id ON limit_mutations(user_id)`
  - [ ] `CREATE INDEX idx_users_email ON users(email)`
  - [ ] `CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id)`

- [ ] Optimize GORM connection pooling
  - [ ] Configure `SetMaxIdleConns`
  - [ ] Configure `SetMaxOpenConns`
  - [ ] Configure `SetConnMaxLifetime`

## Caching (Redis)

- [ ] Setup Redis container in docker-compose
- [ ] Implement Redis client
- [ ] Cache user permissions (TTL: 5 min)
- [ ] Cache JWT token validation
- [ ] Cache tenor limits data

## Security Improvements

- [ ] Add request ID middleware for tracing
- [ ] Implement refresh token rotation
- [ ] Add password reset functionality
- [ ] Add email verification
- [ ] Implement account lockout after failed attempts

## API Improvements

- [ ] Add pagination for list endpoints
  - [ ] GET /api/transaction/
  - [ ] GET /api/limit/
  - [ ] GET /api/logs/*
- [ ] Add filtering & sorting options
- [ ] Add API versioning (v1, v2)
- [ ] Implement soft delete for entities

## Testing

- [ ] Increase test coverage to > 80%
- [ ] Add integration tests
- [ ] Add E2E tests with testcontainers
- [ ] Add benchmark tests

## Monitoring & Observability

- [ ] Add Prometheus metrics endpoint
- [ ] Setup Grafana dashboards
- [ ] Implement distributed tracing (Jaeger/Zipkin)
- [ ] Add health check for database connection

## Documentation

- [ ] Add Swagger/OpenAPI documentation
- [ ] Add API changelog
- [ ] Add contribution guidelines

## DevOps

- [ ] Setup CI/CD pipeline (GitHub Actions)
- [ ] Add Kubernetes manifests
- [ ] Setup staging environment
- [ ] Add automated database migrations

---

## Priority Matrix

| Priority | Task | Impact | Effort |
|----------|------|--------|--------|
| 🔴 High | Database Indexes | High | Low |
| 🔴 High | Connection Pooling | Medium | Low |
| 🟡 Medium | Redis Caching | High | Medium |
| 🟡 Medium | Pagination | Medium | Low |
| 🟡 Medium | Swagger Docs | Medium | Medium |
| 🟢 Low | Prometheus Metrics | Medium | Medium |
| 🟢 Low | Kubernetes | Low | High |
