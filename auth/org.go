package auth

import "context"

// ListOrgs returns all orgs the authenticated user belongs to.
func (s *Service) ListOrgs(ctx context.Context, accessToken string) ([]Org, error) {
	var result []Org
	err := s.transport.RequestWithHeaders(ctx, "GET", "/auth/orgs", nil, &result, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	return result, err
}

// CreateOrg creates a new organisation and makes the caller its owner.
func (s *Service) CreateOrg(ctx context.Context, accessToken string, params *CreateOrgParams) (*Org, error) {
	result := &Org{}
	err := s.transport.RequestWithHeaders(ctx, "POST", "/auth/orgs", params, result, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	return result, err
}

// GetOrg returns a single organisation by ID.
func (s *Service) GetOrg(ctx context.Context, accessToken string, orgID string) (*Org, error) {
	result := &Org{}
	err := s.transport.RequestWithHeaders(ctx, "GET", "/auth/orgs/"+orgID, nil, result, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	return result, err
}

// UpdateOrg updates an organisation's name and/or avatar.
func (s *Service) UpdateOrg(ctx context.Context, accessToken string, orgID string, params *UpdateOrgParams) (*Org, error) {
	result := &Org{}
	err := s.transport.RequestWithHeaders(ctx, "PATCH", "/auth/orgs/"+orgID, params, result, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	return result, err
}

// DeleteOrg deletes an organisation. Caller must be the owner.
func (s *Service) DeleteOrg(ctx context.Context, accessToken string, orgID string) (*DeleteResult, error) {
	result := &DeleteResult{}
	err := s.transport.RequestWithHeaders(ctx, "DELETE", "/auth/orgs/"+orgID, nil, result, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	return result, err
}
