package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/hadi-projects/xyz-finance-go/pkg/logger"
)

const (
	// Cache key prefixes
	UserCachePrefix        = "user:"
	LimitCachePrefix       = "limits:"
	TransactionCachePrefix = "transactions:"

	// Cache TTLs
	UserCacheTTL  = 5 * time.Minute
	LimitCacheTTL = 10 * time.Minute
)

// CacheInvalidator provides methods for invalidating cached data
type CacheInvalidator struct {
	redis     *RedisClient
	permCache *PermissionCache
}

// NewCacheInvalidator creates a new cache invalidator
func NewCacheInvalidator(redis *RedisClient, permCache *PermissionCache) *CacheInvalidator {
	return &CacheInvalidator{
		redis:     redis,
		permCache: permCache,
	}
}

// InvalidateUserCache invalidates all cached data for a user
func (c *CacheInvalidator) InvalidateUserCache(ctx context.Context, userID uint) error {
	var errs []error

	// Invalidate user permissions
	if c.permCache != nil {
		if err := c.permCache.InvalidateUserPermissions(ctx, userID); err != nil {
			logger.SystemLogger.Warn().Err(err).Uint("user_id", userID).Msg("Failed to invalidate user permissions cache")
			errs = append(errs, err)
		}
	}

	// Invalidate user profile cache
	key := fmt.Sprintf("%sprofile:%d", UserCachePrefix, userID)
	if err := c.redis.Delete(ctx, key); err != nil {
		logger.SystemLogger.Warn().Err(err).Uint("user_id", userID).Msg("Failed to invalidate user profile cache")
		errs = append(errs, err)
	}

	// Invalidate user limits cache
	limitKey := fmt.Sprintf("%suser:%d", LimitCachePrefix, userID)
	if err := c.redis.Delete(ctx, limitKey); err != nil {
		logger.SystemLogger.Warn().Err(err).Uint("user_id", userID).Msg("Failed to invalidate user limits cache")
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("cache invalidation had %d errors", len(errs))
	}

	logger.SystemLogger.Debug().Uint("user_id", userID).Msg("User cache invalidated successfully")
	return nil
}

// InvalidateUserLimits invalidates only the limits cache for a user
func (c *CacheInvalidator) InvalidateUserLimits(ctx context.Context, userID uint) error {
	key := fmt.Sprintf("%suser:%d", LimitCachePrefix, userID)
	if err := c.redis.Delete(ctx, key); err != nil {
		logger.SystemLogger.Warn().Err(err).Uint("user_id", userID).Msg("Failed to invalidate user limits cache")
		return err
	}
	logger.SystemLogger.Debug().Uint("user_id", userID).Msg("User limits cache invalidated")
	return nil
}

// InvalidateAllLimits invalidates all limits cache
func (c *CacheInvalidator) InvalidateAllLimits(ctx context.Context) error {
	pattern := LimitCachePrefix + "*"
	if err := c.redis.DeletePattern(ctx, pattern); err != nil {
		logger.SystemLogger.Warn().Err(err).Msg("Failed to invalidate all limits cache")
		return err
	}
	logger.SystemLogger.Debug().Msg("All limits cache invalidated")
	return nil
}

// InvalidateUserTransactions invalidates transaction cache for a user
func (c *CacheInvalidator) InvalidateUserTransactions(ctx context.Context, userID uint) error {
	pattern := fmt.Sprintf("%suser:%d:*", TransactionCachePrefix, userID)
	if err := c.redis.DeletePattern(ctx, pattern); err != nil {
		logger.SystemLogger.Warn().Err(err).Uint("user_id", userID).Msg("Failed to invalidate user transactions cache")
		return err
	}
	logger.SystemLogger.Debug().Uint("user_id", userID).Msg("User transactions cache invalidated")
	return nil
}

// InvalidateOnLimitChange should be called when a user's limit is created, updated, or deleted
func (c *CacheInvalidator) InvalidateOnLimitChange(ctx context.Context, userID uint) error {
	return c.InvalidateUserLimits(ctx, userID)
}

// InvalidateOnTransactionCreate should be called when a new transaction is created
func (c *CacheInvalidator) InvalidateOnTransactionCreate(ctx context.Context, userID uint) error {
	// Invalidate limits (consumed limit changes) and transactions
	if err := c.InvalidateUserLimits(ctx, userID); err != nil {
		return err
	}
	return c.InvalidateUserTransactions(ctx, userID)
}

// InvalidateOnRoleChange should be called when a user's role/permissions change
func (c *CacheInvalidator) InvalidateOnRoleChange(ctx context.Context, userID uint) error {
	if c.permCache != nil {
		return c.permCache.InvalidateUserPermissions(ctx, userID)
	}
	return nil
}
