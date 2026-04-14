package auth

import "context"

// ListMembers returns all members of an organisation. The backend restricts
// this to owner/admin callers; non-privileged users will receive a 403. Use
// GetMyMembership to resolve the caller's own role.
func (s *Service) ListMembers(ctx context.Context, accessToken string, orgID string) ([]Member, error) {
	var result []Member
	err := s.transport.RequestWithHeaders(ctx, "GET", "/auth/orgs/"+orgID+"/members", nil, &result, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	return result, err
}

// GetMyMembership returns the authenticated user's own membership in the given
// organisation, including the primary role and every role assigned to them.
// Any member can call this — there is no owner/admin gate. Returns an *Error
// with Code 404 if the user is not a member of the org.
func (s *Service) GetMyMembership(ctx context.Context, accessToken string, orgID string) (*Member, error) {
	result := &Member{}
	err := s.transport.RequestWithHeaders(ctx, "GET", "/auth/orgs/"+orgID+"/me", nil, result, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
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
