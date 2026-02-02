package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"
)

// JWKSResponse represents the JWKS endpoint response.
type JWKSResponse struct {
	Keys []JWK `json:"keys"`
}

// JWK represents a JSON Web Key.
type JWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// JWKSCache manages caching of JWKS keys.
type JWKSCache struct {
	mu             sync.RWMutex
	keys           map[string]*rsa.PublicKey // kid -> public key
	lastFetch      time.Time
	ttl            time.Duration
	refreshCount   int
	lastRefreshMin time.Time
	refreshLimit   int
}

// NewJWKSCache creates a new JWKS cache.
func NewJWKSCache(ttl time.Duration, refreshLimit int) *JWKSCache {
	return &JWKSCache{
		keys:         make(map[string]*rsa.PublicKey),
		ttl:          ttl,
		refreshLimit: refreshLimit,
	}
}

// GetKey returns the public key for the given kid.
func (c *JWKSCache) GetKey(kid string) (*rsa.PublicKey, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	key, ok := c.keys[kid]
	return key, ok
}

// SetKeys stores the fetched keys.
func (c *JWKSCache) SetKeys(keys map[string]*rsa.PublicKey) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.keys = keys
	c.lastFetch = time.Now()
}

// IsExpired checks if the cache has expired.
func (c *JWKSCache) IsExpired() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return time.Since(c.lastFetch) > c.ttl
}

// CanRefresh checks if refresh is allowed (rate limiting).
func (c *JWKSCache) CanRefresh() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	if now.Sub(c.lastRefreshMin) > time.Minute {
		c.refreshCount = 0
		c.lastRefreshMin = now
	}

	if c.refreshCount >= c.refreshLimit {
		return false
	}

	c.refreshCount++
	return true
}

// JWKSFetcher fetches and parses JWKS from the auth server.
type JWKSFetcher struct {
	url        string
	httpClient *http.Client
	cache      *JWKSCache
}

// NewJWKSFetcher creates a new JWKS fetcher.
func NewJWKSFetcher(url string, cacheTTL time.Duration, refreshLimit int) *JWKSFetcher {
	return &JWKSFetcher{
		url:        url,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		cache:      NewJWKSCache(cacheTTL, refreshLimit),
	}
}

// GetPublicKey returns the public key for the given kid.
func (f *JWKSFetcher) GetPublicKey(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	// Try cache first
	if key, ok := f.cache.GetKey(kid); ok && !f.cache.IsExpired() {
		return key, nil
	}

	// Fetch fresh keys
	if err := f.Refresh(ctx); err != nil {
		return nil, err
	}

	key, ok := f.cache.GetKey(kid)
	if !ok {
		return nil, fmt.Errorf("key not found: %s", kid)
	}

	return key, nil
}

// Refresh fetches keys from JWKS endpoint.
func (f *JWKSFetcher) Refresh(ctx context.Context) error {
	if !f.cache.CanRefresh() {
		return fmt.Errorf("JWKS refresh rate limit exceeded")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", f.url, nil)
	if err != nil {
		return err
	}

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("JWKS fetch failed: status %d", resp.StatusCode)
	}

	var jwks JWKSResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return err
	}

	keys := make(map[string]*rsa.PublicKey)
	for _, jwk := range jwks.Keys {
		if jwk.Kty != "RSA" {
			continue
		}

		pubKey, err := parseRSAPublicKey(jwk)
		if err != nil {
			continue
		}

		keys[jwk.Kid] = pubKey
	}

	f.cache.SetKeys(keys)
	return nil
}

// ForceRefresh forces a refresh on validation failure (key rotation).
func (f *JWKSFetcher) ForceRefresh(ctx context.Context) error {
	return f.Refresh(ctx)
}

func parseRSAPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, err
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, err
	}

	n := new(big.Int).SetBytes(nBytes)
	e := int(new(big.Int).SetBytes(eBytes).Int64())

	return &rsa.PublicKey{N: n, E: e}, nil
}
