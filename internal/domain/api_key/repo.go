package api_key

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hrz8/altalune/internal/postgres"
	"github.com/hrz8/altalune/internal/shared/nanoid"
	"github.com/hrz8/altalune/internal/shared/query"
)

type Repo struct {
	db postgres.DB
}

func NewRepo(db postgres.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Query(ctx context.Context, projectID int64, params *query.QueryParams) (*query.QueryResult[ApiKey], error) {
	// Build the base query
	baseQuery := `
		SELECT
			id,
			public_id,
			name,
			expiration,
			active,
			created_at,
			updated_at
		FROM altalune_project_api_keys
		WHERE project_id = $1
	`

	// Build WHERE conditions for filters and search
	var whereConditions []string
	var args []interface{}
	args = append(args, projectID) // $1
	argCounter := 2

	// Handle keyword search (global search across multiple fields)
	if params.Keyword != "" {
		searchCondition := fmt.Sprintf(`
			(LOWER(name) LIKE $%d)
		`, argCounter)
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
			case "name", "names":
				dbColumn = "name"
			case "status", "statuses":
				// Handle combined status as a special case
				r.handleCombinedStatusFilter(&whereConditions, &args, &argCounter, values)
				continue
			default:
				continue // Skip unknown fields
			}

			// Build IN clause for multiple values
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

	// Combine all WHERE conditions
	if len(whereConditions) > 0 {
		baseQuery += " AND " + strings.Join(whereConditions, " AND ")
	}

	// First, get the total count before pagination
	countQuery := "SELECT COUNT(*) FROM (" + baseQuery + ") as filtered"
	var totalRows int32
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalRows)
	if err != nil {
		return nil, fmt.Errorf("count api keys: %w", err)
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

	// Execute query
	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("query api keys: %w", err)
	}
	defer rows.Close()

	// Scan results
	var results []*ApiKeyQueryResult
	for rows.Next() {
		var result ApiKeyQueryResult
		err := rows.Scan(
			&result.ID,
			&result.PublicID,
			&result.Name,
			&result.Expiration,
			&result.Active,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan api key: %w", err)
		}
		results = append(results, &result)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate api key rows: %w", err)
	}

	// Convert to domain objects
	apiKeys := make([]*ApiKey, 0, len(results))
	for _, result := range results {
		apiKeys = append(apiKeys, result.ToApiKey())
	}

	// Calculate pagination
	pageCount := int32((totalRows + pageSize - 1) / pageSize)

	// Get filters (for dropdown values)
	filters, err := r.getDistinctValues(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("get distinct values: %w", err)
	}

	return &query.QueryResult[ApiKey]{
		Data:       apiKeys,
		TotalRows:  totalRows,
		TotalPages: pageCount,
		Filters:    filters,
	}, nil
}

func (r *Repo) handleExpirationFilter(whereConditions *[]string, args *[]interface{}, argCounter *int, values []string) {
	now := time.Now()
	var expirationConditions []string

	for _, value := range values {
		switch strings.ToLower(value) {
		case "active":
			expirationConditions = append(expirationConditions, fmt.Sprintf("expiration > $%d", *argCounter))
			*args = append(*args, now)
			*argCounter++
		case "expired":
			expirationConditions = append(expirationConditions, fmt.Sprintf("expiration <= $%d", *argCounter))
			*args = append(*args, now)
			*argCounter++
		case "expiring_soon":
			soonThreshold := now.AddDate(0, 0, 30) // 30 days from now
			expirationConditions = append(expirationConditions, fmt.Sprintf("expiration > $%d AND expiration <= $%d", *argCounter, *argCounter+1))
			*args = append(*args, now, soonThreshold)
			*argCounter += 2
		}
	}

	if len(expirationConditions) > 0 {
		*whereConditions = append(*whereConditions, "("+strings.Join(expirationConditions, " OR ")+")")
	}
}

func (r *Repo) handleCombinedStatusFilter(whereConditions *[]string, args *[]interface{}, argCounter *int, values []string) {
	now := time.Now()
	tenDaysFromNow := now.AddDate(0, 0, 10) // 10 days from now
	var statusConditions []string

	for _, value := range values {
		switch strings.ToLower(value) {
		case "active":
			// Active: status == active AND expiration time is still far away (>10 days)
			condition := fmt.Sprintf("(active = $%d AND expiration > $%d)", *argCounter, *argCounter+1)
			statusConditions = append(statusConditions, condition)
			*args = append(*args, true, tenDaysFromNow)
			*argCounter += 2
		case "inactive":
			// Inactive: status == inactive AND expiration time is still future (not expired)
			// This excludes expired items which should only appear in "expired" filter
			condition := fmt.Sprintf("(active = $%d AND expiration > $%d)", *argCounter, *argCounter+1)
			statusConditions = append(statusConditions, condition)
			*args = append(*args, false, now)
			*argCounter += 2
		case "expired":
			// Expired: status == inactive AND expiration time is passed
			condition := fmt.Sprintf("(active = $%d AND expiration <= $%d)", *argCounter, *argCounter+1)
			statusConditions = append(statusConditions, condition)
			*args = append(*args, false, now)
			*argCounter += 2
		case "expiring_soon":
			// Expiring Soon: status == active AND expiration time is near (<=10 days)
			condition := fmt.Sprintf("(active = $%d AND expiration > $%d AND expiration <= $%d)", *argCounter, *argCounter+1, *argCounter+2)
			statusConditions = append(statusConditions, condition)
			*args = append(*args, true, now, tenDaysFromNow)
			*argCounter += 3
		}
	}

	if len(statusConditions) > 0 {
		*whereConditions = append(*whereConditions, "("+strings.Join(statusConditions, " OR ")+")")
	}
}

