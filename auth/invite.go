package auth

import "context"

// SendInvite creates an invite to join an organisation.
func (s *Service) SendInvite(ctx context.Context, accessToken string, orgID string, params *SendInviteParams) (*Invite, error) {
	result := &Invite{}
	err := s.transport.RequestWithHeaders(ctx, "POST", "/auth/orgs/"+orgID+"/invites", params, result, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	return result, err
}

// AcceptInvite accepts a pending invite by its raw token.
func (s *Service) AcceptInvite(ctx context.Context, accessToken string, token string) (*AcceptInviteResult, error) {
	result := &AcceptInviteResult{}
	err := s.transport.RequestWithHeaders(ctx, "POST", "/auth/invites/"+token+"/accept", nil, result, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	return result, err
}
