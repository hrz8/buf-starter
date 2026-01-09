package oauth_seeder

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/hrz8/altalune"
	"github.com/hrz8/altalune/internal/shared/nanoid"
	"github.com/lib/pq"
)

// Seeder handles database seeding operations for OAuth infrastructure
type Seeder struct {
	db     *sql.DB
	config altalune.Config
	logger *slog.Logger
}

// NewSeeder creates a new Seeder instance
func NewSeeder(db *sql.DB, cfg altalune.Config) (*Seeder, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}

	logger := slog.Default()

	return &Seeder{
		db:     db,
		config: cfg,
		logger: logger,
	}, nil
}

// Seed executes all seeding operations in a transaction
func (s *Seeder) Seed(ctx context.Context) error {
	s.logger.Info("Starting database seeding...")

	// Start transaction for atomic seeding
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Execute seeding steps
	if err := s.seedSuperadmin(ctx, tx); err != nil {
		return fmt.Errorf("seed superadmin: %w", err)
	}

	if err := s.seedDefaultOAuthClient(ctx, tx); err != nil {
		return fmt.Errorf("seed default OAuth client: %w", err)
	}

	if err := s.seedOAuthProviders(ctx, tx); err != nil {
		return fmt.Errorf("seed OAuth providers: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	s.logger.Info("Database seeding completed successfully")
	return nil
}

// seedSuperadmin updates the superadmin user email, creates user identity, and project membership
// NOTE: Superadmin user (user_id=1) is created by SQL migration 20260105000001_seed_iam_data.sql
// This function only UPDATES the email address for environment-specific configuration
func (s *Seeder) seedSuperadmin(ctx context.Context, tx *sql.Tx) error {
	s.logger.Info("Updating superadmin user email...")

	email := s.config.GetSuperadminEmail()

	// Check if superadmin exists (user_id=1 from SQL migration)
	var userID int64
	err := tx.QueryRowContext(ctx, `
		SELECT id FROM altalune_users WHERE id = 1
	`).Scan(&userID)

	if err == sql.ErrNoRows {
		// Superadmin doesn't exist yet (fresh install, SQL migration hasn't run)
		s.logger.Info("Superadmin user (user_id=1) doesn't exist yet, skipping custom seeding")
		return nil
	}

	if err != nil {
		return fmt.Errorf("check superadmin existence: %w", err)
	}

	// Update superadmin email from config (allows different emails per environment)
	_, err = tx.ExecContext(ctx, `
		UPDATE altalune_users
		SET email = $1, updated_at = NOW()
		WHERE id = 1
	`, email)

	if err != nil {
		return fmt.Errorf("update superadmin email: %w", err)
	}

	s.logger.Info("Superadmin email updated successfully", "userId", userID, "email", email)

	// Ensure user identity exists
	if err := s.ensureSuperadminIdentity(ctx, tx, userID); err != nil {
		return err
	}

	// Ensure project membership exists
	if err := s.ensureSuperadminMembership(ctx, tx, userID); err != nil {
		return err
	}

	return nil
}

// createSuperadminIdentity creates a user identity for the superadmin user
func (s *Seeder) createSuperadminIdentity(ctx context.Context, tx *sql.Tx, userID int64) error {
	s.logger.Info("Creating superadmin user identity...")

	// Generate public_id at runtime using nanoid
	publicID, err := nanoid.GeneratePublicID()
	if err != nil {
		return fmt.Errorf("generate public_id: %w", err)
	}

	// Get user's name from database (fixed in SQL migration)
	var firstName, lastName string
	err = tx.QueryRowContext(ctx, `
		SELECT first_name, last_name FROM altalune_users WHERE id = $1
	`, userID).Scan(&firstName, &lastName)

	if err != nil {
		return fmt.Errorf("get user name: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO altalune_user_identities (
			public_id, user_id, provider, provider_user_id,
			email, first_name, last_name, oauth_client_id,
			last_login_at, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NULL, NOW(), NOW(), NOW())
	`, publicID, userID, "system", "superadmin",
		s.config.GetSuperadminEmail(), firstName, lastName)

	if err != nil {
		return fmt.Errorf("create user identity: %w", err)
	}

	s.logger.Info("Superadmin user identity created successfully")
	return nil
}

// ensureSuperadminIdentity ensures user identity exists for superadmin
func (s *Seeder) ensureSuperadminIdentity(ctx context.Context, tx *sql.Tx, userID int64) error {
	var count int
	err := tx.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM altalune_user_identities
		WHERE user_id = $1 AND provider = 'system'
	`, userID).Scan(&count)

	if err != nil {
		return fmt.Errorf("check user identity: %w", err)
	}

	if count > 0 {
		s.logger.Info("Superadmin user identity already exists, skipping")
		return nil
	}

	return s.createSuperadminIdentity(ctx, tx, userID)
}

// createSuperadminMembership creates a project membership for the superadmin user
func (s *Seeder) createSuperadminMembership(ctx context.Context, tx *sql.Tx, userID int64) error {
	s.logger.Info("Creating superadmin project membership...")

	// Get the first project (or create one if none exists)
	var projectID int64
	err := tx.QueryRowContext(ctx, `
		SELECT id FROM altalune_projects ORDER BY id LIMIT 1
	`).Scan(&projectID)

	if err == sql.ErrNoRows {
		// No project exists, this is expected on fresh installation
		s.logger.Info("No projects exist yet, skipping project membership creation")
		return nil
	}

	if err != nil {
		return fmt.Errorf("get first project: %w", err)
	}

	// Generate public_id at runtime using nanoid
	publicID, err := nanoid.GeneratePublicID()
	if err != nil {
		return fmt.Errorf("generate public_id: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO altalune_project_members (public_id, project_id, user_id, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
	`, publicID, projectID, userID, "owner")

	if err != nil {
		return fmt.Errorf("create project membership: %w", err)
	}

	s.logger.Info("Superadmin project membership created successfully", "projectId", projectID, "role", "owner")
	return nil
}

// ensureSuperadminMembership ensures project membership exists for superadmin
func (s *Seeder) ensureSuperadminMembership(ctx context.Context, tx *sql.Tx, userID int64) error {
	var count int
	err := tx.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM altalune_project_members WHERE user_id = $1 AND role = 'owner'
	`, userID).Scan(&count)

	if err != nil {
		return fmt.Errorf("check project membership: %w", err)
	}

	if count > 0 {
		s.logger.Info("Superadmin project membership already exists, skipping")
		return nil
	}

	return s.createSuperadminMembership(ctx, tx, userID)
}

// seedDefaultOAuthClient creates the default OAuth client for the dashboard
func (s *Seeder) seedDefaultOAuthClient(ctx context.Context, tx *sql.Tx) error {
	s.logger.Info("Checking default OAuth client...")

	clientID := s.config.GetDefaultOAuthClientID()
	clientUUID, err := uuid.Parse(clientID)
	if err != nil {
		return fmt.Errorf("parse client UUID: %w", err)
	}

	// Check if client already exists
	var existingID int64
	err = tx.QueryRowContext(ctx, `
		SELECT id FROM altalune_oauth_clients WHERE client_id = $1 LIMIT 1
	`, clientUUID).Scan(&existingID)

	if err == nil {
		s.logger.Info("Default OAuth client already exists, skipping", "clientId", clientID)
		return nil
	}

	if err != sql.ErrNoRows {
		return fmt.Errorf("check OAuth client existence: %w", err)
	}

	// Get the first project
	var projectID int64
	err = tx.QueryRowContext(ctx, `
		SELECT id FROM altalune_projects ORDER BY id LIMIT 1
	`).Scan(&projectID)

	if err == sql.ErrNoRows {
		s.logger.Info("No projects exist yet, skipping OAuth client creation")
		return nil
	}

	if err != nil {
		return fmt.Errorf("get first project: %w", err)
	}

	// Get client configuration from interface
	clientName := s.config.GetDefaultOAuthClientName()
	clientSecret := s.config.GetDefaultOAuthClientSecret()
	redirectURIs := s.config.GetDefaultOAuthClientRedirectURIs()
	pkceRequired := s.config.GetDefaultOAuthClientPKCERequired()

	// Hash the client secret
	s.logger.Info("Creating default OAuth client...", "name", clientName)
	secretHash, err := HashClientSecret(clientSecret)
	if err != nil {
		return fmt.Errorf("hash client secret: %w", err)
	}

	// Generate public_id at runtime using nanoid
	publicID, err := nanoid.GeneratePublicID()
	if err != nil {
		return fmt.Errorf("generate public_id: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO altalune_oauth_clients (
			project_id, public_id, name, client_id, client_secret_hash,
			redirect_uris, pkce_required, is_default, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
	`, projectID, publicID, clientName, clientUUID, secretHash,
		pq.Array(redirectURIs), pkceRequired, true)

	if err != nil {
		return fmt.Errorf("create OAuth client: %w", err)
	}

	s.logger.Info("Default OAuth client created successfully",
		"clientId", clientID,
		"name", clientName,
		"pkceRequired", pkceRequired)

	return nil
}

// seedOAuthProviders creates OAuth provider configurations
func (s *Seeder) seedOAuthProviders(ctx context.Context, tx *sql.Tx) error {
	s.logger.Info("Seeding OAuth providers...")

	for _, provider := range s.config.GetOAuthProviders() {
		if !provider.Enabled {
			s.logger.Info("Skipping disabled OAuth provider", "provider", provider.Provider)
			continue
		}

		if err := s.seedOAuthProvider(ctx, tx, provider); err != nil {
			return fmt.Errorf("seed provider %s: %w", provider.Provider, err)
		}
	}

	s.logger.Info("OAuth providers seeded successfully")
	return nil
}

// seedOAuthProvider creates a single OAuth provider configuration
func (s *Seeder) seedOAuthProvider(ctx context.Context, tx *sql.Tx, provider altalune.OAuthProviderConfig) error {
	s.logger.Info("Checking OAuth provider...", "provider", provider.Provider)

	// Check if provider already exists
	var existingID int64
	err := tx.QueryRowContext(ctx, `
		SELECT id FROM altalune_oauth_providers WHERE provider_type = $1
	`, provider.Provider).Scan(&existingID)

	if err == nil {
		s.logger.Info("OAuth provider already exists, skipping", "provider", provider.Provider)
		return nil
	}

	if err != sql.ErrNoRows {
		return fmt.Errorf("check provider existence: %w", err)
	}

	// Encrypt the provider secret
	s.logger.Info("Creating OAuth provider...", "provider", provider.Provider)
	encryptedSecret, err := EncryptProviderSecret(provider.ClientSecret, s.config.GetIAMEncryptionKey())
	if err != nil {
		return fmt.Errorf("encrypt provider secret: %w", err)
	}

	// Generate public_id at runtime using nanoid
	publicID, err := nanoid.GeneratePublicID()
	if err != nil {
		return fmt.Errorf("generate public_id: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO altalune_oauth_providers (
			public_id, provider_type, client_id, client_secret,
			redirect_url, scopes, enabled, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	`, publicID, provider.Provider, provider.ClientID, encryptedSecret,
		provider.RedirectURL, provider.Scopes, provider.Enabled)

	if err != nil {
		return fmt.Errorf("create provider: %w", err)
	}

	s.logger.Info("OAuth provider created successfully", "provider", provider.Provider)
	return nil
}
