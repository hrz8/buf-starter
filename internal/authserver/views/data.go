package views

type BaseData struct {
	Title   string
	Message string
}

type LoginPageData struct {
	BaseData
	Providers    []Provider
	ErrorMessage string
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
