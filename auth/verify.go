package auth

import "context"

// VerifyToken validates a project-scoped JWT server-side.
// This is the primary method for backend services to verify end-user tokens.
func (s *Service) VerifyToken(ctx context.Context, params *VerifyTokenParams) (*TokenClaims, error) {
	result := &TokenClaims{}
	err := s.transport.Request(ctx, "POST", "/auth/verify-token", params, result)
	return result, err
}
