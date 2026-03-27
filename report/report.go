package report

import "context"

// Requester sends HTTP requests to the SaaS Support API.
type Requester interface {
	Request(ctx context.Context, method, path string, body interface{}, dest interface{}) error
	RequestWithHeaders(ctx context.Context, method, path string, body interface{}, dest interface{}, headers map[string]string) error
}

// Service provides access to Report module endpoints.
type Service struct {
	transport Requester
}

// NewService creates a new Report service. Called by saassupport.NewClient.
func NewService(transport Requester) *Service {
	return &Service{transport: transport}
}
