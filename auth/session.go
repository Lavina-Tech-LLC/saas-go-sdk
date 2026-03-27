package auth

import "context"

// Refresh exchanges a refresh token for a new access/refresh token pair.
func (s *Service) Refresh(ctx context.Context, params *RefreshParams) (*RefreshResult, error) {
	result := &RefreshResult{}
	err := s.transport.Request(ctx, "POST", "/auth/refresh", params, result)
	return result, err
}

// Logout revokes the session associated with the refresh token.
func (s *Service) Logout(ctx context.Context, params *LogoutParams) (*LogoutResult, error) {
	result := &LogoutResult{}
	err := s.transport.Request(ctx, "POST", "/auth/logout", params, result)
	return result, err
}
