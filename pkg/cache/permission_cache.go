package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	permissionCachePrefix = "user:permissions:"
	permissionCacheTTL    = 5 * time.Minute
)

type PermissionCache struct {
	redis *RedisClient
}

// NewPermissionCache creates a new permission cache
func NewPermissionCache(redis *RedisClient) *PermissionCache {
	return &PermissionCache{redis: redis}
}

// GetUserPermissions retrieves cached permissions for a user
func (c *PermissionCache) GetUserPermissions(ctx context.Context, userID uint) ([]string, error) {
	key := fmt.Sprintf("%s%d", permissionCachePrefix, userID)

	data, err := c.redis.Get(ctx, key)
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var permissions []string
	if err := json.Unmarshal([]byte(data), &permissions); err != nil {
		return nil, err
	}

	return permissions, nil
}

// SetUserPermissions caches permissions for a user
func (c *PermissionCache) SetUserPermissions(ctx context.Context, userID uint, permissions []string) error {
	key := fmt.Sprintf("%s%d", permissionCachePrefix, userID)

	data, err := json.Marshal(permissions)
	if err != nil {
		return err
	}

	return c.redis.Set(ctx, key, string(data), permissionCacheTTL)
}

// InvalidateUserPermissions removes cached permissions for a user
func (c *PermissionCache) InvalidateUserPermissions(ctx context.Context, userID uint) error {
	key := fmt.Sprintf("%s%d", permissionCachePrefix, userID)
	return c.redis.Delete(ctx, key)
}

// InvalidateAllPermissions removes all cached permissions
func (c *PermissionCache) InvalidateAllPermissions(ctx context.Context) error {
	pattern := permissionCachePrefix + "*"
	return c.redis.DeletePattern(ctx, pattern)
}
