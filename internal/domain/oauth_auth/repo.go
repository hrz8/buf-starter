package oauth_auth

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/hrz8/altalune/internal/postgres"
	"github.com/lib/pq"
)

type repo struct {
	db postgres.DB
}

// NewRepo creates a new OAuth auth repository.
func NewRepo(db postgres.DB) Repositor {
	return &repo{db: db}
}

// CreateAuthorizationCode stores a new authorization code in the database.
func (r *repo) CreateAuthorizationCode(ctx context.Context, input *CreateAuthCodeInput) (*AuthorizationCode, error) {
	code := uuid.New()

	query := `
		INSERT INTO altalune_oauth_authorization_codes (
			code, client_id, user_id, redirect_uri, scope,
			nonce, code_challenge, code_challenge_method, expires_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at
	`

	var id int64
	var createdAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query,
		code,
		input.ClientID,
		input.UserID,
		input.RedirectURI,
		input.Scope,
		input.Nonce,
		input.CodeChallenge,
		input.CodeChallengeMethod,
		input.ExpiresAt,
	).Scan(&id, &createdAt)

	if err != nil {
		return nil, fmt.Errorf("insert authorization code: %w", err)
	}

	return &AuthorizationCode{
		ID:                  id,
		Code:                code,
		ClientID:            input.ClientID,
		UserID:              input.UserID,
		RedirectURI:         input.RedirectURI,
		Scope:               input.Scope,
		Nonce:               input.Nonce,
		CodeChallenge:       input.CodeChallenge,
		CodeChallengeMethod: input.CodeChallengeMethod,
		ExpiresAt:           input.ExpiresAt,
		CreatedAt:           createdAt.Time,
	}, nil
}

// GetAuthorizationCodeByCode retrieves a valid, unexpired, unused authorization code.
func (r *repo) GetAuthorizationCodeByCode(ctx context.Context, code uuid.UUID) (*AuthorizationCode, error) {
	query := `
		SELECT id, code, client_id, user_id, redirect_uri, scope,
		       nonce, code_challenge, code_challenge_method,
		       expires_at, exchange_at, created_at
		FROM altalune_oauth_authorization_codes
		WHERE code = $1
		  AND exchange_at IS NULL
		  AND expires_at > NOW()
	`

	var ac AuthorizationCode
	var nonce, codeChallenge, codeChallengeMethod sql.NullString
	var exchangeAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&ac.ID,
		&ac.Code,
		&ac.ClientID,
		&ac.UserID,
		&ac.RedirectURI,
		&ac.Scope,
		&nonce,
		&codeChallenge,
		&codeChallengeMethod,
		&ac.ExpiresAt,
		&exchangeAt,
		&ac.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrAuthorizationCodeNotFound
		}
		return nil, fmt.Errorf("get authorization code: %w", err)
	}

	if nonce.Valid {
		ac.Nonce = &nonce.String
	}
	if codeChallenge.Valid {
		ac.CodeChallenge = &codeChallenge.String
	}
	if codeChallengeMethod.Valid {
		ac.CodeChallengeMethod = &codeChallengeMethod.String
	}
	if exchangeAt.Valid {
		ac.ExchangeAt = &exchangeAt.Time
	}

	return &ac, nil
}

// MarkCodeExchanged marks an authorization code as used by setting exchange_at.
func (r *repo) MarkCodeExchanged(ctx context.Context, code uuid.UUID) error {
	query := `
		UPDATE altalune_oauth_authorization_codes
		SET exchange_at = NOW(), updated_at = NOW()
		WHERE code = $1 AND exchange_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, code)
	if err != nil {
		return fmt.Errorf("mark code exchanged: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrAuthorizationCodeNotFound
	}

	return nil
}

// CreateRefreshToken stores a new refresh token in the database.
func (r *repo) CreateRefreshToken(ctx context.Context, input *CreateRefreshTokenInput) (*RefreshToken, error) {
	token := uuid.New()

	query := `
		INSERT INTO altalune_oauth_refresh_tokens (
			token, client_id, user_id, scope, nonce, expires_at
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	var id int64
	var createdAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query,
		token,
		input.ClientID,
		input.UserID,
		input.Scope,
		input.Nonce,
		input.ExpiresAt,
	).Scan(&id, &createdAt)

	if err != nil {
		return nil, fmt.Errorf("insert refresh token: %w", err)
	}

	return &RefreshToken{
		ID:        id,
		Token:     token,
		ClientID:  input.ClientID,
		UserID:    input.UserID,
		Scope:     input.Scope,
		Nonce:     input.Nonce,
		ExpiresAt: input.ExpiresAt,
		CreatedAt: createdAt.Time,
	}, nil
}

// GetRefreshTokenByToken retrieves a valid, unexpired, unused refresh token.
func (r *repo) GetRefreshTokenByToken(ctx context.Context, token uuid.UUID) (*RefreshToken, error) {
	query := `
		SELECT id, token, client_id, user_id, scope, nonce,
		       expires_at, exchange_at, created_at
		FROM altalune_oauth_refresh_tokens
		WHERE token = $1
		  AND exchange_at IS NULL
		  AND expires_at > NOW()
	`

	var rt RefreshToken
	var nonce sql.NullString
	var exchangeAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&rt.ID,
		&rt.Token,
		&rt.ClientID,
		&rt.UserID,
		&rt.Scope,
		&nonce,
		&rt.ExpiresAt,
		&exchangeAt,
		&rt.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRefreshTokenNotFound
		}
		return nil, fmt.Errorf("get refresh token: %w", err)
	}

	if nonce.Valid {
		rt.Nonce = &nonce.String
	}
	if exchangeAt.Valid {
		rt.ExchangeAt = &exchangeAt.Time
	}

	return &rt, nil
}

