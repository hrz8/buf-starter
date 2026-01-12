package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/hrz8/altalune/internal/postgres"
	"github.com/hrz8/altalune/internal/shared/nanoid"
)

// GetUserIdentityByProvider retrieves a user identity by provider and provider user ID
func (r *Repo) GetUserIdentityByProvider(ctx context.Context, provider, providerUserID string) (*UserIdentity, error) {
	query := `
		SELECT
			id, public_id, user_id, provider, provider_user_id,
			email, first_name, last_name, oauth_client_id, origin_oauth_client_name, last_login_at,
			created_at, updated_at
		FROM altalune_user_identities
		WHERE provider = $1 AND provider_user_id = $2
		LIMIT 1
	`

	var identity UserIdentity
	var email, firstName, lastName, oauthClientID, originOAuthClientName sql.NullString
	var lastLoginAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, provider, providerUserID).Scan(
		&identity.ID,
		&identity.PublicID,
		&identity.UserID,
		&identity.Provider,
		&identity.ProviderUserID,
		&email,
		&firstName,
		&lastName,
		&oauthClientID,
		&originOAuthClientName,
		&lastLoginAt,
		&identity.CreatedAt,
		&identity.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user identity: %w", err)
	}

	if email.Valid {
		identity.Email = email.String
	}
	if firstName.Valid {
		identity.FirstName = firstName.String
	}
	if lastName.Valid {
		identity.LastName = lastName.String
	}
	if oauthClientID.Valid {
		identity.OAuthClientID = &oauthClientID.String
	}
	if originOAuthClientName.Valid {
		identity.OriginOAuthClientName = &originOAuthClientName.String
	}
	if lastLoginAt.Valid {
		identity.LastLoginAt = &lastLoginAt.Time
	}

	return &identity, nil
}

// GetUserIdentities retrieves all identities for a user
func (r *Repo) GetUserIdentities(ctx context.Context, userID int64) ([]*UserIdentity, error) {
	query := `
		SELECT
			id, public_id, user_id, provider, provider_user_id,
			email, first_name, last_name, oauth_client_id, origin_oauth_client_name, last_login_at,
			created_at, updated_at
		FROM altalune_user_identities
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query user identities: %w", err)
	}
	defer rows.Close()

	var identities []*UserIdentity
	for rows.Next() {
		var identity UserIdentity
		var email, firstName, lastName, oauthClientID, originOAuthClientName sql.NullString
		var lastLoginAt sql.NullTime

		err := rows.Scan(
			&identity.ID,
			&identity.PublicID,
			&identity.UserID,
			&identity.Provider,
			&identity.ProviderUserID,
			&email,
			&firstName,
			&lastName,
			&oauthClientID,
			&originOAuthClientName,
			&lastLoginAt,
			&identity.CreatedAt,
			&identity.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("scan user identity: %w", err)
		}

		if email.Valid {
			identity.Email = email.String
		}
		if firstName.Valid {
			identity.FirstName = firstName.String
		}
		if lastName.Valid {
			identity.LastName = lastName.String
		}
		if oauthClientID.Valid {
			identity.OAuthClientID = &oauthClientID.String
		}
		if originOAuthClientName.Valid {
			identity.OriginOAuthClientName = &originOAuthClientName.String
		}
		if lastLoginAt.Valid {
			identity.LastLoginAt = &lastLoginAt.Time
		}

		identities = append(identities, &identity)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate user identities: %w", err)
	}

	return identities, nil
}

// CreateUserIdentity creates a new user identity record
func (r *Repo) CreateUserIdentity(ctx context.Context, input *CreateUserIdentityInput) error {
	publicID, _ := nanoid.GeneratePublicID()

	query := `
		INSERT INTO altalune_user_identities (
			public_id, user_id, provider, provider_user_id,
			email, first_name, last_name, oauth_client_id, origin_oauth_client_name, last_login_at,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW(), NOW())
	`

	_, err := r.db.ExecContext(ctx, query,
		publicID,
		input.UserID,
		input.Provider,
		input.ProviderUserID,
		input.Email,
		input.FirstName,
		input.LastName,
		input.OAuthClientID,
		input.OriginOAuthClientName,
	)

	if err != nil {
		return fmt.Errorf("create user identity: %w", err)
	}

	return nil
}

// UpdateUserIdentityLastLogin updates the last_login_at timestamp for a user identity
func (r *Repo) UpdateUserIdentityLastLogin(ctx context.Context, userID int64, provider string) error {
	query := `
		UPDATE altalune_user_identities
		SET last_login_at = NOW(), updated_at = NOW()
		WHERE user_id = $1 AND provider = $2
	`

	result, err := r.db.ExecContext(ctx, query, userID, provider)
	if err != nil {
		return fmt.Errorf("update user identity last login: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// AddProjectMember adds a user to a project with the specified role
func (r *Repo) AddProjectMember(ctx context.Context, projectID, userID int64, role string) error {
	publicID, _ := nanoid.GeneratePublicID()

	query := `
		INSERT INTO altalune_project_members (
			public_id, project_id, user_id, role, created_at, updated_at
		) VALUES ($1, $2, $3, $4, NOW(), NOW())
	`

	_, err := r.db.ExecContext(ctx, query, publicID, projectID, userID, role)
	if err != nil {
		if postgres.IsUniqueViolation(err) {
			return nil
		}
		return fmt.Errorf("add project member: %w", err)
	}

	return nil
}
