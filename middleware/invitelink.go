package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	saassupport "github.com/Lavina-Tech-LLC/saas-go-sdk"
)

// InviteLinkInfoHandler returns an http.HandlerFunc that retrieves public info
// about an invite link. The codeExtractor function extracts the invite code from
// the request (router-agnostic).
//
// Usage with chi:
//
//	r.Get("/invite/{code}", middleware.InviteLinkInfoHandler(client, func(r *http.Request) string {
//	    return chi.URLParam(r, "code")
//	}))
func InviteLinkInfoHandler(client *saassupport.Client, codeExtractor func(*http.Request) string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := codeExtractor(r)
		if code == "" {
			writeError(w, http.StatusBadRequest, "Invite code is required")
			return
		}

		info, err := client.Auth.GetInviteLinkInfo(r.Context(), code)
		if err != nil {
			writeError(w, http.StatusNotFound, "Invite link not found or expired")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(info)
	}
}

// UseInviteLinkHandler returns an http.HandlerFunc that uses an invite link to
// join an organisation. Requires prior WithAuth middleware. The codeExtractor
// function extracts the invite code from the request.
//
// Usage with chi:
//
//	r.With(middleware.WithAuth(client)).
//	    Post("/invite/{code}/use", middleware.UseInviteLinkHandler(client, func(r *http.Request) string {
//	        return chi.URLParam(r, "code")
//	    }))
func UseInviteLinkHandler(client *saassupport.Client, codeExtractor func(*http.Request) string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := codeExtractor(r)
		if code == "" {
			writeError(w, http.StatusBadRequest, "Invite code is required")
			return
		}

		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			writeError(w, http.StatusUnauthorized, "Authorization required")
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		result, err := client.Auth.UseInviteLink(r.Context(), token, code)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}
