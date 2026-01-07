package oauth_provider

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hrz8/altalune/internal/postgres"
	"github.com/hrz8/altalune/internal/shared/crypto"
	"github.com/hrz8/altalune/internal/shared/nanoid"
	"github.com/hrz8/altalune/internal/shared/query"
)

type Repo struct {
	db            postgres.DB
	encryptionKey []byte
}

func NewRepo(db postgres.DB, encryptionKey []byte) *Repo {
	return &Repo{
		db:            db,
		encryptionKey: encryptionKey,
	}
}

func (r *Repo) Query(ctx context.Context, params *query.QueryParams) (*query.QueryResult[OAuthProvider], error) {
	// Build the base query - NEVER select client_secret
	baseQuery := `
		SELECT
			id,
			public_id,
			provider_type,
			client_id,
			redirect_url,
			scopes,
			enabled,
			created_at,
			updated_at
		FROM altalune_oauth_providers
		WHERE 1=1
	`

	// Build WHERE conditions for filters and search
	var whereConditions []string
	var args []interface{}
	argCounter := 1

	// Handle keyword search (search in provider_type, client_id, redirect_url)
	if params.Keyword != "" {
		searchCondition := fmt.Sprintf(`
			(LOWER(provider_type) LIKE $%d OR
			 LOWER(client_id) LIKE $%d OR
			 LOWER(redirect_url) LIKE $%d)
		`, argCounter, argCounter, argCounter)
		whereConditions = append(whereConditions, searchCondition)
		searchPattern := "%" + strings.ToLower(params.Keyword) + "%"
		args = append(args, searchPattern)
		argCounter++
	}

	// Handle column-specific filters
	if params.Filters != nil {
		for field, values := range params.Filters {
			if len(values) == 0 {
				continue
			}

			// Map field names to database columns
			var dbColumn string
			switch field {
			case "provider_type", "providerType":
				dbColumn = "provider_type"
			case "enabled":
				dbColumn = "enabled"
			default:
				continue // Skip unknown fields
			}

			// Build IN clause for multiple values
			placeholders := make([]string, len(values))
			for i, value := range values {
				placeholders[i] = fmt.Sprintf("$%d", argCounter)

				// For boolean fields, convert string to boolean
				if dbColumn == "enabled" {
					boolValue := value == "true" || value == "1"
					args = append(args, boolValue)
				} else {
					args = append(args, strings.ToLower(value))
				}
				argCounter++
			}

			if dbColumn == "enabled" {
				filterCondition := fmt.Sprintf("%s IN (%s)", dbColumn, strings.Join(placeholders, ","))
				whereConditions = append(whereConditions, filterCondition)
			} else {
				filterCondition := fmt.Sprintf("LOWER(%s) IN (%s)", dbColumn, strings.Join(placeholders, ","))
				whereConditions = append(whereConditions, filterCondition)
			}
		}
	}

	// Combine all WHERE conditions
	if len(whereConditions) > 0 {
		baseQuery += " AND " + strings.Join(whereConditions, " AND ")
	}

	// First, get the total count before pagination
	countQuery := "SELECT COUNT(*) FROM (" + baseQuery + ") as filtered"
	var totalRows int32
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalRows)
	if err != nil {
		return nil, fmt.Errorf("count oauth providers: %w", err)
	}

	// Add ORDER BY clause
	orderClause := r.buildOrderClause(params.Sorting)
	baseQuery += orderClause

	// Add pagination
	pageSize := params.Pagination.PageSize
	page := params.Pagination.Page
	offset := (page - 1) * pageSize
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCounter, argCounter+1)
	args = append(args, pageSize, offset)

	// Execute the main query
	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("query oauth providers: %w", err)
	}
	defer rows.Close()

	// Collect queryResults
	queryResults := make([]*OAuthProviderQueryResult, 0)
	for rows.Next() {
		var provider OAuthProviderQueryResult

		err := rows.Scan(
			&provider.ID,
			&provider.PublicID,
			&provider.ProviderType,
			&provider.ClientID,
			&provider.RedirectURL,
			&provider.Scopes,
			&provider.Enabled,
			&provider.CreatedAt,
			&provider.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan oauth provider row: %w", err)
		}

		queryResults = append(queryResults, &provider)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate oauth provider rows: %w", err)
	}

	// Calculate total pages
	var totalPages int32
	if totalRows > 0 {
		totalPages = (totalRows + pageSize - 1) / pageSize
	}

	// Get distinct values for filters
	filters, err := r.getDistinctValues(ctx)
	if err != nil {
		// Don't fail the entire query if we can't get filters
		filters = make(map[string][]string)
	}

	// Convert to domain models
	results := make([]*OAuthProvider, 0)
	for _, v := range queryResults {
		results = append(results, v.ToOAuthProvider())
	}

	return &query.QueryResult[OAuthProvider]{
		Data:       results,
		TotalRows:  totalRows,
		TotalPages: totalPages,
		Filters:    filters,
	}, nil
}

