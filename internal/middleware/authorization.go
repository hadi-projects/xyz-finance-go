package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/xyz-finance-go/internal/repository"
	"github.com/hadi-projects/xyz-finance-go/pkg/cache"
	"github.com/hadi-projects/xyz-finance-go/pkg/logger"
)

func PermissionMiddleware(userRepo repository.UserRepository, permCache *cache.PermissionCache, requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		uid := userID.(uint)
		ctx := context.Background()

		// Try to get permissions from cache first
		if permCache != nil {
			permissions, err := permCache.GetUserPermissions(ctx, uid)
			if err != nil {
				logger.SystemLogger.Warn().Err(err).Msg("Failed to get permissions from cache")
			} else if permissions != nil {
				// Cache hit - check permission
				for _, perm := range permissions {
					if perm == requiredPermission {
						c.Next()
						return
					}
				}
				c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient permissions"})
				c.Abort()
				return
			}
		}

		// Cache miss - get from database
		user, err := userRepo.FindByID(uid)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Extract permission names
		var permissionNames []string
		for _, perm := range user.Role.Permissions {
			permissionNames = append(permissionNames, perm.Name)
		}

		// Cache permissions for future requests
		if permCache != nil {
			if err := permCache.SetUserPermissions(ctx, uid, permissionNames); err != nil {
				logger.SystemLogger.Warn().Err(err).Msg("Failed to cache permissions")
			}
		}

		// Check permission
		hasPermission := false
		for _, perm := range permissionNames {
			if perm == requiredPermission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
