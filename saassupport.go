// Package saassupport provides a Go SDK for the SaaS Support platform.
//
// Create a client using your API key:
//
//	client := saassupport.NewClient("pk_live_xxx")
//	result, err := client.Auth.Login(ctx, &auth.LoginParams{
//	    Email:    "user@example.com",
//	    Password: "secret",
//	})
package saassupport

const (
	// Version is the current SDK version.
	Version = "0.1.0"

	// DefaultURL is the default SaaS Support API base URL.
	DefaultURL = "https://api.saas-support.com/v1"
)
