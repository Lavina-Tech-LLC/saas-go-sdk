package middleware

import (
	"context"
	"net/http"
	"strings"

	saassupport "github.com/Lavina-Tech-LLC/saas-go-sdk"
	"github.com/Lavina-Tech-LLC/saas-go-sdk/auth"
)

const (
	apiKeyClaimsContextKey contextKey = "saassupport_apikey_claims"
	apiKeyHeader                      = "X-API-Key"
)

// WithAuthOrAPIKey returns an http.Handler middleware that accepts either an
// X-API-Key header (user-issued uk_live_... key) or an Authorization: Bearer
// JWT. It tries X-API-Key first; if absent, it falls back to the Bearer token.
// Responds 401 when neither header is present or valid.
//
// On API-key success, both GetAPIKeyClaims and GetUserID/GetEmail resolve.
// On JWT success, GetClaims and GetUserID/GetEmail resolve; GetAPIKeyClaims is nil.
func WithAuthOrAPIKey(client *saassupport.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if apiKey := r.Header.Get(apiKeyHeader); apiKey != "" {
				claims, err := client.VerifyAPIKey(r.Context(), apiKey)
				if err != nil || claims == nil || !claims.Valid {
					writeError(w, http.StatusUnauthorized, "Invalid or expired API key")
					return
				}
				ctx := context.WithValue(r.Context(), apiKeyClaimsContextKey, claims)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				writeError(w, http.StatusUnauthorized, "Authorization required")
				return
			}
			token := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := client.VerifyToken(r.Context(), token)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), claimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetAPIKeyClaims retrieves verified API-key claims from the request context.
// Returns nil if the request was authenticated via JWT or not authenticated at all.
func GetAPIKeyClaims(r *http.Request) *auth.APIKeyClaims {
	claims, _ := r.Context().Value(apiKeyClaimsContextKey).(*auth.APIKeyClaims)
	return claims
}
