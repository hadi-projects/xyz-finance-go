# XYZ Finance - Performance Improvement TODO

## Current Performance (After BCrypt Optimization)

| Endpoint | Before | After | Target | Status |
|----------|--------|-------|--------|--------|
| Login | 117ms | **54ms** | < 80ms | ✅ Achieved |
| Get Profile | 8ms | 13ms | < 15ms | ✅ Achieved |
| Get Limits | 21ms | 16ms | < 20ms | ✅ Achieved |
| Get Transactions | 37ms | 31ms | < 35ms | ✅ Achieved |
| Create Transaction | 30ms | 27ms | < 30ms | ✅ Achieved |

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
- [x] Reduce bcrypt cost dari 10 ke 8 (development)
- [x] Keep cost 10-12 untuk production

### JWT Caching
- [x] Cache parsed JWT tokens (in-memory)
- [x] Cache user permissions lookup

---

## 3. Caching Layer (Redis)

### Setup
- [x] Add Redis to docker-compose
- [x] Implement Redis client wrapper

### Cache Strategies
- [x] Cache user profile (TTL: 5 min)
- [x] Cache tenor limits (TTL: 10 min)
- [x] Cache permissions (TTL: 5 min)
- [ ] Implement cache invalidation

---

## 4. API Response Optimization

### Response Compression
- [ ] Enable Gzip compression middleware
- [ ] Compress responses > 1KB

### Pagination
- [x] Add pagination for GET /api/transaction/
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
