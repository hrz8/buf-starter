package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ServerConfig struct {
	Host           string `yaml:"host" validate:"required,hostname|ip"`
	Port           int    `yaml:"port" validate:"required,gte=1,lte=65535"`
	LogLevel       string `yaml:"logLevel" validate:"oneof=debug info warn error"`
	HTTPLogging    bool   `yaml:"httpLogging"`
	EnableCORS     bool   `yaml:"enableCORS"`
	ReadTimeout    int    `yaml:"readTimeout" validate:"gte=1"`
	WriteTimeout   int    `yaml:"writeTimeout" validate:"gte=1"`
	IdleTimeout    int    `yaml:"idleTimeout" validate:"gte=1"`
	CleanupTimeout int    `yaml:"cleanupTimeout" validate:"gte=1,lte=300"`
}

func (c *ServerConfig) setDefaults() {
	if c.Host == "" {
		c.Host = "localhost"
	}
	if c.Port == 0 {
		c.Port = 3100
	}
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}
	c.HTTPLogging = false
	c.EnableCORS = true
	if c.ReadTimeout == 0 {
		c.ReadTimeout = 15
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = 15
	}
	if c.IdleTimeout == 0 {
		c.IdleTimeout = 60
	}
	if c.CleanupTimeout == 0 {
		c.CleanupTimeout = 10
	}
}

type DatabaseConfig struct {
	URL            string `yaml:"url" validate:"omitempty,url"`
	MaxConnections int    `yaml:"maxConnections" validate:"gte=1,lte=100"`
	MaxIdleTime    int    `yaml:"maxIdleTime" validate:"gte=1"`
	ConnectTimeout int    `yaml:"connectTimeout" validate:"gte=1"`
}

func (c *DatabaseConfig) setDefaults() {
	if c.MaxConnections == 0 {
		c.MaxConnections = 25
	}
	if c.MaxIdleTime == 0 {
		c.MaxIdleTime = 300
	}
	if c.ConnectTimeout == 0 {
		c.ConnectTimeout = 10
	}
}

type SecurityConfig struct {
	AllowedOrigins    []string `yaml:"allowedOrigins" validate:"required,min=1,dive,required"`
	IAMEncryptionKey  string   `yaml:"iamEncryptionKey" validate:"required,len=44"` // base64-encoded 32-byte key = 44 chars
	JWTPrivateKeyPath string   `yaml:"jwtPrivateKeyPath" validate:"required"`
	JWTPublicKeyPath  string   `yaml:"jwtPublicKeyPath" validate:"required"`
	JWKSKid           string   `yaml:"jwksKid" validate:"required"`
}

func (c *SecurityConfig) setDefaults() {
	if len(c.AllowedOrigins) == 0 {
		c.AllowedOrigins = []string{"*"}
		return
	}
	for _, origin := range c.AllowedOrigins {
		if origin == "" || origin == "*" {
			c.AllowedOrigins = []string{"*"}
			return
		}
	}
}

type AuthConfig struct {
	Host               string `yaml:"host" validate:"required,hostname|ip"`
	Port               int    `yaml:"port" validate:"required,gte=1,lte=65535"`
	SessionSecret      string `yaml:"sessionSecret" validate:"required,min=32"`
	CodeExpiry         int    `yaml:"codeExpiry" validate:"gte=1"`
	AccessTokenExpiry  int    `yaml:"accessTokenExpiry" validate:"gte=1"`
	RefreshTokenExpiry int    `yaml:"refreshTokenExpiry" validate:"gte=1"`
}

func (c *AuthConfig) setDefaults() {
	if c.Host == "" {
		c.Host = "localhost"
	}
	if c.Port == 0 {
		c.Port = 3101
	}
	if c.CodeExpiry == 0 {
		c.CodeExpiry = 600 // 10 minutes
	}
	if c.AccessTokenExpiry == 0 {
		c.AccessTokenExpiry = 3600 // 1 hour
	}
	if c.RefreshTokenExpiry == 0 {
		c.RefreshTokenExpiry = 2592000 // 30 days
	}
}

type SuperadminConfig struct {
	Email string `yaml:"email" validate:"required,email"`
}

type OAuthProviderConfig struct {
	Provider     string `yaml:"provider" validate:"required,oneof=google github"`
	ClientID     string `yaml:"clientId" validate:"required"`
	ClientSecret string `yaml:"clientSecret" validate:"required"`
	RedirectURL  string `yaml:"redirectUrl" validate:"required,url"`
	Scopes       string `yaml:"scopes" validate:"required"`
	Enabled      bool   `yaml:"enabled"`
}

type SeederConfig struct {
	Superadmin     SuperadminConfig      `yaml:"superadmin" validate:"required"`
	OAuthProviders []OAuthProviderConfig `yaml:"oauthProviders" validate:"dive"`
}

// DashboardOAuthConfig contains OAuth configuration for the Default Dashboard client
type DashboardOAuthConfig struct {
	ExternalServer bool     `yaml:"externalServer"`
	Server         string   `yaml:"server" validate:"required,url"`
	Name           string   `yaml:"name" validate:"required"`
	ClientID       string   `yaml:"clientId" validate:"required,uuid"`
	ClientSecret   string   `yaml:"clientSecret" validate:"required,min=32"`
	RedirectURIs   []string `yaml:"redirectUris" validate:"required,min=1,dive,required,url"`
	PKCERequired   bool     `yaml:"pkceRequired"`
}

type AppConfig struct {
	Server         *ServerConfig         `yaml:"server" validate:"required"`
	Database       *DatabaseConfig       `yaml:"database" validate:"required"`
	Security       *SecurityConfig       `yaml:"security" validate:"required"`
	Auth           *AuthConfig           `yaml:"auth" validate:"required"`
	Seeder         *SeederConfig         `yaml:"seeder" validate:"required"`
	DashboardOAuth *DashboardOAuthConfig `yaml:"dashboardOauth" validate:"required"`
}

func (c *AppConfig) setDefaults() {
	c.Server.setDefaults()
	c.Database.setDefaults()
	c.Security.setDefaults()
	c.Auth.setDefaults()
}

func (c *AppConfig) Validate() error {
	validate := validator.New()

	_ = validate.RegisterValidation("allowed_origin", func(fl validator.FieldLevel) bool {
		origins := fl.Field().Interface().([]string)
		for _, origin := range origins {
			if origin == "*" || origin != "" {
				continue
			}
			return false
		}
		return true
	})

	if err := validate.Struct(c); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	return nil
}
