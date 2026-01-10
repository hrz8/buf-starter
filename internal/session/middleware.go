package session

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"slices"
)

// Middleware provides HTTP middleware for session management.
type Middleware struct {
	store       *Store
	loginPath   string
	excludePath []string
}

// NewMiddleware creates a new session middleware with the given store and configuration.
func NewMiddleware(store *Store, loginPath string, excludePath []string) *Middleware {
	return &Middleware{
		store:       store,
		loginPath:   loginPath,
		excludePath: excludePath,
	}
}

// LoadSession loads session data from the cookie into the request context.
func (m *Middleware) LoadSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := m.store.GetData(r)
		if err != nil {
			data = &Data{}
		}
		ctx := WithSession(r.Context(), data)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuth redirects unauthenticated users to the login page.
func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if slices.Contains(m.excludePath, r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		if !m.store.IsAuthenticated(r) {
			data, _ := m.store.GetData(r)
			if data == nil {
				data = &Data{}
			}
			data.OriginalURL = r.URL.String()
			_ = m.store.SetData(r, w, data)

			http.Redirect(w, r, m.loginPath, http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ValidateCSRF validates the CSRF token on POST requests.
func (m *Middleware) ValidateCSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			next.ServeHTTP(w, r)
			return
		}

		data := FromContext(r.Context())
		if data == nil {
			http.Error(w, "session not found", http.StatusForbidden)
			return
		}

		formToken := r.FormValue("csrf_token")
		if formToken == "" {
			formToken = r.Header.Get("X-CSRF-Token")
		}

		if formToken == "" || formToken != data.CSRFToken {
			http.Error(w, "invalid CSRF token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GenerateCSRFToken generates a cryptographically secure CSRF token.
func GenerateCSRFToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// GenerateState generates a cryptographically secure OAuth state parameter.
func GenerateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
