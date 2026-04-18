package ginmiddleware

import (
	"strings"

	saassupport "github.com/Lavina-Tech-LLC/saas-go-sdk"
	"github.com/Lavina-Tech-LLC/saas-go-sdk/auth"
	"github.com/gin-gonic/gin"
)

const apiKeyHeader = "X-API-Key"

// AuthOrAPIKey returns a Gin middleware that accepts either an X-API-Key header
// (user-issued uk_live_... key) or an Authorization: Bearer JWT. It tries
// X-API-Key first; if absent, it falls back to the Bearer token. Aborts 401
// when neither is present or valid.
//
// On success (both paths) it sets "userId" and "email" in the Gin context,
// so RequireRole and existing handler code that reads those keys work
// transparently. On the JWT path "claims" is set (*auth.TokenClaims); on the
// API-key path "apiKeyClaims" is set (*auth.APIKeyClaims) and "orgId" is set.
//
// Usage:
//
//	client := saassupport.NewClient("pk_live_xxx")
//	r.Use(ginmiddleware.AuthOrAPIKey(client))
func AuthOrAPIKey(client *saassupport.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		if apiKey := c.GetHeader(apiKeyHeader); apiKey != "" {
			claims, err := client.VerifyAPIKey(c.Request.Context(), apiKey)
			if err != nil || claims == nil || !claims.Valid {
				c.AbortWithStatusJSON(401, gin.H{"error": "Invalid or expired API key"})
				return
			}
			c.Set("apiKeyClaims", claims)
			c.Set("userId", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("orgId", claims.OrgID)
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization required"})
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := client.VerifyToken(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Set("claims", claims)
		c.Set("userId", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}

// GetAPIKeyClaims retrieves verified API-key claims from the Gin context.
// Returns nil if the request was authenticated via JWT or not authenticated.
func GetAPIKeyClaims(c *gin.Context) *auth.APIKeyClaims {
	val, exists := c.Get("apiKeyClaims")
	if !exists {
		return nil
	}
	claims, _ := val.(*auth.APIKeyClaims)
	return claims
}
