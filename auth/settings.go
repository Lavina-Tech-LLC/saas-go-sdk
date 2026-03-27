package auth

import "context"

// GetSettings returns the public auth configuration for the project.
func (s *Service) GetSettings(ctx context.Context) (*Settings, error) {
	result := &Settings{}
	err := s.transport.Request(ctx, "GET", "/auth/settings", nil, result)
	return result, err
}
