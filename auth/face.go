package auth

import "context"

// EnrollFace stores a user's face template.
//
// Set FaceToken (from a face-required login response) to enrol during sign-in;
// the call then also completes the sign-in and returns a session. Without a
// token the request must carry the user's own access token and enrols the
// already signed-in user.
//
// Consent must be true: face data is only stored on an explicit opt-in.
func (s *Service) EnrollFace(ctx context.Context, params *FaceEnrollParams) (*FaceEnrollResult, error) {
	result := &FaceEnrollResult{}
	err := s.transport.Request(ctx, "POST", "/auth/face/enroll", params, result)
	return result, err
}

// VerifyFace completes a sign-in by matching a live descriptor against the
// stored template.
func (s *Service) VerifyFace(ctx context.Context, params *FaceVerifyParams) (*AuthResult, error) {
	result := &AuthResult{}
	err := s.transport.Request(ctx, "POST", "/auth/face/verify", params, result)
	return result, err
}

// FaceStatus reports whether the signed-in user has a face template.
func (s *Service) FaceStatus(ctx context.Context) (*FaceStatusResult, error) {
	result := &FaceStatusResult{}
	err := s.transport.Request(ctx, "GET", "/auth/face/status", nil, result)
	return result, err
}

// DeleteFace removes the signed-in user's face template.
func (s *Service) DeleteFace(ctx context.Context) error {
	return s.transport.Request(ctx, "DELETE", "/auth/face", nil, nil)
}