func (r *Repo) buildOrderClause(sorting *query.SortingParams) string {
	if sorting == nil || sorting.Field == "" {
		return " ORDER BY created_at DESC" // Default sorting
	}

	// Map field to database column
	var dbColumn string
	switch sorting.Field {
	case "providerType", "provider_type":
		dbColumn = "provider_type"
	case "clientId", "client_id":
		dbColumn = "client_id"
	case "redirectUrl", "redirect_url":
		dbColumn = "redirect_url"
	case "enabled":
		dbColumn = "enabled"
	case "createdAt", "created_at":
		dbColumn = "created_at"
	case "updatedAt", "updated_at":
		dbColumn = "updated_at"
	case "id":
		dbColumn = "id"
	default:
		dbColumn = "created_at" // Fallback to default
	}

	// Determine sort direction
	direction := "ASC"
	if sorting.Order == query.SortOrderDesc {
		direction = "DESC"
	}

	return fmt.Sprintf(" ORDER BY %s %s", dbColumn, direction)
}

func (r *Repo) getDistinctValues(ctx context.Context) (map[string][]string, error) {
	filters := make(map[string][]string)

	// provider_type filter
	providerTypeQuery := `
		SELECT DISTINCT provider_type
		FROM altalune_oauth_providers
		ORDER BY provider_type
	`
	rows, err := r.db.QueryContext(ctx, providerTypeQuery)
	if err != nil {
		return nil, fmt.Errorf("get distinct provider types: %w", err)
	}
	defer rows.Close()

	providerTypes := make([]string, 0)
	for rows.Next() {
		var providerType string
		if err := rows.Scan(&providerType); err != nil {
			return nil, fmt.Errorf("scan provider type: %w", err)
		}
		providerTypes = append(providerTypes, providerType)
	}
	filters["provider_type"] = providerTypes

	// enabled filter (boolean)
	filters["enabled"] = []string{"true", "false"}

	return filters, nil
}

// Create creates a new OAuth provider with encrypted client_secret
func (r *Repo) Create(ctx context.Context, input *CreateOAuthProviderInput) (*CreateOAuthProviderResult, error) {
	// Check for duplicate provider_type
	existing, err := r.GetByProviderType(ctx, input.ProviderType)
	if err == nil && existing != nil {
		return nil, ErrDuplicateProviderType
	}
	if err != nil && !errors.Is(err, ErrOAuthProviderNotFound) {
		return nil, fmt.Errorf("check duplicate provider type: %w", err)
	}

	// Generate public ID
	publicID, _ := nanoid.GeneratePublicID()

	// CRITICAL: Encrypt client_secret before storing
	encryptedSecret, err := crypto.Encrypt(input.ClientSecret, r.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrEncryptionFailed, err)
	}

	// Insert query
	insertQuery := `
		INSERT INTO altalune_oauth_providers (
			public_id,
			provider_type,
			client_id,
			client_secret,
			redirect_url,
			scopes,
			enabled,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, public_id, provider_type, client_id, redirect_url, scopes, enabled, created_at, updated_at
	`

	now := time.Now()
	var result CreateOAuthProviderResult

	err = r.db.QueryRowContext(
		ctx,
		insertQuery,
		publicID,
		string(input.ProviderType),
		input.ClientID,
		encryptedSecret,
		input.RedirectURL,
		input.Scopes,
		input.Enabled,
		now,
		now,
	).Scan(
		&result.ID,
		&result.PublicID,
		&result.ProviderType,
		&result.ClientID,
		&result.RedirectURL,
		&result.Scopes,
		&result.Enabled,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		// Check for unique constraint violation
		if postgres.IsUniqueViolation(err) {
			// Check if it's the provider_type constraint
			if strings.Contains(err.Error(), "ux_altalune_oauth_providers_provider_type") {
				return nil, ErrDuplicateProviderType
			}
		}
		return nil, fmt.Errorf("create oauth provider: %w", err)
	}

	return &result, nil
}

// GetByID retrieves an OAuth provider by public ID
// CRITICAL: NEVER selects client_secret column
func (r *Repo) GetByID(ctx context.Context, publicID string) (*OAuthProvider, error) {
	sqlQuery := `
		SELECT
			public_id,
			provider_type,
			client_id,
			redirect_url,
			scopes,
			enabled,
			created_at,
			updated_at
		FROM altalune_oauth_providers
		WHERE public_id = $1
	`

	var provider OAuthProvider
	var providerTypeStr string

	err := r.db.QueryRowContext(ctx, sqlQuery, publicID).Scan(
		&provider.ID,
		&providerTypeStr,
		&provider.ClientID,
		&provider.RedirectURL,
		&provider.Scopes,
		&provider.Enabled,
		&provider.CreatedAt,
		&provider.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrOAuthProviderNotFound
		}
		return nil, fmt.Errorf("failed to get oauth provider: %w", err)
	}

	provider.ProviderType = ProviderType(providerTypeStr)
	provider.ClientSecretSet = true // If record exists, secret is set

	return &provider, nil
}

