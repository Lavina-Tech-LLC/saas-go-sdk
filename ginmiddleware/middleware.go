package ginmiddleware

import (
	"strings"

	saassupport "github.com/Lavina-Tech-LLC/saas-support-go"
	"github.com/Lavina-Tech-LLC/saas-support-go/auth"
	"github.com/gin-gonic/gin"
)

// Auth returns a Gin middleware that verifies the Bearer token and sets
// "claims", "userId", and "email" in the Gin context.
//
// Usage:
//
//	client := saassupport.NewClient("pk_live_xxx")
//	r.Use(ginmiddleware.Auth(client))
//
//	// In handler:
//	claims := ginmiddleware.GetClaims(c)
func Auth(client *saassupport.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
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

// OptionalAuth is like Auth but does not reject unauthenticated requests.
// If a valid token is present, claims are set in the Gin context.
func OptionalAuth(client *saassupport.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := client.VerifyToken(c.Request.Context(), token)
			if err == nil {
				c.Set("claims", claims)
				c.Set("userId", claims.UserID)
				c.Set("email", claims.Email)
			}
		}
		c.Next()
	}
}

// GetClaims retrieves verified token claims from the Gin context.
// Returns nil if the request was not authenticated.
func GetClaims(c *gin.Context) *auth.TokenClaims {
	val, exists := c.Get("claims")
	if !exists {
		return nil
	}
	claims, _ := val.(*auth.TokenClaims)
	return claims
}

// MustGetClaims is like GetClaims but panics if claims are not present.
// Use only in routes protected by Auth middleware.
func MustGetClaims(c *gin.Context) *auth.TokenClaims {
	return c.MustGet("claims").(*auth.TokenClaims)
}