func (r *Repo) buildOrderClause(sorting *query.SortingParams) string {
	if sorting == nil {
		return " ORDER BY updated_at DESC" // Default sorting
	}

	// Map sorting field to database column
	var dbColumn string
	switch sorting.Field {
	case "name":
		dbColumn = "name"
	case "expiration":
		dbColumn = "expiration"
	case "createdAt", "created_at":
		dbColumn = "created_at"
	case "updatedAt", "updated_at":
		dbColumn = "updated_at"
	case "id":
		dbColumn = "id"
	default:
		dbColumn = "updated_at" // Fallback to default
	}

	// Determine sort direction
	direction := "ASC"
	if sorting.Order == query.SortOrderDesc {
		direction = "DESC"
	}

	return fmt.Sprintf(" ORDER BY %s %s", dbColumn, direction)
}

func (r *Repo) getDistinctValues(ctx context.Context, projectID int64) (map[string][]string, error) {
	filters := make(map[string][]string)

	// Get distinct names
	namesQuery := `
		SELECT DISTINCT name
		FROM altalune_project_api_keys
		WHERE project_id = $1 AND name IS NOT NULL
		ORDER BY name
	`
	names, err := r.queryDistinctValues(ctx, namesQuery, projectID)
	if err != nil {
		return nil, fmt.Errorf("get distinct names: %w", err)
	}
	filters["names"] = names

	// Set combined statuses (computed based on both active field and expiration)
	filters["statuses"] = []string{"active", "inactive", "expired", "expiring_soon"}

	return filters, nil
}

func (r *Repo) queryDistinctValues(ctx context.Context, query string, projectID int64) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := make([]string, 0) // Initialize as empty slice
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, rows.Err()
}