// MarkRefreshTokenExchanged marks a refresh token as used by setting exchange_at.
func (r *repo) MarkRefreshTokenExchanged(ctx context.Context, token uuid.UUID) error {
	query := `
		UPDATE altalune_oauth_refresh_tokens
		SET exchange_at = NOW(), updated_at = NOW()
		WHERE token = $1 AND exchange_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("mark refresh token exchanged: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrRefreshTokenNotFound
	}

	return nil
}

// GetUserConsent retrieves a user's consent record for a specific client.
func (r *repo) GetUserConsent(ctx context.Context, userID int64, clientID uuid.UUID) (*UserConsent, error) {
	query := `
		SELECT id, user_id, client_id, scope, granted_at, revoked_at, created_at
		FROM altalune_oauth_user_consents
		WHERE user_id = $1 AND client_id = $2
	`

	var uc UserConsent
	var revokedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, userID, clientID).Scan(
		&uc.ID,
		&uc.UserID,
		&uc.ClientID,
		&uc.Scope,
		&uc.GrantedAt,
		&revokedAt,
		&uc.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserConsentNotFound
		}
		return nil, fmt.Errorf("get user consent: %w", err)
	}

	if revokedAt.Valid {
		uc.RevokedAt = &revokedAt.Time
	}

	return &uc, nil
}

// CreateOrUpdateUserConsent creates or updates a user's consent for a client.
func (r *repo) CreateOrUpdateUserConsent(ctx context.Context, input *UserConsentInput) (*UserConsent, error) {
	query := `
		INSERT INTO altalune_oauth_user_consents (user_id, client_id, scope, granted_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id, client_id)
		DO UPDATE SET scope = EXCLUDED.scope, granted_at = NOW(), revoked_at = NULL, updated_at = NOW()
		RETURNING id, user_id, client_id, scope, granted_at, revoked_at, created_at
	`

	var uc UserConsent
	var revokedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query,
		input.UserID,
		input.ClientID,
		input.Scope,
	).Scan(
		&uc.ID,
		&uc.UserID,
		&uc.ClientID,
		&uc.Scope,
		&uc.GrantedAt,
		&revokedAt,
		&uc.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("create or update user consent: %w", err)
	}

	if revokedAt.Valid {
		uc.RevokedAt = &revokedAt.Time
	}

	return &uc, nil
}

// GetUserConsents retrieves all consents for a user with client details.
func (r *repo) GetUserConsents(ctx context.Context, userID int64) ([]*UserConsentWithClient, error) {
	query := `
		SELECT
			uc.id,
			uc.user_id,
			uc.client_id,
			c.name as client_name,
			uc.scope,
			uc.granted_at,
			uc.revoked_at,
			uc.created_at
		FROM altalune_oauth_user_consents uc
		INNER JOIN altalune_oauth_clients c ON uc.client_id = c.client_id
		WHERE uc.user_id = $1 AND uc.revoked_at IS NULL
		ORDER BY uc.granted_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get user consents: %w", err)
	}
	defer rows.Close()

	var consents []*UserConsentWithClient
	for rows.Next() {
		var uc UserConsentWithClient
		var revokedAt sql.NullTime

		if err := rows.Scan(
			&uc.ID,
			&uc.UserID,
			&uc.ClientID,
			&uc.ClientName,
			&uc.Scope,
			&uc.GrantedAt,
			&revokedAt,
			&uc.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan user consent: %w", err)
		}

		if revokedAt.Valid {
			uc.RevokedAt = &revokedAt.Time
		}

		consents = append(consents, &uc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate user consents: %w", err)
	}

	return consents, nil
}

// RevokeUserConsent revokes a user's consent for a specific client.
func (r *repo) RevokeUserConsent(ctx context.Context, userID int64, clientID uuid.UUID) error {
	query := `
		UPDATE altalune_oauth_user_consents
		SET revoked_at = NOW(), updated_at = NOW()
		WHERE user_id = $1 AND client_id = $2 AND revoked_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, userID, clientID)
	if err != nil {
		return fmt.Errorf("revoke user consent: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}

	if rows == 0 {
		return ErrUserConsentNotFound
	}

	return nil
}

// GetOAuthClientByClientID retrieves an OAuth client by its client_id.
func (r *repo) GetOAuthClientByClientID(ctx context.Context, clientID uuid.UUID) (*OAuthClientInfo, error) {
	query := `
		SELECT id, client_id, name, client_secret_hash,
		       redirect_uris, pkce_required, is_default, confidential
		FROM altalune_oauth_clients
		WHERE client_id = $1
	`

	var oc OAuthClientInfo
	var redirectURIs pq.StringArray
	var secretHash sql.NullString

	err := r.db.QueryRowContext(ctx, query, clientID).Scan(
		&oc.ID,
		&oc.ClientID,
		&oc.Name,
		&secretHash,
		&redirectURIs,
		&oc.PKCERequired,
		&oc.IsDefault,
		&oc.Confidential,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrOAuthClientNotFound
		}
		return nil, fmt.Errorf("get oauth client: %w", err)
	}

	if secretHash.Valid {
		oc.SecretHash = &secretHash.String
	}
	oc.RedirectURIs = []string(redirectURIs)

	return &oc, nil
}
