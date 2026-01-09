package oauth_client

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hrz8/altalune/internal/postgres"
	"github.com/hrz8/altalune/internal/shared/nanoid"
	"github.com/hrz8/altalune/internal/shared/password"
	"github.com/hrz8/altalune/internal/shared/query"
	"github.com/lib/pq"
)

type repo struct {
	db postgres.DB
}

// NewRepo creates a new OAuth client repository
func NewRepo(db postgres.DB) Repositor {
	return &repo{db: db}
}

// Create creates a new OAuth client with generated client_id and hashed secret
func (r *repo) Create(ctx context.Context, input *CreateOAuthClientInput) (*CreateOAuthClientResult, error) {
	// 1. Generate public ID (nanoid)
	publicID, err := nanoid.GeneratePublicID()
	if err != nil {
		return nil, fmt.Errorf("generate public id: %w", err)
	}

	// 2. Generate UUID client_id for OAuth flow
	clientID := uuid.New()

	// 3. Generate secure random client secret (minimum 32 characters)
	clientSecret := generateSecureRandom(32)

	// 4. Hash secret with Argon2id (using production parameters from T19/oauth_seeder)
	hashedSecret, err := password.HashPassword(clientSecret, password.HashOption{
		Iterations: 2,         // Time cost
		Memory:     64 * 1024, // 64MB
		Threads:    4,         // Parallelism
		Len:        32,        // Hash length
	})
	if err != nil {
		return nil, fmt.Errorf("hash client secret: %w", err)
	}

	// 5. Insert into partitioned table
	insertQuery := `
		INSERT INTO altalune_oauth_clients (
			project_id, public_id, name, client_id,
			client_secret_hash, redirect_uris, pkce_required, is_default
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	var id int64
	var createdAt, updatedAt sql.NullTime

	err = r.db.QueryRowContext(ctx, insertQuery,
		input.ProjectID,
		publicID,
		input.Name,
		clientID,
		hashedSecret,
		pq.Array(input.RedirectURIs),
		input.PKCERequired,
		false, // is_default (always false for user-created clients)
	).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		if postgres.IsUniqueViolation(err) {
			return nil, ErrOAuthClientAlreadyExists
		}
		return nil, fmt.Errorf("insert oauth client: %w", err)
	}

	// 6. Build domain model
	client := &OAuthClient{
		ID:           publicID,
		ProjectID:    input.ProjectPublicID,
		Name:         input.Name,
		ClientID:     clientID,
		RedirectURIs: input.RedirectURIs,
		PKCERequired: input.PKCERequired,
		IsDefault:    false,
		CreatedAt:    createdAt.Time,
		UpdatedAt:    updatedAt.Time,
	}

	// 7. Return client with PLAINTEXT secret (ONLY time it's returned)
	return &CreateOAuthClientResult{
		Client:       client,
		ClientSecret: clientSecret,
	}, nil
}

// Query returns a paginated list of OAuth clients for a project
func (r *repo) Query(ctx context.Context, projectID int64, params *query.QueryParams) (*query.QueryResult[OAuthClient], error) {
	// Base query WITHOUT client_secret_hash (security: never expose secret hash)
	baseQuery := `
		SELECT id, public_id, name, client_id,
		       redirect_uris, pkce_required, is_default,
		       created_at, updated_at
		FROM altalune_oauth_clients
		WHERE project_id = $1
	`

	var whereConditions []string
	var args []interface{}
	args = append(args, projectID) // $1
	argCounter := 2

	// Handle keyword search (search by name)
	if params.Keyword != "" {
		searchCondition := fmt.Sprintf("LOWER(name) LIKE $%d", argCounter)
		whereConditions = append(whereConditions, searchCondition)
		searchPattern := "%" + strings.ToLower(params.Keyword) + "%"
		args = append(args, searchPattern)
		argCounter++
	}

	// Handle filters
	if params.Filters != nil {
		for field, values := range params.Filters {
			if len(values) == 0 {
				continue
			}

			var dbColumn string
			switch field {
			case "name", "names":
				dbColumn = "name"
			case "pkce_required":
				// Handle boolean filter
				for _, value := range values {
					if value == "true" || value == "false" {
						whereConditions = append(whereConditions, fmt.Sprintf("pkce_required = $%d", argCounter))
						args = append(args, value == "true")
						argCounter++
					}
				}
				continue
			case "is_default":
				// Handle boolean filter
				for _, value := range values {
					if value == "true" || value == "false" {
						whereConditions = append(whereConditions, fmt.Sprintf("is_default = $%d", argCounter))
						args = append(args, value == "true")
						argCounter++
					}
				}
				continue
			default:
				continue
			}

			// Handle string filters
			placeholders := make([]string, len(values))
			for i, value := range values {
				placeholders[i] = fmt.Sprintf("$%d", argCounter)
				args = append(args, strings.ToLower(value))
				argCounter++
			}
			filterCondition := fmt.Sprintf("LOWER(%s) IN (%s)", dbColumn, strings.Join(placeholders, ","))
			whereConditions = append(whereConditions, filterCondition)
		}
	}

	// Combine WHERE conditions
	if len(whereConditions) > 0 {
		baseQuery += " AND " + strings.Join(whereConditions, " AND ")
	}

	// Get total count BEFORE pagination
	countQuery := "SELECT COUNT(*) FROM (" + baseQuery + ") as filtered"
	var totalRows int32
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalRows)
	if err != nil {
		return nil, fmt.Errorf("count oauth clients: %w", err)
	}

	// Add sorting (default: created_at DESC)
	orderBy := "created_at DESC"
	if params.Sorting != nil && params.Sorting.Field != "" {
		direction := "ASC"
		if params.Sorting.Order == query.SortOrderDesc {
			direction = "DESC"
		}
		// Validate sort column
		switch params.Sorting.Field {
		case "name", "created_at", "updated_at":
			orderBy = fmt.Sprintf("%s %s", params.Sorting.Field, direction)
		}
	}
	baseQuery += fmt.Sprintf(" ORDER BY %s", orderBy)

	// Add pagination
	pageSize := params.Pagination.PageSize
	offset := (params.Pagination.Page - 1) * pageSize

	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCounter, argCounter+1)
	args = append(args, pageSize, offset)

	// Execute query
	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("query oauth clients: %w", err)
	}
	defer rows.Close()

	// Get project public ID for response
	var projectPublicID string
	projectQuery := "SELECT public_id FROM altalune_projects WHERE id = $1"
	err = r.db.QueryRowContext(ctx, projectQuery, projectID).Scan(&projectPublicID)
	if err != nil {
		return nil, fmt.Errorf("get project public id: %w", err)
	}

	// Scan rows
	data := make([]*OAuthClient, 0)
	for rows.Next() {
		var result OAuthClientQueryResult
		var redirectURIs pq.StringArray

		err := rows.Scan(
			&result.ID,
			&result.PublicID,
			&result.Name,
			&result.ClientID,
			&redirectURIs,
			&result.PKCERequired,
			&result.IsDefault,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan oauth client: %w", err)
		}

		result.ProjectID = projectID
		result.ProjectPublicID = projectPublicID
		result.RedirectURIs = []string(redirectURIs)

		data = append(data, result.ToOAuthClient())
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	// Calculate page count
	pageCount := (totalRows + pageSize - 1) / pageSize
	if pageCount == 0 {
		pageCount = 1
	}

	return &query.QueryResult[OAuthClient]{
		Data:       data,
		TotalRows:  totalRows,
		TotalPages: pageCount,
		Filters:    params.Filters,
	}, nil
}

// GetByPublicID retrieves an OAuth client by its public nanoid
func (r *repo) GetByPublicID(ctx context.Context, projectID int64, publicID string) (*OAuthClient, error) {
	// Query WITHOUT client_secret_hash
	selectQuery := `
		SELECT id, public_id, name, client_id,
		       redirect_uris, pkce_required, is_default,
		       created_at, updated_at
		FROM altalune_oauth_clients
		WHERE project_id = $1 AND public_id = $2
	`

	var result OAuthClientQueryResult
	var redirectURIs pq.StringArray

	err := r.db.QueryRowContext(ctx, selectQuery, projectID, publicID).Scan(
		&result.ID,
		&result.PublicID,
		&result.Name,
		&result.ClientID,
		&redirectURIs,
		&result.PKCERequired,
		&result.IsDefault,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrOAuthClientNotFound
		}
		return nil, fmt.Errorf("get oauth client: %w", err)
	}

	// Get project public ID
	var projectPublicID string
	projectQuery := "SELECT public_id FROM altalune_projects WHERE id = $1"
	err = r.db.QueryRowContext(ctx, projectQuery, projectID).Scan(&projectPublicID)
	if err != nil {
		return nil, fmt.Errorf("get project public id: %w", err)
	}

	result.ProjectID = projectID
	result.ProjectPublicID = projectPublicID
	result.RedirectURIs = []string(redirectURIs)

	return result.ToOAuthClient(), nil
}

// GetByClientID retrieves an OAuth client by its UUID client_id (for OAuth flows)
func (r *repo) GetByClientID(ctx context.Context, clientID string) (*OAuthClient, error) {
	// Parse UUID
	clientUUID, err := uuid.Parse(clientID)
	if err != nil {
		return nil, fmt.Errorf("invalid client_id UUID: %w", err)
	}

	// Query WITHOUT client_secret_hash
	selectQuery := `
		SELECT oc.id, oc.public_id, oc.project_id, p.public_id as project_public_id,
		       oc.name, oc.client_id, oc.redirect_uris, oc.pkce_required,
		       oc.is_default, oc.created_at, oc.updated_at
		FROM altalune_oauth_clients oc
		JOIN altalune_projects p ON oc.project_id = p.id
		WHERE oc.client_id = $1
	`

	var result OAuthClientQueryResult
	var redirectURIs pq.StringArray

	err = r.db.QueryRowContext(ctx, selectQuery, clientUUID).Scan(
		&result.ID,
		&result.PublicID,
		&result.ProjectID,
		&result.ProjectPublicID,
		&result.Name,
		&result.ClientID,
		&redirectURIs,
		&result.PKCERequired,
		&result.IsDefault,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrOAuthClientNotFound
		}
		return nil, fmt.Errorf("get oauth client by client_id: %w", err)
	}

	result.RedirectURIs = []string(redirectURIs)

	return result.ToOAuthClient(), nil
}

// Update updates an existing OAuth client
func (r *repo) Update(ctx context.Context, input *UpdateOAuthClientInput) (*OAuthClient, error) {
	// Build dynamic UPDATE query
	setClauses := []string{}
	args := []interface{}{input.ProjectID, input.PublicID} // $1, $2
	argCounter := 3

	if input.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argCounter))
		args = append(args, *input.Name)
		argCounter++
	}

	if len(input.RedirectURIs) > 0 {
		setClauses = append(setClauses, fmt.Sprintf("redirect_uris = $%d", argCounter))
		args = append(args, pq.Array(input.RedirectURIs))
		argCounter++
	}

	if input.PKCERequired != nil {
		setClauses = append(setClauses, fmt.Sprintf("pkce_required = $%d", argCounter))
		args = append(args, *input.PKCERequired)
		argCounter++
	}

	// Always update updated_at
	setClauses = append(setClauses, "updated_at = CURRENT_TIMESTAMP")

	if len(setClauses) == 1 { // Only updated_at
		return nil, fmt.Errorf("no fields to update")
	}

	updateQuery := fmt.Sprintf(`
		UPDATE altalune_oauth_clients
		SET %s
		WHERE project_id = $1 AND public_id = $2
		RETURNING id, public_id, name, client_id,
		          redirect_uris, pkce_required, is_default,
		          created_at, updated_at
	`, strings.Join(setClauses, ", "))

	var result OAuthClientQueryResult
	var redirectURIs pq.StringArray

	err := r.db.QueryRowContext(ctx, updateQuery, args...).Scan(
		&result.ID,
		&result.PublicID,
		&result.Name,
		&result.ClientID,
		&redirectURIs,
		&result.PKCERequired,
		&result.IsDefault,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrOAuthClientNotFound
		}
		if postgres.IsUniqueViolation(err) {
			return nil, ErrOAuthClientAlreadyExists
		}
		return nil, fmt.Errorf("update oauth client: %w", err)
	}

	// Get project public ID
	var projectPublicID string
	projectQuery := "SELECT public_id FROM altalune_projects WHERE id = $1"
	err = r.db.QueryRowContext(ctx, projectQuery, input.ProjectID).Scan(&projectPublicID)
	if err != nil {
		return nil, fmt.Errorf("get project public id: %w", err)
	}

	result.ProjectID = input.ProjectID
	result.ProjectPublicID = projectPublicID
	result.RedirectURIs = []string(redirectURIs)

	return result.ToOAuthClient(), nil
}

// Delete deletes an OAuth client (with default client protection)
func (r *repo) Delete(ctx context.Context, projectID int64, publicID string) error {
	// First check if it's the default client
	var isDefault bool
	checkQuery := "SELECT is_default FROM altalune_oauth_clients WHERE project_id = $1 AND public_id = $2"
	err := r.db.QueryRowContext(ctx, checkQuery, projectID, publicID).Scan(&isDefault)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrOAuthClientNotFound
		}
		return fmt.Errorf("check default client: %w", err)
	}

	// Protect default client from deletion
	if isDefault {
		return ErrDefaultClientCannotBeDeleted
	}

	// Delete client
	deleteQuery := "DELETE FROM altalune_oauth_clients WHERE project_id = $1 AND public_id = $2"
	result, err := r.db.ExecContext(ctx, deleteQuery, projectID, publicID)
	if err != nil {
		return fmt.Errorf("delete oauth client: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrOAuthClientNotFound
	}

	return nil
}

// RevealClientSecret retrieves the hashed client secret (with audit logging)
// NOTE: Returns the Argon2id hash. In production, you might want to:
// 1. Add audit logging here
// 2. Return encrypted secret if you store it encrypted
// 3. Add rate limiting for this operation
func (r *repo) RevealClientSecret(ctx context.Context, projectID int64, publicID string) (string, error) {
	// Query for client_secret_hash (this is the ONLY method that returns it)
	selectQuery := `
		SELECT client_secret_hash
		FROM altalune_oauth_clients
		WHERE project_id = $1 AND public_id = $2
	`

	var hashedSecret string
	err := r.db.QueryRowContext(ctx, selectQuery, projectID, publicID).Scan(&hashedSecret)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrOAuthClientNotFound
		}
		return "", fmt.Errorf("reveal client secret: %w", err)
	}

	// TODO: Add audit logging here
	// logger.Info("oauth_client_secret_revealed",
	//     "project_id", projectID,
	//     "client_public_id", publicID,
	//     "user_id", userIDFromContext,
	// )

	// Return hashed secret (Argon2id PHC string format)
	// Frontend will display this as-is (no way to recover plaintext after creation)
	return hashedSecret, nil
}

// generateSecureRandom generates a cryptographically secure random string
func generateSecureRandom(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(fmt.Sprintf("failed to generate random bytes: %v", err))
	}
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}
