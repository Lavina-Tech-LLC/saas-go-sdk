package report

import "context"

// CreateEmbedToken issues a short-lived token for embedded dashboard access.
func (s *Service) CreateEmbedToken(ctx context.Context, params *CreateEmbedTokenParams) (*EmbedTokenResult, error) {
	result := &EmbedTokenResult{}
	err := s.transport.Request(ctx, "POST", "/reports/embed-tokens", params, result)
	return result, err
}

// ListEmbedTokens returns all active embed tokens for the project.
func (s *Service) ListEmbedTokens(ctx context.Context) ([]EmbedToken, error) {
	var result []EmbedToken
	err := s.transport.Request(ctx, "GET", "/reports/embed-tokens", nil, &result)
	return result, err
}

// RevokeEmbedToken invalidates a single embed token.
func (s *Service) RevokeEmbedToken(ctx context.Context, tokenID string) (*RevokeResult, error) {
	result := &RevokeResult{}
	err := s.transport.Request(ctx, "DELETE", "/reports/embed-tokens/"+tokenID, nil, result)
	return result, err
}
