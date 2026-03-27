package auth

import "context"

// GetMe returns the authenticated end-user's profile.
func (s *Service) GetMe(ctx context.Context, accessToken string) (*User, error) {
	result := &User{}
	err := s.transport.RequestWithHeaders(ctx, "GET", "/auth/me", nil, result, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	return result, err
}

// UpdateMe updates the authenticated end-user's metadata.
func (s *Service) UpdateMe(ctx context.Context, accessToken string, params *UpdateProfileParams) (*User, error) {
	result := &User{}
	err := s.transport.RequestWithHeaders(ctx, "PATCH", "/auth/me", params, result, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	return result, err
}