func (r *Repo) Create(ctx context.Context, input *CreateApiKeyInput) (*CreateApiKeyResult, error) {
	// Generate public ID
	publicID, _ := nanoid.GeneratePublicID()

	// Generate secure API key
	key, err := r.generateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("generate api key: %w", err)
	}

	// Check for name uniqueness within project
	nameCheckQuery := `SELECT COUNT(*) FROM altalune_project_api_keys WHERE project_id = $1 AND LOWER(name) = LOWER($2)`
	var count int
	err = r.db.QueryRowContext(ctx, nameCheckQuery, input.ProjectID, input.Name).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("check name uniqueness: %w", err)
	}
	if count > 0 {
		return nil, ErrApiKeyAlreadyExists
	}

	// Insert query
	insertQuery := `
		INSERT INTO altalune_project_api_keys (
			public_id,
			project_id,
			name,
			expiration,
			key,
			active,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	var result CreateApiKeyResult
	err = r.db.QueryRowContext(
		ctx,
		insertQuery,
		publicID,
		input.ProjectID,
		input.Name,
		input.Expiration,
		key,
		true, // New API keys are active by default
		now,
		now,
	).Scan(&result.ID, &result.CreatedAt, &result.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("insert api key: %w", err)
	}

	result.PublicID = publicID
	result.Name = input.Name
	result.Key = key
	result.Expiration = input.Expiration

	return &result, nil
}

func (r *Repo) generateAPIKey() (string, error) {
	// Generate 32 bytes of randomness
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Encode as base64 and add sk- prefix (OpenAI style)
	encoded := base64.URLEncoding.EncodeToString(randomBytes)
	// Remove padding and ensure consistent length
	encoded = strings.TrimRight(encoded, "=")

	return "sk-" + encoded, nil
}

func (r *Repo) GetByID(ctx context.Context, projectID int64, publicID string) (*ApiKey, error) {
	query := `
		SELECT
			public_id,
			name,
			expiration,
			active,
			created_at,
			updated_at
		FROM altalune_project_api_keys
		WHERE project_id = $1 AND public_id = $2
	`

	var result ApiKeyQueryResult
	err := r.db.QueryRowContext(ctx, query, projectID, publicID).Scan(
		&result.PublicID,
		&result.Name,
		&result.Expiration,
		&result.Active,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrApiKeyNotFound
		}
		return nil, fmt.Errorf("get api key by id: %w", err)
	}

	return result.ToApiKey(), nil
}

func (r *Repo) GetByKey(ctx context.Context, key string) (*ApiKey, error) {
	query := `
		SELECT
			public_id,
			name,
			expiration,
			active,
			created_at,
			updated_at
		FROM altalune_project_api_keys
		WHERE key = $1
	`

	var result ApiKeyQueryResult
	err := r.db.QueryRowContext(ctx, query, key).Scan(
		&result.PublicID,
		&result.Name,
		&result.Expiration,
		&result.Active,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrApiKeyNotFound
		}
		return nil, fmt.Errorf("get api key by key: %w", err)
	}

	return result.ToApiKey(), nil
}

func (r *Repo) Update(ctx context.Context, input *UpdateApiKeyInput) (*UpdateApiKeyResult, error) {
	// Check for name uniqueness within project (excluding current record)
	nameCheckQuery := `
		SELECT COUNT(*)
		FROM altalune_project_api_keys
		WHERE project_id = $1 AND LOWER(name) = LOWER($2) AND public_id != $3
	`
	var count int
	err := r.db.QueryRowContext(ctx, nameCheckQuery, input.ProjectID, input.Name, input.PublicID).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("check name uniqueness: %w", err)
	}
	if count > 0 {
		return nil, ErrApiKeyAlreadyExists
	}

	// Update query
	updateQuery := `
		UPDATE altalune_project_api_keys
		SET
			name = $1,
			expiration = $2,
			updated_at = $3
		WHERE project_id = $4 AND public_id = $5
		RETURNING id, active, created_at, updated_at
	`

	now := time.Now()
	var result UpdateApiKeyResult
	err = r.db.QueryRowContext(
		ctx,
		updateQuery,
		input.Name,
		input.Expiration,
		now,
		input.ProjectID,
		input.PublicID,
	).Scan(&result.ID, &result.Active, &result.CreatedAt, &result.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrApiKeyNotFound
		}
		return nil, fmt.Errorf("update api key: %w", err)
	}

	result.PublicID = input.PublicID
	result.Name = input.Name
	result.Expiration = input.Expiration

	return &result, nil
}

func (r *Repo) Delete(ctx context.Context, input *DeleteApiKeyInput) error {
	deleteQuery := `
		DELETE FROM altalune_project_api_keys
		WHERE project_id = $1 AND public_id = $2
	`

	result, err := r.db.ExecContext(ctx, deleteQuery, input.ProjectID, input.PublicID)
	if err != nil {
		return fmt.Errorf("delete api key: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrApiKeyNotFound
	}

	return nil
}

func (r *Repo) Activate(ctx context.Context, input *ActivateApiKeyInput) (*ActivateApiKeyResult, error) {
	// Set expiration to 1 year from now when reactivating (in case it was set to epoch time)
	now := time.Now()
	oneYearFromNow := now.AddDate(1, 0, 0)
	updateQuery := `
		UPDATE altalune_project_api_keys
		SET active = true, expiration = $1, updated_at = $2
		WHERE project_id = $3 AND public_id = $4
		RETURNING id, name, expiration, created_at, updated_at
	`

	var result ActivateApiKeyResult
	err := r.db.QueryRowContext(
		ctx,
		updateQuery,
		oneYearFromNow,
		now,
		input.ProjectID,
		input.PublicID,
	).Scan(&result.ID, &result.Name, &result.Expiration, &result.CreatedAt, &result.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrApiKeyNotFound
		}
		return nil, fmt.Errorf("activate api key: %w", err)
	}

	result.PublicID = input.PublicID
	result.Active = true

	return &result, nil
}

func (r *Repo) Deactivate(ctx context.Context, input *DeactivateApiKeyInput) (*DeactivateApiKeyResult, error) {
	// Set expiration to epoch time (1970-01-01) when deactivating
	epochTime := time.Unix(0, 0).UTC()
	updateQuery := `
		UPDATE altalune_project_api_keys
		SET active = false, expiration = $1, updated_at = $2
		WHERE project_id = $3 AND public_id = $4
		RETURNING id, name, expiration, created_at, updated_at
	`

	now := time.Now()
	var result DeactivateApiKeyResult
	err := r.db.QueryRowContext(
		ctx,
		updateQuery,
		epochTime,
		now,
		input.ProjectID,
		input.PublicID,
	).Scan(&result.ID, &result.Name, &result.Expiration, &result.CreatedAt, &result.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrApiKeyNotFound
		}
		return nil, fmt.Errorf("deactivate api key: %w", err)
	}

	result.PublicID = input.PublicID
	result.Active = false

	return &result, nil
}
