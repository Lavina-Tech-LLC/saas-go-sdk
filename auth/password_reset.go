package auth

import "context"

// PasswordResetSend generates a password-reset token and emails it.
func (s *Service) PasswordResetSend(ctx context.Context, params *PasswordResetSendParams) (*PasswordResetSendResult, error) {
	result := &PasswordResetSendResult{}
	err := s.transport.Request(ctx, "POST", "/auth/password-reset/send", params, result)
	return result, err
}

// PasswordResetVerify validates a reset token and sets a new password.
func (s *Service) PasswordResetVerify(ctx context.Context, params *PasswordResetVerifyParams) (*PasswordResetVerifyResult, error) {
	result := &PasswordResetVerifyResult{}
	err := s.transport.Request(ctx, "POST", "/auth/password-reset/verify", params, result)
	return result, err
}
