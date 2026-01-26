package views

// BrandingData contains branding information for templates.
type BrandingData struct {
	Name string // Auth server branding name
}

type BaseData struct {
	Title    string
	Message  string
	Branding BrandingData
}

type LoginPageData struct {
	BaseData
	Providers    []Provider
	ErrorMessage string
	ClientName   string
}

type Provider struct {
	Name    string
	Label   string
	IconSVG string
}

type ConsentPageData struct {
	BaseData
	ClientName          string
	Scopes              []ScopeInfo
	CSRFToken           string
	ClientID            string
	RedirectURI         string
	Scope               string
	State               string
	Nonce               *string
	CodeChallenge       *string
	CodeChallengeMethod *string
}

type ScopeInfo struct {
	Name        string
	Description string
}

type ErrorPageData struct {
	BaseData
	Error            string
	ErrorDescription string
	ShowBackToLogin  bool
}

type ProfileData struct {
	BaseData
	User                       any
	Identities                 any
	Consents                   any
	ShowEmailVerificationAlert bool   // Show alert if email not verified
	UserEmail                  string // For resend verification link
	VerificationEmailSent      bool   // Show success message after resending
	VerificationEmailError     bool   // Show error message if resend failed
}

// EmailLoginPageData is the data structure for the email login page.
type EmailLoginPageData struct {
	BaseData
	Error string
}

// OTPPageData is the data structure for the OTP verification page.
type OTPPageData struct {
	BaseData
	Email      string // Masked email (e.g., j***n@example.com)
	Error      string
	ExpiryMins int
}

// VerifyEmailResultData is the data structure for the email verification result page.
type VerifyEmailResultData struct {
	BaseData
	Success bool
	Error   string
}

// PendingActivationData is the data structure for the pending activation page.
type PendingActivationData struct {
	BaseData
	UserEmail string // Email of the user awaiting activation
}

// EditProfileData is the data structure for the edit profile page.
type EditProfileData struct {
	BaseData
	User         any
	ErrorMessage string
	Success      bool
}
