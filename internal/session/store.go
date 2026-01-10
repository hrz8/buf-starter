package session

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
)

const (
	CookieName = "altalune_auth"

	keyUserID          = "user_id"
	keyAuthenticatedAt = "authenticated_at"
	keyOAuthState      = "oauth_state"
	keyOriginalURL     = "original_url"
	keyCSRFToken       = "csrf_token"
)

type ctxKey string

const sessionCtxKey ctxKey = "session"

// Data holds session information for authenticated users.
type Data struct {
	UserID          int64
	AuthenticatedAt time.Time
	OAuthState      string
	OriginalURL     string
	CSRFToken       string
}

// Store wraps gorilla/sessions for cookie-based session management.
type Store struct {
	store *sessions.CookieStore
}

// NewStore creates a new session store with the given secret and options.
func NewStore(secret string, secure bool, maxAge int) *Store {
	store := sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
	return &Store{store: store}
}

// Get retrieves the session from the request cookie.
func (s *Store) Get(r *http.Request) (*sessions.Session, error) {
	return s.store.Get(r, CookieName)
}

// Save persists the session to the response cookie.
func (s *Store) Save(r *http.Request, w http.ResponseWriter, sess *sessions.Session) error {
	return s.store.Save(r, w, sess)
}

// GetData extracts session data from the request.
func (s *Store) GetData(r *http.Request) (*Data, error) {
	sess, err := s.Get(r)
	if err != nil {
		return nil, err
	}

	data := &Data{}

	if v, ok := sess.Values[keyUserID].(int64); ok {
		data.UserID = v
	}
	if v, ok := sess.Values[keyAuthenticatedAt].(int64); ok {
		data.AuthenticatedAt = time.Unix(v, 0)
	}
	if v, ok := sess.Values[keyOAuthState].(string); ok {
		data.OAuthState = v
	}
	if v, ok := sess.Values[keyOriginalURL].(string); ok {
		data.OriginalURL = v
	}
	if v, ok := sess.Values[keyCSRFToken].(string); ok {
		data.CSRFToken = v
	}

	return data, nil
}

// SetData stores session data in the response cookie.
func (s *Store) SetData(r *http.Request, w http.ResponseWriter, data *Data) error {
	sess, err := s.Get(r)
	if err != nil {
		return err
	}

	sess.Values[keyUserID] = data.UserID
	sess.Values[keyAuthenticatedAt] = data.AuthenticatedAt.Unix()
	sess.Values[keyOAuthState] = data.OAuthState
	sess.Values[keyOriginalURL] = data.OriginalURL
	sess.Values[keyCSRFToken] = data.CSRFToken

	return s.Save(r, w, sess)
}

// Clear removes all session data and invalidates the cookie.
func (s *Store) Clear(r *http.Request, w http.ResponseWriter) error {
	sess, err := s.Get(r)
	if err != nil {
		return err
	}

	sess.Options.MaxAge = -1
	for key := range sess.Values {
		delete(sess.Values, key)
	}

	return s.Save(r, w, sess)
}

// IsAuthenticated returns true if the request has a valid authenticated session.
func (s *Store) IsAuthenticated(r *http.Request) bool {
	data, err := s.GetData(r)
	if err != nil {
		return false
	}
	return data.UserID > 0
}

// WithSession returns a new context with session data attached.
func WithSession(ctx context.Context, data *Data) context.Context {
	return context.WithValue(ctx, sessionCtxKey, data)
}

// FromContext extracts session data from the context.
func FromContext(ctx context.Context) *Data {
	if data, ok := ctx.Value(sessionCtxKey).(*Data); ok {
		return data
	}
	return nil
}
