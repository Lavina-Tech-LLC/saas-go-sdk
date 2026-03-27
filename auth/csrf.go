package auth

import "context"

// GetCSRFToken returns a CSRF token for the authenticated user's session.
func (s *Service) GetCSRFToken(ctx context.Context, accessToken string) (*CSRFTokenResult, error) {
	result := &CSRFTokenResult{}
	err := s.transport.RequestWithHeaders(ctx, "GET", "/auth/csrf", nil, result, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	return result, err
}
