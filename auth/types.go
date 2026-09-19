package auth

import (
	"encoding/json"
	"time"
)

// --- User types ---

// User represents an authenticated end-user. A user is identified by an
// e-mail, a phone number, or both; the unused field is empty.
type User struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Phone         string `json:"phone,omitempty"` // E.164
	ProjectID     string `json:"projectId"`
	EmailVerified bool   `json:"emailVerified"`
	PhoneVerified bool   `json:"phoneVerified,omitempty"`
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

// LoginResult is returned by Login. Check MFARequired and FaceRequired to
// determine which step, if any, still has to be completed.
type LoginResult struct {
	// MFARequired is true when the user has MFA enabled.
	MFARequired bool   `json:"mfaRequired"`
	MFAToken    string `json:"mfaToken,omitempty"`

	// FaceRequired is true when the project uses face control. Complete the
	// sign-in with VerifyFace, or with EnrollFace when FaceEnrolled is false.
	FaceRequired bool     `json:"faceRequired"`
	FaceToken    string   `json:"faceToken,omitempty"`
	FaceEnrolled bool     `json:"enrolled,omitempty"`
	FacePoses    []string `json:"poses,omitempty"`

	// Populated when neither step is pending.
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	User         *User  `json:"user,omitempty"`
}

// --- Request param types ---

// RegisterParams are the parameters for Register. Set exactly one of
// Email / Phone — whichever the project's auth settings allow.
type RegisterParams struct {
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"` // accepted in any format; normalised to E.164
	Password string `json:"password"`
	// OTPToken proves ownership of Phone, obtained from VerifyPhoneOTP. Only
	// required when Settings.PhoneOtpRequired is true.
	OTPToken string `json:"otpToken,omitempty"`
	// InviteToken registers the user via a per-email invite.
	//
	// Deprecated: use InviteCode, which also accepts reusable link codes.
	InviteToken string `json:"inviteToken,omitempty"`
	// InviteCode registers the user via an invite — either a per-email invite
	// token or a reusable invite link code. Registering this way attaches the
	// user to the inviting organisation instead of creating a new one, and
	// works even when self-service registration is disabled.
	InviteCode string `json:"inviteCode,omitempty"`
}

