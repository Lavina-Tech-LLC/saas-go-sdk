package saassupport

import (
	"context"
	"net/http"
	"time"

	"github.com/Lavina-Tech-LLC/saas-go-sdk/auth"
	"github.com/Lavina-Tech-LLC/saas-go-sdk/billing"
	"github.com/Lavina-Tech-LLC/saas-go-sdk/report"
)

// ClientOption configures the Client.
type ClientOption func(*Client)

// WithBaseURL overrides the default API base URL.
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithHTTPClient provides a custom *http.Client.
func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// WithUserAgent sets a custom User-Agent string.
func WithUserAgent(ua string) ClientOption {
	return func(c *Client) {
		c.userAgent = ua
	}
}

// Client is the top-level SDK client for the SaaS Support platform.
type Client struct {
	// Auth provides access to the Auth module endpoints.
	Auth *auth.Service
	// Billing provides access to the Billing module endpoints.
	Billing *billing.Service
	// Report provides access to the Reports module endpoints.
	Report *report.Service

	apiKey     string
	baseURL    string
	httpClient *http.Client
	userAgent  string
}

// NewClient creates a new SDK client with the given API key.
func NewClient(apiKey string, opts ...ClientOption) *Client {
	c := &Client{
		apiKey:  apiKey,
		baseURL: DefaultURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		userAgent: "saas-support-go/" + Version,
	}

	for _, opt := range opts {
		opt(c)
	}

	transport := &Transport{
		APIKey:     c.apiKey,
		BaseURL:    c.baseURL,
		HTTPClient: c.httpClient,
		UserAgent:  c.userAgent,
	}

	c.Auth = auth.NewService(transport)
	c.Billing = billing.NewService(transport)
	c.Report = report.NewService(transport)

	return c
}

// VerifyToken is a convenience method that calls Auth.VerifyToken.
// This is the most common server-side operation.
func (c *Client) VerifyToken(ctx context.Context, token string) (*auth.TokenClaims, error) {
	return c.Auth.VerifyToken(ctx, &auth.VerifyTokenParams{Token: token})
}

// VerifyAPIKey is a convenience method that calls Auth.VerifyAPIKey.
// Used by middleware that accepts X-API-Key-authenticated requests.
func (c *Client) VerifyAPIKey(ctx context.Context, apiKey string) (*auth.APIKeyClaims, error) {
	return c.Auth.VerifyAPIKey(ctx, &auth.VerifyAPIKeyParams{APIKey: apiKey})
}
