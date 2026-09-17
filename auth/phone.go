package auth

import "context"

// SendPhoneOTP sends a one-time code to a phone number.
//
// Returns a 409 error when the project has no SMS provider configured. That is
// a configuration state, not a failure — registration in such a project does
// not require a code, so callers should skip the verification step rather than
// surface an error. Check Settings.PhoneOtpRequired to know in advance.
func (s *Service) SendPhoneOTP(ctx context.Context, params *PhoneOTPSendParams) (*PhoneOTPSendResult, error) {
	result := &PhoneOTPSendResult{}
	err := s.transport.Request(ctx, "POST", "/auth/phone/otp/send", params, result)
	return result, err
}

// VerifyPhoneOTP exchanges a delivered code for a short-lived token that proves
// ownership of the number. Pass the token to Register or PhonePasswordReset.
func (s *Service) VerifyPhoneOTP(ctx context.Context, params *PhoneOTPVerifyParams) (*PhoneOTPVerifyResult, error) {
	result := &PhoneOTPVerifyResult{}
	err := s.transport.Request(ctx, "POST", "/auth/phone/otp/verify", params, result)
	return result, err
}

// PhonePasswordReset sets a new password for the account owning a phone number,
// authorised by an OTP token issued for the "reset" purpose.
func (s *Service) PhonePasswordReset(ctx context.Context, params *PhonePasswordResetParams) (*PasswordResetVerifyResult, error) {
	result := &PasswordResetVerifyResult{}
	err := s.transport.Request(ctx, "POST", "/auth/password-reset/phone", params, result)
	return result, err
}
