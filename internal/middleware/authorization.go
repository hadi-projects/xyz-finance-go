package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/xyz-finance-go/internal/repository"
)

func PermissionMiddleware(userRepo repository.UserRepository, requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		user, err := userRepo.FindByID(userID.(uint))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// DEBUG LOG
		// fmt.Printf("DEBUG: UserID: %d, Role: %s, Permissions: %v\n", user.ID, user.Role.Name, user.Role.Permissions)

		hasPermission := false
		for _, perm := range user.Role.Permissions {
			if perm.Name == requiredPermission {
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
