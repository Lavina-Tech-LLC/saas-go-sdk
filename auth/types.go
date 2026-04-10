package auth

import (
	"encoding/json"
	"time"
)

// --- User types ---

// User represents an authenticated end-user.
type User struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	ProjectID     string `json:"projectId"`
	EmailVerified bool   `json:"emailVerified"`
	Provider      string `json:"provider"`
	Source        string `json:"source,omitempty"` // "self" | "invite"
	Metadata      string `json:"metadata,omitempty"`
	MFAEnabled    bool   `json:"mfaEnabled,omitempty"`
}

// ParseMetadata unmarshals the user's Metadata JSON string into a map.
func (u *User) ParseMetadata() (map[string]interface{}, error) {
	if u.Metadata == "" {
		return map[string]interface{}{}, nil
	}
	var m map[string]interface{}
	err := json.Unmarshal([]byte(u.Metadata), &m)
	return m, err
}

// ParseMetadataTo unmarshals the user's Metadata JSON string into dest.
func (u *User) ParseMetadataTo(dest interface{}) error {
	if u.Metadata == "" {
		return nil
	}
	return json.Unmarshal([]byte(u.Metadata), dest)
}

// --- Auth result types ---

// AuthResult is returned by Register, LoginMFA, MagicLinkVerify, OAuthCallback.
type AuthResult struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	User         User   `json:"user"`
}

// LoginResult is returned by Login. Check MFARequired to determine the flow.
type LoginResult struct {
	// MFARequired is true when the user has MFA enabled.
	MFARequired bool   `json:"mfaRequired"`
	MFAToken    string `json:"mfaToken,omitempty"`

	// Populated when MFARequired is false.
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	User         *User  `json:"user,omitempty"`
}

// --- Request param types ---

// RegisterParams are the parameters for Register.
type RegisterParams struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	InviteToken string `json:"inviteToken,omitempty"` // optional: register via invite
}

// LoginParams are the parameters for Login.
type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginMFAParams are the parameters for LoginMFA.
type LoginMFAParams struct {
	MFAToken string `json:"mfaToken"`
	Code     string `json:"code"`
}

// RefreshParams are the parameters for Refresh.
type RefreshParams struct {
	RefreshToken string `json:"refreshToken"`
}

// RefreshResult is returned by Refresh.
type RefreshResult struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// LogoutParams are the parameters for Logout.
type LogoutParams struct {
	RefreshToken string `json:"refreshToken"`
}

// LogoutResult is returned by Logout.
type LogoutResult struct {
	LoggedOut bool `json:"loggedOut"`
}

// VerifyTokenParams are the parameters for VerifyToken.
type VerifyTokenParams struct {
	Token string `json:"token"`
}

// TokenClaims is returned by VerifyToken.
type TokenClaims struct {
	Valid     bool      `json:"valid"`
	UserID    string    `json:"userId"`
	Email     string    `json:"email"`
	ProjectID string    `json:"projectId"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// MagicLinkSendParams are the parameters for MagicLinkSend.
type MagicLinkSendParams struct {
	Email       string `json:"email"`
	RedirectURL string `json:"redirectUrl"`
}

// MagicLinkSendResult is returned by MagicLinkSend.
type MagicLinkSendResult struct {
	Sent bool `json:"sent"`
}

// MagicLinkVerifyParams are the parameters for MagicLinkVerify.
type MagicLinkVerifyParams struct {
	Token string `json:"token"`
}

// PasswordResetSendParams are the parameters for PasswordResetSend.
type PasswordResetSendParams struct {
	Email       string `json:"email"`
	RedirectURL string `json:"redirectUrl"`
}

// PasswordResetSendResult is returned by PasswordResetSend.
type PasswordResetSendResult struct {
	Sent bool `json:"sent"`
}

// PasswordResetVerifyParams are the parameters for PasswordResetVerify.
type PasswordResetVerifyParams struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

// PasswordResetVerifyResult is returned by PasswordResetVerify.
type PasswordResetVerifyResult struct {
	Reset bool `json:"reset"`
}

// UpdateProfileParams are the parameters for UpdateMe.
type UpdateProfileParams struct {
	Metadata string `json:"metadata"`
}

// --- MFA types ---

// MFASetupResult is returned by MFASetup.
type MFASetupResult struct {
	Secret string `json:"secret"`
	URI    string `json:"uri"`
}

// MFAVerifyParams are the parameters for MFAVerify.
type MFAVerifyParams struct {
	Code string `json:"code"`
}

// MFAVerifyResult is returned by MFAVerify.
type MFAVerifyResult struct {
	BackupCodes []string `json:"backupCodes"`
}

// MFADisableParams are the parameters for MFADisable.
type MFADisableParams struct {
	Code string `json:"code"`
}

// MFADisableResult is returned by MFADisable.
type MFADisableResult struct {
	Disabled bool `json:"disabled"`
}

// --- Org types ---

// Org represents an organisation.
type Org struct {
	ID        string `json:"id"`
	ProjectID string `json:"projectId"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	AvatarURL string `json:"avatarUrl,omitempty"`
	Metadata  string `json:"metadata,omitempty"`
}

