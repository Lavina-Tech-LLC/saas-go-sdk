package auth

import "context"

// Login authenticates an end-user. If the user has MFA enabled,
// LoginResult.MFARequired will be true and you must call LoginMFA.
func (s *Service) Login(ctx context.Context, params *LoginParams) (*LoginResult, error) {
	result := &LoginResult{}
	err := s.transport.Request(ctx, "POST", "/auth/login", params, result)
	return result, err
}

// LoginMFA completes login for users with MFA enabled.
func (s *Service) LoginMFA(ctx context.Context, params *LoginMFAParams) (*AuthResult, error) {
	result := &AuthResult{}
	err := s.transport.Request(ctx, "POST", "/auth/login/mfa", params, result)
	return result, err
}
