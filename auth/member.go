package auth

import "context"

// ListMembers returns all members of an organisation.
func (s *Service) ListMembers(ctx context.Context, accessToken string, orgID string) ([]Member, error) {
	var result []Member
	err := s.transport.RequestWithHeaders(ctx, "GET", "/auth/orgs/"+orgID+"/members", nil, &result, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	return result, err
}

// UpdateMemberRole changes a member's role within an organisation.
func (s *Service) UpdateMemberRole(ctx context.Context, accessToken string, orgID string, userID string, params *UpdateMemberRoleParams) (*UpdateResult, error) {
	result := &UpdateResult{}
	err := s.transport.RequestWithHeaders(ctx, "PATCH", "/auth/orgs/"+orgID+"/members/"+userID, params, result, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	return result, err
}

// RemoveMember removes a user from an organisation.
func (s *Service) RemoveMember(ctx context.Context, accessToken string, orgID string, userID string) (*RemoveResult, error) {
	result := &RemoveResult{}
	err := s.transport.RequestWithHeaders(ctx, "DELETE", "/auth/orgs/"+orgID+"/members/"+userID, nil, result, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	return result, err
}