// GetByProviderType retrieves an OAuth provider by provider type
// CRITICAL: NEVER selects client_secret column
func (r *Repo) GetByProviderType(ctx context.Context, providerType ProviderType) (*OAuthProvider, error) {
	sqlQuery := `
		SELECT
			public_id,
			provider_type,
			client_id,
			redirect_url,
			scopes,
			enabled,
			created_at,
			updated_at
		FROM altalune_oauth_providers
		WHERE provider_type = $1
		LIMIT 1
	`

	var provider OAuthProvider
	var providerTypeStr string

	err := r.db.QueryRowContext(ctx, sqlQuery, string(providerType)).Scan(
		&provider.ID,
		&providerTypeStr,
		&provider.ClientID,
		&provider.RedirectURL,
		&provider.Scopes,
		&provider.Enabled,
		&provider.CreatedAt,
		&provider.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrOAuthProviderNotFound
		}
		return nil, fmt.Errorf("get oauth provider by type: %w", err)
	}

	provider.ProviderType = ProviderType(providerTypeStr)
	provider.ClientSecretSet = true // If record exists, secret is set

	return &provider, nil
}

// Update updates an OAuth provider
// CRITICAL: Re-encrypts client_secret if provided (non-empty)
func (r *Repo) Update(ctx context.Context, input *UpdateOAuthProviderInput) (*UpdateOAuthProviderResult, error) {
	// Build dynamic UPDATE query based on whether client_secret is provided
	var sqlQuery string
	var args []interface{}

	if input.ClientSecret != "" {
		// CRITICAL: Re-encrypt new client_secret
		encryptedSecret, err := crypto.Encrypt(input.ClientSecret, r.encryptionKey)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrEncryptionFailed, err)
		}

		sqlQuery = `
			UPDATE altalune_oauth_providers
			SET client_id = $1, client_secret = $2, redirect_url = $3, scopes = $4, enabled = $5, updated_at = CURRENT_TIMESTAMP
			WHERE public_id = $6
			RETURNING id, public_id, client_id, redirect_url, scopes, enabled, updated_at
		`
		args = []interface{}{
			input.ClientID,
			encryptedSecret,
			input.RedirectURL,
			input.Scopes,
			input.Enabled,
			input.PublicID,
		}
	} else {
		// Keep existing client_secret (don't update it)
		sqlQuery = `
			UPDATE altalune_oauth_providers
			SET client_id = $1, redirect_url = $2, scopes = $3, enabled = $4, updated_at = CURRENT_TIMESTAMP
			WHERE public_id = $5
			RETURNING id, public_id, client_id, redirect_url, scopes, enabled, updated_at
		`
		args = []interface{}{
			input.ClientID,
			input.RedirectURL,
			input.Scopes,
			input.Enabled,
			input.PublicID,
		}
	}

	var result UpdateOAuthProviderResult

	err := r.db.QueryRowContext(ctx, sqlQuery, args...).Scan(
		&result.ID,
		&result.PublicID,
		&result.ClientID,
		&result.RedirectURL,
		&result.Scopes,
		&result.Enabled,
		&result.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrOAuthProviderNotFound
		}
		return nil, fmt.Errorf("failed to update oauth provider: %w", err)
	}

	return &result, nil
}

// Delete deletes an OAuth provider from the database
func (r *Repo) Delete(ctx context.Context, input *DeleteOAuthProviderInput) error {
	sqlQuery := `DELETE FROM altalune_oauth_providers WHERE public_id = $1`

	result, err := r.db.ExecContext(ctx, sqlQuery, input.PublicID)
	if err != nil {
		return fmt.Errorf("failed to delete oauth provider: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrOAuthProviderNotFound
	}

	return nil
}

// RevealClientSecret decrypts and returns the plaintext client secret
// CRITICAL: Only method that SELECTs client_secret column
func (r *Repo) RevealClientSecret(ctx context.Context, publicID string) (string, error) {
	sqlQuery := `
		SELECT client_secret
		FROM altalune_oauth_providers
		WHERE public_id = $1
	`

	var encryptedSecret string
	err := r.db.QueryRowContext(ctx, sqlQuery, publicID).Scan(&encryptedSecret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrOAuthProviderNotFound
		}
		return "", fmt.Errorf("get encrypted client secret: %w", err)
	}

	// CRITICAL: Decrypt the client_secret
	plaintext, err := crypto.Decrypt(encryptedSecret, r.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrDecryptionFailed, err)
	}

	return plaintext, nil
}
