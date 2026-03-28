package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	saassupport "github.com/Lavina-Tech-LLC/saas-go-sdk"
	"github.com/Lavina-Tech-LLC/saas-go-sdk/auth"
)

type contextKey string

const (
	claimsContextKey contextKey = "saassupport_claims"
)

// WithAuth returns an http.Handler middleware that verifies the Bearer token
// in the Authorization header and injects token claims into the request context.
//
// Usage:
//
//	client := saassupport.NewClient("pk_live_xxx")
//	mux.Handle("/api/", middleware.WithAuth(client)(myHandler))
func WithAuth(client *saassupport.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

// WithOptionalAuth is like WithAuth but does not reject unauthenticated requests.
// If a valid token is present, claims are injected into context.
// If no token is present, the request proceeds without claims.
func WithOptionalAuth(client *saassupport.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				token := strings.TrimPrefix(authHeader, "Bearer ")
				claims, err := client.VerifyToken(r.Context(), token)
				if err == nil {
					ctx := context.WithValue(r.Context(), claimsContextKey, claims)
					r = r.WithContext(ctx)
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// GetClaims retrieves the verified token claims from the request context.
// Returns nil if the request was not authenticated.
func GetClaims(r *http.Request) *auth.TokenClaims {
	claims, _ := r.Context().Value(claimsContextKey).(*auth.TokenClaims)
	return claims
}

// GetUserID returns the user ID from the context.
// Returns an empty string if the request was not authenticated.
func GetUserID(r *http.Request) string {
	if claims := GetClaims(r); claims != nil {
		return claims.UserID
	}
	return ""
}

// GetEmail returns the email from the context.
// Returns an empty string if the request was not authenticated.
func GetEmail(r *http.Request) string {
	if claims := GetClaims(r); claims != nil {
		return claims.Email
	}
	return ""
}

func writeError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
