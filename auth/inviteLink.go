package auth

import "context"

// CreateInviteLink creates a reusable invite link for an organisation.
func (s *Service) CreateInviteLink(ctx context.Context, accessToken string, orgID string, params *CreateInviteLinkParams) (*InviteLink, error) {
	result := &InviteLink{}
	err := s.transport.RequestWithHeaders(ctx, "POST", "/auth/orgs/"+orgID+"/invite-links", params, result, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	return result, err
}

// ListInviteLinks returns all active invite links for an organisation.
func (s *Service) ListInviteLinks(ctx context.Context, accessToken string, orgID string) ([]InviteLink, error) {
	var result []InviteLink
	err := s.transport.RequestWithHeaders(ctx, "GET", "/auth/orgs/"+orgID+"/invite-links", nil, &result, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	return result, err
}

// RevokeInviteLink deletes an invite link.
func (s *Service) RevokeInviteLink(ctx context.Context, accessToken string, orgID string, linkID string) error {
	err := s.transport.RequestWithHeaders(ctx, "DELETE", "/auth/orgs/"+orgID+"/invite-links/"+linkID, nil, nil, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	return err
}

// GetInviteLinkInfo returns public info about an invite link (no auth required).
func (s *Service) GetInviteLinkInfo(ctx context.Context, code string) (*InviteLinkInfo, error) {
	result := &InviteLinkInfo{}
	err := s.transport.Request(ctx, "GET", "/auth/invite-links/"+code+"/info", nil, result)
	return result, err
}

// UseInviteLink uses an invite link to join the associated organisation.
func (s *Service) UseInviteLink(ctx context.Context, accessToken string, code string) (*UseInviteLinkResult, error) {
	result := &UseInviteLinkResult{}
	err := s.transport.RequestWithHeaders(ctx, "POST", "/auth/invite-links/"+code+"/use", nil, result, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	return result, err
}
