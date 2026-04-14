package ginmiddleware

import (
	"strings"

	saassupport "github.com/Lavina-Tech-LLC/saas-go-sdk"
	"github.com/Lavina-Tech-LLC/saas-go-sdk/auth"
	"github.com/gin-gonic/gin"
)

// RequireRole returns a Gin middleware that checks the authenticated user's
// organisation role against allowedRoles. It must be used after Auth() and
// after the consuming application sets "orgId" in the Gin context (typically
// from an X-Org-Id header).
//
// On success the user's role is stored in the context and can be read with
// GetRole(). If no allowedRoles are specified the middleware resolves and
// stores the role without performing an access check.
//
// Usage:
//
//	r.Use(ginmiddleware.Auth(client))
//	r.Use(func(c *gin.Context) {
//	    c.Set("orgId", c.GetHeader("X-Org-Id"))
//	    c.Next()
//	})
//	admin := r.Group("/admin", ginmiddleware.RequireRole(client, "owner", "admin"))
func RequireRole(client *saassupport.Client, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("userId")
		if userID == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authentication required"})
			return
		}

		orgID := c.GetString("orgId")
		if orgID == "" {
			c.AbortWithStatusJSON(403, gin.H{"error": "Organization ID required"})
			return
		}

		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authentication required"})
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		members, err := client.Auth.ListMembers(c.Request.Context(), token, orgID)
		if err != nil {
			c.AbortWithStatusJSON(403, gin.H{"error": "Insufficient permissions"})
			return
		}

		var member *auth.Member
		for i := range members {
			if members[i].UserID == userID {
				member = &members[i]
				break
			}
		}
		if member == nil || member.Role == "" {
			c.AbortWithStatusJSON(403, gin.H{"error": "Insufficient permissions"})
			return
		}

		c.Set("role", member.Role)

		if len(allowedRoles) == 0 {
			c.Next()
			return
		}
		if member.HasAnyRole(allowedRoles...) {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(403, gin.H{"error": "Insufficient permissions"})
	}
}

// GetRole returns the user's organisation role from the Gin context.
// Returns an empty string if the role has not been resolved yet.
func GetRole(c *gin.Context) string {
	return c.GetString("role")
}
