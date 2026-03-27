package auth

import "context"

// MFASetup generates a TOTP secret. MFA is not enabled until MFAVerify succeeds.
func (s *Service) MFASetup(ctx context.Context, accessToken string) (*MFASetupResult, error) {
	result := &MFASetupResult{}
	err := s.transport.RequestWithHeaders(ctx, "POST", "/auth/mfa/setup", nil, result, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	return result, err
}

// MFAVerify validates a TOTP code, enables MFA, and returns backup codes.
func (s *Service) MFAVerify(ctx context.Context, accessToken string, params *MFAVerifyParams) (*MFAVerifyResult, error) {
	result := &MFAVerifyResult{}
	err := s.transport.RequestWithHeaders(ctx, "POST", "/auth/mfa/verify", params, result, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	return result, err
}

// MFADisable disables MFA after validating a TOTP code.
func (s *Service) MFADisable(ctx context.Context, accessToken string, params *MFADisableParams) (*MFADisableResult, error) {
	result := &MFADisableResult{}
	err := s.transport.RequestWithHeaders(ctx, "POST", "/auth/mfa/disable", params, result, map[string]string{
		"Authorization": "Bearer " + accessToken,
	})
	return result, err
}
