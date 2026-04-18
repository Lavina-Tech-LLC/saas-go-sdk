package auth

import "context"

// VerifyAPIKey validates a user-issued API key (uk_live_…) server-side.
// The backend returns {Valid: false} for unknown, revoked, or expired keys
// rather than a transport error — callers should branch on Valid, not err.
func (s *Service) VerifyAPIKey(ctx context.Context, params *VerifyAPIKeyParams) (*APIKeyClaims, error) {
	result := &APIKeyClaims{}
	err := s.transport.Request(ctx, "POST", "/auth/verify-api-key", params, result)
	return result, err
}
