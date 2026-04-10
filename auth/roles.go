package auth

import "context"

// ListRoles returns all roles defined for the project.
func (s *Service) ListRoles(ctx context.Context, accessToken string) ([]Role, error) {
	var result []Role
	err := s.transport.RequestWithHeaders(ctx, "GET", "/auth/roles", nil, &result, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	return result, err
}
