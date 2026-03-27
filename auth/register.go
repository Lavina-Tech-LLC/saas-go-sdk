package auth

import "context"

// Register creates a new end-user account and returns auth tokens.
func (s *Service) Register(ctx context.Context, params *RegisterParams) (*AuthResult, error) {
	result := &AuthResult{}
	err := s.transport.Request(ctx, "POST", "/auth/register", params, result)
	return result, err
}