// ParseMetadata unmarshals the org's Metadata JSON string into a map.
func (o *Org) ParseMetadata() (map[string]interface{}, error) {
	if o.Metadata == "" {
		return map[string]interface{}{}, nil
	}
	var m map[string]interface{}
	err := json.Unmarshal([]byte(o.Metadata), &m)
	return m, err
}

// ParseMetadataTo unmarshals the org's Metadata JSON string into dest.
func (o *Org) ParseMetadataTo(dest interface{}) error {
	if o.Metadata == "" {
		return nil
	}
	return json.Unmarshal([]byte(o.Metadata), dest)
}

// CreateOrgParams are the parameters for CreateOrg.
type CreateOrgParams struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// UpdateOrgParams are the parameters for UpdateOrg.
type UpdateOrgParams struct {
	Name      string `json:"name,omitempty"`
	AvatarURL string `json:"avatarUrl,omitempty"`
}

// DeleteResult is returned by delete operations.
type DeleteResult struct {
	Deleted bool `json:"deleted"`
}

// --- Member types ---

// Member represents an organisation member.
type Member struct {
	UserID   string `json:"userId"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	RoleID   string `json:"roleId,omitempty"`
	RoleName string `json:"roleName,omitempty"`
}

// UpdateMemberRoleParams are the parameters for UpdateMemberRole.
type UpdateMemberRoleParams struct {
	Role   string `json:"role,omitempty"`   // backward compat: system role key
	RoleID string `json:"roleId,omitempty"` // preferred: UUID of AuthRole
}

// UpdateResult is returned by update operations.
type UpdateResult struct {
	Updated bool `json:"updated"`
}

// RemoveResult is returned by remove operations.
type RemoveResult struct {
	Removed bool `json:"removed"`
}

// --- Invite types ---

// SendInviteParams are the parameters for SendInvite.
type SendInviteParams struct {
	Email  string `json:"email"`
	Role   string `json:"role,omitempty"`   // backward compat: system role key
	RoleID string `json:"roleId,omitempty"` // preferred: UUID of AuthRole
}

// Invite represents an organisation invite.
type Invite struct {
	InviteID  string    `json:"inviteId"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	RoleID    string    `json:"roleId,omitempty"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// AcceptInviteResult is returned by AcceptInvite.
type AcceptInviteResult struct {
	OrgID  string `json:"orgId"`
	Role   string `json:"role"`
	RoleID string `json:"roleId,omitempty"`
}

// --- Invite Link types ---

// CreateInviteLinkParams are the parameters for CreateInviteLink.
type CreateInviteLinkParams struct {
	Role    string `json:"role,omitempty"`
	RoleID  string `json:"roleId,omitempty"`
	MaxUses int    `json:"maxUses,omitempty"`
}

// InviteLink represents a reusable invite link.
type InviteLink struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Role      string    `json:"role"`
	RoleID    string    `json:"roleId,omitempty"`
	RoleName  string    `json:"roleName,omitempty"`
	MaxUses   int       `json:"maxUses"`
	UseCount  int       `json:"useCount"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// InviteLinkInfo is the public info about an invite link (no org ID).
type InviteLinkInfo struct {
	OrgName      string    `json:"orgName"`
	OrgAvatarUrl string    `json:"orgAvatarUrl,omitempty"`
	Role         string    `json:"role"`
	RoleName     string    `json:"roleName,omitempty"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

// UseInviteLinkResult is returned by UseInviteLink.
type UseInviteLinkResult struct {
	OrgID   string `json:"orgId"`
	OrgName string `json:"orgName"`
	Role    string `json:"role"`
	RoleID  string `json:"roleId,omitempty"`
}

// --- OAuth types ---

// OAuthInitiateParams are the parameters for OAuthInitiate.
type OAuthInitiateParams struct {
	Provider    string // Path param
	RedirectURI string // Query param
}

// OAuthInitiateResult is returned by OAuthInitiate.
type OAuthInitiateResult struct {
	AuthURL string `json:"authUrl"`
	State   string `json:"state"`
}

// OAuthCallbackParams are the parameters for OAuthCallback.
type OAuthCallbackParams struct {
	Provider string `json:"-"` // Path param
	Code     string `json:"code"`
	State    string `json:"state"`
}

// --- Settings types ---

// Settings represents the public auth configuration for a project.
type Settings struct {
	GoogleEnabled     bool   `json:"googleEnabled"`
	GitHubEnabled     bool   `json:"githubEnabled"`
	EmailEnabled      bool   `json:"emailEnabled"`
	MFAEnforced       bool   `json:"mfaEnforced"`
	PasswordMinLength int    `json:"passwordMinLength"`
	EmailVerification bool   `json:"emailVerification"`
	OrgCreationPolicy string `json:"orgCreationPolicy"` // "anyone" | "self_registered_only"
}

// --- Role types ---

// Role represents a project-scoped role definition.
type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	Description string `json:"description,omitempty"`
	IsSystem    bool   `json:"isSystem"`
}

// --- CSRF types ---

// CSRFTokenResult is returned by GetCSRFToken.
type CSRFTokenResult struct {
	CSRFToken string `json:"csrfToken"`
}
