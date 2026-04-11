package ginmiddleware

import (
	"strings"

	saassupport "github.com/Lavina-Tech-LLC/saas-go-sdk"
	"github.com/gin-gonic/gin"
)

// InviteLinkInfo returns a gin.HandlerFunc that retrieves public info about
// an invite link. Expects a `:code` URL parameter.
//
// Usage:
//
//	r.GET("/invite/:code", ginmiddleware.InviteLinkInfo(client))
func InviteLinkInfo(client *saassupport.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Param("code")
		if code == "" {
			c.JSON(400, gin.H{"error": "Invite code is required"})
			return
		}

		info, err := client.Auth.GetInviteLinkInfo(c.Request.Context(), code)
		if err != nil {
			c.JSON(404, gin.H{"error": "Invite link not found or expired"})
			return
		}

		c.JSON(200, info)
	}
}

// UseInviteLink returns a gin.HandlerFunc that uses an invite link to join
// an organisation. Requires prior Auth() middleware.
//
// Usage:
//
//	protected := r.Group("/", ginmiddleware.Auth(client))
//	protected.POST("/invite/:code/use", ginmiddleware.UseInviteLink(client))
func UseInviteLink(client *saassupport.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Param("code")
		if code == "" {
			c.JSON(400, gin.H{"error": "Invite code is required"})
			return
		}

		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{"error": "Authorization required"})
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		result, err := client.Auth.UseInviteLink(c.Request.Context(), token, code)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, result)
	}
}
