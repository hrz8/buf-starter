# Task T30: Go HTML Templates (Login, Consent, Error Pages)

**Story Reference:** US7-oauth-authorization-server.md
**Type:** Frontend (Go Templates)
**Priority:** Medium
**Estimated Effort:** 3-4 hours
**Prerequisites:** T26 (serve-auth CLI Command)

## Objective

Create HTML templates for the OAuth server pages using Go's html/template package. Templates should be clean, responsive, and consistent with Altalune branding.

## Acceptance Criteria

- [ ] Template loading and caching system
- [ ] Base layout with Tailwind CSS
- [ ] Login page with provider buttons
- [ ] Consent page with scope descriptions
- [ ] Error page with user-friendly messages
- [ ] Logout success page
- [ ] Responsive design (mobile-friendly)
- [ ] Proper HTML escaping for security

## Technical Requirements

### Template Package (`internal/authserver/templates/`)

#### Template Manager (`templates.go`)

```go
package templates

import (
    "embed"
    "html/template"
    "io"
    "sync"
)

//go:embed html/*.html
var templateFS embed.FS

var (
    templates *template.Template
    once      sync.Once
)

func Load() error {
    var loadErr error
    once.Do(func() {
        templates, loadErr = template.New("").
            Funcs(template.FuncMap{
                "join": strings.Join,
            }).
            ParseFS(templateFS, "html/*.html")
    })
    return loadErr
}

func Render(w io.Writer, name string, data interface{}) error {
    return templates.ExecuteTemplate(w, name, data)
}
```

#### Data Structures (`data.go`)

```go
package templates

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
    Name     string // "google", "github"
    Label    string // "Continue with Google"
    IconSVG  string // Inline SVG for provider logo
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
    Nonce               string
    CodeChallenge       string
    CodeChallengeMethod string
}

type ScopeInfo struct {
    Name        string
    Description string
    Icon        string // Icon name or SVG
}

type ErrorPageData struct {
    BaseData
    Error            string
    ErrorDescription string
    ShowBackToLogin  bool
}

type LogoutSuccessData struct {
    BaseData
    LoginURL string
}
```

### HTML Templates

#### Base Layout (`html/base.html`)

```html
{{define "base"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - Altalune</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <style>
        /* Custom styles for Altalune branding */
        .btn-primary {
            @apply bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-lg transition-colors;
        }
        .btn-secondary {
            @apply bg-gray-200 hover:bg-gray-300 text-gray-800 font-medium py-2 px-4 rounded-lg transition-colors;
        }
        .btn-provider {
            @apply flex items-center justify-center gap-3 w-full py-3 px-4 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors;
        }
    </style>
</head>
<body class="min-h-screen bg-gray-50">
    <div class="min-h-screen flex flex-col items-center justify-center p-4">
        {{template "content" .}}
    </div>
</body>
</html>
{{end}}
```

#### Login Page (`html/login.html`)

```html
{{define "content"}}
<div class="w-full max-w-md">
    <!-- Logo/Branding -->
    <div class="text-center mb-8">
        <h1 class="text-2xl font-bold text-gray-900">Sign in to Altalune</h1>
        <p class="text-gray-600 mt-2">Choose your authentication method</p>
    </div>

    <!-- Error Message -->
    {{if .ErrorMessage}}
    <div class="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
        <p class="text-red-700 text-sm">{{.ErrorMessage}}</p>
    </div>
    {{end}}

    <!-- Provider Buttons -->
    <div class="bg-white shadow-md rounded-lg p-6 space-y-4">
        {{range .Providers}}
        <a href="/login/{{.Name}}" class="btn-provider">
            {{.IconSVG}}
            <span class="text-gray-700">{{.Label}}</span>
        </a>
        {{end}}
    </div>

    <!-- Footer -->
    <p class="text-center text-sm text-gray-500 mt-6">
        By signing in, you agree to our Terms of Service and Privacy Policy
    </p>
</div>
{{end}}

{{template "base" .}}
```

#### Consent Page (`html/consent.html`)

```html
{{define "content"}}
<div class="w-full max-w-md">
    <!-- Header -->
    <div class="text-center mb-6">
        <h1 class="text-2xl font-bold text-gray-900">Authorize {{.ClientName}}</h1>
        <p class="text-gray-600 mt-2">This application wants to access your account</p>
    </div>

    <!-- Consent Form -->
    <div class="bg-white shadow-md rounded-lg p-6">
        <!-- Requested Permissions -->
        <div class="mb-6">
            <h2 class="text-sm font-medium text-gray-500 uppercase tracking-wide mb-3">
                Requested Permissions
            </h2>
            <ul class="space-y-3">
                {{range .Scopes}}
                <li class="flex items-start gap-3">
                    <svg class="w-5 h-5 text-green-500 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                    </svg>
                    <div>
                        <p class="font-medium text-gray-900">{{.Name}}</p>
                        <p class="text-sm text-gray-500">{{.Description}}</p>
                    </div>
                </li>
                {{end}}
            </ul>
        </div>

        <!-- Action Buttons -->
        <form method="POST" action="/oauth/authorize">
            <input type="hidden" name="csrf_token" value="{{.CSRFToken}}">
            <input type="hidden" name="client_id" value="{{.ClientID}}">
            <input type="hidden" name="redirect_uri" value="{{.RedirectURI}}">
            <input type="hidden" name="scope" value="{{.Scope}}">
            <input type="hidden" name="state" value="{{.State}}">
            <input type="hidden" name="nonce" value="{{.Nonce}}">
            <input type="hidden" name="code_challenge" value="{{.CodeChallenge}}">
            <input type="hidden" name="code_challenge_method" value="{{.CodeChallengeMethod}}">

            <div class="flex gap-3">
                <button type="submit" name="decision" value="deny" class="btn-secondary flex-1">
                    Deny
                </button>
                <button type="submit" name="decision" value="allow" class="btn-primary flex-1">
                    Allow
                </button>
            </div>
        </form>
    </div>

    <!-- Learn More -->
    <p class="text-center text-sm text-gray-500 mt-6">
        <a href="#" class="text-blue-600 hover:underline">Learn more</a> about OAuth permissions
    </p>
</div>
{{end}}

{{template "base" .}}
```

