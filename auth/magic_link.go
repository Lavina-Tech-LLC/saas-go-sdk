package auth

import "context"

// MagicLinkSend generates a magic-link token and emails it to the user.
func (s *Service) MagicLinkSend(ctx context.Context, params *MagicLinkSendParams) (*MagicLinkSendResult, error) {
	result := &MagicLinkSendResult{}
	err := s.transport.Request(ctx, "POST", "/auth/magic-link/send", params, result)
	return result, err
}

// MagicLinkVerify validates a magic-link token and returns auth tokens.
func (s *Service) MagicLinkVerify(ctx context.Context, params *MagicLinkVerifyParams) (*AuthResult, error) {
	result := &AuthResult{}
	err := s.transport.Request(ctx, "POST", "/auth/magic-link/verify", params, result)
	return result, err
}
