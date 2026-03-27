package auth

import (
	"context"
	"net/url"
)

// OAuthInitiate generates an authorization URL for the specified provider.
func (s *Service) OAuthInitiate(ctx context.Context, params *OAuthInitiateParams) (*OAuthInitiateResult, error) {
	result := &OAuthInitiateResult{}
	path := "/auth/oauth/" + params.Provider + "?redirect_uri=" + url.QueryEscape(params.RedirectURI)
	err := s.transport.Request(ctx, "GET", path, nil, result)
	return result, err
}

// OAuthCallback exchanges the authorization code for tokens.
func (s *Service) OAuthCallback(ctx context.Context, params *OAuthCallbackParams) (*AuthResult, error) {
	result := &AuthResult{}
	path := "/auth/oauth/" + params.Provider + "/callback"
	err := s.transport.Request(ctx, "POST", path, params, result)
	return result, err
}