#### Error Page (`html/error.html`)

```html
{{define "content"}}
<div class="w-full max-w-md text-center">
    <!-- Error Icon -->
    <div class="mb-6">
        <svg class="w-16 h-16 mx-auto text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path>
        </svg>
    </div>

    <!-- Error Message -->
    <div class="bg-white shadow-md rounded-lg p-6">
        <h1 class="text-xl font-bold text-gray-900 mb-2">
            {{if .Error}}
                {{.Error}}
            {{else}}
                Something went wrong
            {{end}}
        </h1>
        {{if .ErrorDescription}}
        <p class="text-gray-600">{{.ErrorDescription}}</p>
        {{end}}
    </div>

    <!-- Back Link -->
    {{if .ShowBackToLogin}}
    <a href="/login" class="inline-block mt-6 text-blue-600 hover:underline">
        Back to Login
    </a>
    {{end}}
</div>
{{end}}

{{template "base" .}}
```

#### Logout Success Page (`html/logout_success.html`)

```html
{{define "content"}}
<div class="w-full max-w-md text-center">
    <!-- Success Icon -->
    <div class="mb-6">
        <svg class="w-16 h-16 mx-auto text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
        </svg>
    </div>

    <!-- Message -->
    <div class="bg-white shadow-md rounded-lg p-6">
        <h1 class="text-xl font-bold text-gray-900 mb-2">You have been logged out</h1>
        <p class="text-gray-600">Your session has been ended successfully.</p>
    </div>

    <!-- Login Link -->
    <a href="{{.LoginURL}}" class="inline-block mt-6 btn-primary">
        Login again
    </a>
</div>
{{end}}

{{template "base" .}}
```

### Provider Icons

```go
// internal/authserver/templates/icons.go
package templates

const GoogleIconSVG = `<svg class="w-5 h-5" viewBox="0 0 24 24">
    <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
    <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
    <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
    <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
</svg>`

const GitHubIconSVG = `<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
    <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
</svg>`
```

## Files to Create

- `internal/authserver/templates/templates.go` - Template loading and rendering
- `internal/authserver/templates/data.go` - Data structures for templates
- `internal/authserver/templates/icons.go` - Provider icon SVGs
- `internal/authserver/templates/html/base.html` - Base layout
- `internal/authserver/templates/html/login.html` - Login page
- `internal/authserver/templates/html/consent.html` - Consent page
- `internal/authserver/templates/html/error.html` - Error page
- `internal/authserver/templates/html/logout_success.html` - Logout success

## Files to Modify

- `internal/authserver/handlers/*.go` - Use template rendering

## Testing Requirements

- Visual testing in browser
- Test responsive design on mobile viewport
- Test all error states
- Test template escaping (XSS prevention)

## Commands to Run

```bash
# Build application
make build

# Start auth server
./bin/app serve-auth -c config.yaml

# Test login page
open http://localhost:3101/login

# Test error page (simulate error)
open "http://localhost:3101/login?error=test_error"
```

## Validation Checklist

- [ ] Templates load without errors
- [ ] Login page displays provider buttons
- [ ] Consent page shows client name and scopes
- [ ] Error page shows error message
- [ ] Logout success page has login link
- [ ] Responsive on mobile devices
- [ ] No XSS vulnerabilities (proper escaping)

## Definition of Done

- [ ] Template package with loading and caching
- [ ] Base layout with Tailwind CSS
- [ ] Login page with Google and GitHub buttons
- [ ] Consent page with scope list and Allow/Deny buttons
- [ ] Error page with user-friendly messages
- [ ] Logout success page with login link
- [ ] All templates properly escape user data
- [ ] Responsive design works on mobile

## Dependencies

- T26: serve-auth server to serve templates
- Tailwind CSS via CDN (no build step required)
- Go embed for template files

## Risk Factors

- **Low Risk**: Standard Go templates
- **Low Risk**: Tailwind via CDN is simple and reliable

## Notes

- Using Tailwind CSS via CDN for simplicity (no build process)
- Templates are embedded in binary using Go 1.16+ embed
- Provider icons are inline SVG for simplicity
- CSRF token must be included in consent form
- All user-provided data must be escaped by template engine
- Consider adding dark mode support in future