// LoginParams are the parameters for Login. Set exactly one of Email / Phone.
type LoginParams struct {
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"` // accepted in any format; normalised to E.164
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
	Phone     string    `json:"phone,omitempty"`
	Name      string    `json:"name"`
	ProjectID string    `json:"projectId"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// VerifyAPIKeyParams are the parameters for VerifyAPIKey.
type VerifyAPIKeyParams struct {
	APIKey string `json:"apiKey"`
}

// APIKeyClaims is returned by VerifyAPIKey. Valid reports whether the key
// was recognised, not revoked, and (if ExpiresAt is set) not expired. When
// Valid is false the remaining fields are zero-valued.
type APIKeyClaims struct {
	Valid     bool       `json:"valid"`
	UserID    string     `json:"userId"`
	Email     string     `json:"email"`
	OrgID     string     `json:"orgId"`
	ProjectID string     `json:"projectId"`
	Roles     []Role     `json:"roles"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
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

// --- Phone verification types ---

// PhoneOTPPurpose values accepted by SendPhoneOTP / VerifyPhoneOTP.
const (
	PhoneOTPPurposeRegister = "register"
	PhoneOTPPurposeReset    = "reset"
)

// PhoneOTPSendParams are the parameters for SendPhoneOTP.
type PhoneOTPSendParams struct {
	Phone string `json:"phone"`
	// Purpose defaults to "register" when empty.
	Purpose string `json:"purpose,omitempty"`
}

// PhoneOTPSendResult is returned by SendPhoneOTP.
type PhoneOTPSendResult struct {
	Sent bool `json:"sent"`
	// ResendAfterSeconds is how long to wait before requesting another code.
	ResendAfterSeconds int `json:"resendAfterSeconds"`
	// ExpiresInSeconds is how long the delivered code stays valid.
	ExpiresInSeconds int `json:"expiresInSeconds"`
}

// PhoneOTPVerifyParams are the parameters for VerifyPhoneOTP.
type PhoneOTPVerifyParams struct {
	Phone   string `json:"phone"`
	Code    string `json:"code"`
	Purpose string `json:"purpose,omitempty"`
}

// PhoneOTPVerifyResult is returned by VerifyPhoneOTP.
type PhoneOTPVerifyResult struct {
	Verified bool   `json:"verified"`
	OTPToken string `json:"otpToken"`
	// ExpiresInSeconds is how long OTPToken stays usable.
	ExpiresInSeconds int `json:"expiresInSeconds"`
}

// --- Face control types ---

// Face verification modes reported by Settings.FaceVerificationMode.
const (
	FaceModeOff      = "off"
	FaceModeOptional = "optional"
	FaceModeRequired = "required"
)

// FaceSample is one captured head pose. Only the descriptor is transmitted —
// camera frames stay in the browser.
type FaceSample struct {
	Pose       string    `json:"pose"`
	Descriptor []float64 `json:"descriptor"`
	Quality    float64   `json:"quality"`
}

// FaceEnrollParams are the parameters for EnrollFace.
type FaceEnrollParams struct {
	// FaceToken authorises enrolment during a sign-in. Leave empty to enrol the
	// caller identified by the request's access token.
	FaceToken string       `json:"faceToken,omitempty"`
	Model     string       `json:"model,omitempty"`
	Samples   []FaceSample `json:"samples"`
	// Consent must be true — face data is special-category personal data.
	Consent bool `json:"consent"`
}

// FaceEnrollResult is returned by EnrollFace. The session fields are populated
// only when the enrolment also completed a sign-in.
type FaceEnrollResult struct {
	Enrolled     bool   `json:"enrolled"`
	PoseCount    int    `json:"poseCount,omitempty"`
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	User         *User  `json:"user,omitempty"`
}

// FaceVerifyParams are the parameters for VerifyFace.
type FaceVerifyParams struct {
	FaceToken  string    `json:"faceToken"`
	Descriptor []float64 `json:"descriptor"`
}

// FaceStatusResult is returned by FaceStatus.
type FaceStatusResult struct {
	Mode           string     `json:"mode"`
	Enrolled       bool       `json:"enrolled"`
	EnrolledAt     *time.Time `json:"enrolledAt,omitempty"`
	PoseCount      int        `json:"poseCount,omitempty"`
	LastVerifiedAt *time.Time `json:"lastVerifiedAt,omitempty"`
	Poses          []string   `json:"poses,omitempty"`
}

// PhonePasswordResetParams are the parameters for PhonePasswordReset.
type PhonePasswordResetParams struct {
	Phone       string `json:"phone"`
	OTPToken    string `json:"otpToken"`
	NewPassword string `json:"newPassword"`
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

// MemberRole is one element of Member.Roles — represents a single role
// assigned to the member (either a system role or a project-defined custom
// role).
type MemberRole struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

// Member represents an organisation member.
type Member struct {
	UserID   string       `json:"userId"`
	Email    string       `json:"email"`
	Role     string       `json:"role"`
	RoleID   string       `json:"roleId,omitempty"`
	RoleName string       `json:"roleName,omitempty"`
	Roles    []MemberRole `json:"roles,omitempty"`
}

// HasAnyRole reports whether the member matches any of the supplied role
// identifiers. Each identifier is compared against the primary role key,
// primary role ID, and every key/ID in the Roles array. Empty identifiers
// never match.
func (m *Member) HasAnyRole(allowed ...string) bool {
	for _, a := range allowed {
		if a == "" {
			continue
		}
		if a == m.Role || a == m.RoleID {
			return true
		}
		for _, r := range m.Roles {
			if a == r.Key || a == r.ID {
				return true
			}
		}
	}
	return false
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

// SendInviteParams are the parameters for SendInvite. Set exactly one
// identifier — Email or Phone — and it has to be one the project's sign-in
// methods accept. Re-inviting the same person replaces their pending invite
// with a fresh link.
type SendInviteParams struct {
	Email  string `json:"email,omitempty"`
	Phone  string `json:"phone,omitempty"`  // accepted in any format; normalised to E.164
	Role   string `json:"role,omitempty"`   // backward compat: system role key
	RoleID string `json:"roleId,omitempty"` // preferred: UUID of AuthRole
}

// Invite represents an organisation invite. Nothing is delivered by the
// platform: pass URL on to the invitee however you like.
type Invite struct {
	InviteID  string    `json:"inviteId"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone,omitempty"`
	URL       string    `json:"url,omitempty"`
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
	// ExpiresInHours overrides the project default lifetime. Point it at 0 for
	// a link that never expires; leave nil to use the project default.
	ExpiresInHours *int `json:"expiresInHours,omitempty"`
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
	// InviteCode keeps an OAuth sign-up on the invite path: the user is added
	// to the inviting organisation and no personal organisation is created.
	InviteCode string `json:"inviteCode,omitempty"`
}

// --- Settings types ---

// Settings represents the public auth configuration for a project.
type Settings struct {
	GoogleEnabled bool `json:"googleEnabled"`
	GitHubEnabled bool `json:"githubEnabled"`
	// EmailEnabled is the legacy name of EmailAuthEnabled.
	//
	// Deprecated: use EmailAuthEnabled.
	EmailEnabled bool `json:"emailEnabled"`
	// EmailAuthEnabled reports whether an e-mail address is accepted as the
	// login identifier; PhoneAuthEnabled does the same for phone numbers.
	EmailAuthEnabled bool `json:"emailAuthEnabled"`
	PhoneAuthEnabled bool `json:"phoneAuthEnabled"`
	// DefaultPhoneCountryCode ("+992") is applied to numbers typed without one.
	DefaultPhoneCountryCode string `json:"defaultPhoneCountryCode,omitempty"`
	// PhoneOtpRequired reports whether a phone sign-up must carry an
	// SMS-verified code. It is already resolved against the project's SMS
	// provider status, so a project asking for verification without a working
	// provider reports false. Never re-derive this client-side.
	PhoneOtpRequired bool `json:"phoneOtpRequired"`
	// FaceVerificationMode is "off", "optional" or "required".
	FaceVerificationMode string `json:"faceVerificationMode"`
	// FaceModelURL overrides where a browser SDK loads the face model from.
	FaceModelURL      string `json:"faceModelUrl,omitempty"`
	MFAEnforced       bool   `json:"mfaEnforced"`
	PasswordMinLength int    `json:"passwordMinLength"`
	EmailVerification bool   `json:"emailVerification"`
	OrgCreationPolicy string `json:"orgCreationPolicy"` // "anyone" | "self_registered_only"
	InviteLinkBaseURL string `json:"inviteLinkBaseUrl,omitempty"`
	// RegistrationEnabled reports whether self-service sign-up is open.
	// Registration through an invite works even when this is false.
	RegistrationEnabled bool `json:"registrationEnabled"`
	// CreateOrgOnRegistration reports whether a personal organisation is
	// created for a new self-registered user. When false, a user may legitimately
	// have no organisation until they accept an invite.
	CreateOrgOnRegistration bool `json:"createOrgOnRegistration"`
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
