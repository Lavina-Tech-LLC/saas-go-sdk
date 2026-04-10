package middleware

import (
	"context"
	"net/http"
	"strings"

	saassupport "github.com/Lavina-Tech-LLC/saas-go-sdk"
)

const (
	roleContextKey  contextKey = "saassupport_role"
	orgIDContextKey contextKey = "saassupport_orgId"
)

// WithRequireRole returns an http.Handler middleware that checks the
// authenticated user's organisation role against allowedRoles. It must be
// applied after WithAuth and after the consuming application stores the org ID
// via WithOrgID.
//
// On success the user's role is injected into the request context and can be
// read with GetRole(). If no allowedRoles are specified the middleware resolves
// and stores the role without performing an access check.
//
// Usage:
//
//	mux.Handle("/admin/", middleware.WithRequireRole(client, "owner", "admin")(handler))
func WithRequireRole(client *saassupport.Client, allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaims(r)
			if claims == nil {
				writeError(w, http.StatusUnauthorized, "Authentication required")
				return
			}

			orgID := GetOrgID(r)
			if orgID == "" {
				writeError(w, http.StatusForbidden, "Organization ID required")
				return
			}

			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				writeError(w, http.StatusUnauthorized, "Authentication required")
				return
			}
			token := strings.TrimPrefix(authHeader, "Bearer ")

			members, err := client.Auth.ListMembers(r.Context(), token, orgID)
			if err != nil {
				writeError(w, http.StatusForbidden, "Insufficient permissions")
				return
			}

			var role string
			for _, m := range members {
				if m.UserID == claims.UserID {
					role = m.Role
					break
				}
			}
			if role == "" {
				writeError(w, http.StatusForbidden, "Insufficient permissions")
				return
			}

			ctx := context.WithValue(r.Context(), roleContextKey, role)
			r = r.WithContext(ctx)

			if len(allowedRoles) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			for _, allowed := range allowedRoles {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}

			writeError(w, http.StatusForbidden, "Insufficient permissions")
		})
	}
}

// GetRole returns the user's organisation role from the request context.
// Returns an empty string if the role has not been resolved yet.
func GetRole(r *http.Request) string {
	role, _ := r.Context().Value(roleContextKey).(string)
	return role
}

// WithOrgID stores the organisation ID in the request context.
// Consuming applications should call this to make the org ID available
// to WithRequireRole (typically extracted from an X-Org-Id header).
func WithOrgID(r *http.Request, orgID string) *http.Request {
	ctx := context.WithValue(r.Context(), orgIDContextKey, orgID)
	return r.WithContext(ctx)
}

// GetOrgID returns the organisation ID from the request context.
// Returns an empty string if no org ID has been set.
func GetOrgID(r *http.Request) string {
	id, _ := r.Context().Value(orgIDContextKey).(string)
	return id
}
