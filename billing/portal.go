package billing

import "context"

// CreatePortalToken issues a short-lived portal access token for a customer.
func (s *Service) CreatePortalToken(ctx context.Context, params *CreatePortalTokenParams) (*PortalTokenResult, error) {
	result := &PortalTokenResult{}
	err := s.transport.Request(ctx, "POST", "/billing/portal/token", params, result)
	return result, err
}
