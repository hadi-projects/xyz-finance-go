# XYZ Finance - Performance Improvement TODO

## Current Performance Baseline

| Endpoint | Current p95 | Target p95 |
|----------|-------------|------------|
| Login | 117ms | < 80ms |
| Get Profile | 8ms | < 5ms |
| Get Limits | 21ms | < 15ms |
| Get Transactions | 37ms | < 20ms |
| Create Transaction | 30ms | < 25ms |

---

## 1. Database Optimization

### Indexes
- [ ] Add index on `transactions.user_id`
- [ ] Add index on `limit_mutations.user_id`  
- [ ] Add index on `users.email`
- [ ] Add index on `refresh_tokens.user_id`
- [ ] Add composite index on `transactions(user_id, created_at)`

### Connection Pooling
- [ ] Configure `SetMaxIdleConns(10)`
- [ ] Configure `SetMaxOpenConns(100)`
- [ ] Configure `SetConnMaxLifetime(1h)`

### Query Optimization
- [ ] Use `Select()` untuk limit kolom yang di-query
- [ ] Replace N+1 queries dengan `Preload`
- [ ] Use raw SQL untuk complex queries

---

## 2. Authentication Optimization

### BCrypt Cost
- [ ] Reduce bcrypt cost dari 10 ke 8 (development)
- [ ] Keep cost 10-12 untuk production

### JWT Caching
- [ ] Cache parsed JWT tokens (in-memory)
- [ ] Cache user permissions lookup

---

## 3. Caching Layer (Redis)

### Setup
- [ ] Add Redis to docker-compose
- [ ] Implement Redis client wrapper

### Cache Strategies
- [ ] Cache user profile (TTL: 5 min)
- [ ] Cache tenor limits (TTL: 10 min)
- [ ] Cache permissions (TTL: 5 min)
- [ ] Implement cache invalidation

---

## 4. API Response Optimization

### Response Compression
- [ ] Enable Gzip compression middleware
- [ ] Compress responses > 1KB

### Pagination
- [ ] Add pagination for GET /api/transaction/
- [ ] Add pagination for GET /api/limit/
- [ ] Default limit: 20, Max: 100

---

## 5. Code Optimization

### Goroutines
- [ ] Parallel database queries where applicable
- [ ] Use worker pools untuk batch operations

### Memory
- [ ] Use sync.Pool untuk frequently allocated objects
- [ ] Reduce allocations in hot paths

---

## Priority & Effort Matrix

| Task | Impact | Effort | Priority |
|------|--------|--------|----------|
| Database Indexes | ⬆️ High | 🟢 Low | 🔴 P1 |
| Connection Pooling | ⬆️ Medium | 🟢 Low | � P1 |
| BCrypt Cost (dev) | ⬆️ High (login) | 🟢 Low | � P1 |
| Gzip Compression | ⬆️ Medium | 🟢 Low | 🟡 P2 |
| Query Optimization | ⬆️ High | 🟡 Medium | 🟡 P2 |
| Redis Caching | ⬆️ High | 🔴 High | 🟢 P3 |
| Pagination | ⬆️ Medium | 🟡 Medium | 🟢 P3 |
