package user

import (
	"context"
	"database/sql"
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

func (r *Repo) GetIDByPublicID(ctx context.Context, publicID string) (int64, error) {
	query := `
		SELECT id
		FROM altalune_users
		WHERE public_id = $1
	`

	var userID int64
	err := r.db.QueryRowContext(ctx, query, publicID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrUserNotFound
		}
		return 0, fmt.Errorf("get user ID by public ID: %w", err)
	}

	return userID, nil
}

// GetInternalIDByEmail retrieves the internal user ID by email (case-insensitive)
func (r *Repo) GetInternalIDByEmail(ctx context.Context, email string) (int64, error) {
	query := `
		SELECT id
		FROM altalune_users
		WHERE LOWER(email) = LOWER($1)
	`

	var userID int64
	err := r.db.QueryRowContext(ctx, query, email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrUserNotFound
		}
		return 0, fmt.Errorf("get user ID by email: %w", err)
	}

	return userID, nil
}

func (r *Repo) Query(ctx context.Context, params *query.QueryParams) (*query.QueryResult[User], error) {
	// Build the base query - NO project_id filtering
	baseQuery := `
		SELECT
			id,
			public_id,
			email,
			first_name,
			last_name,
			is_active,
			email_verified,
			created_at,
			updated_at
		FROM altalune_users
		WHERE 1=1
	`

	// Build WHERE conditions for filters and search
	var whereConditions []string
	var args []interface{}
	argCounter := 1

	// Handle keyword search (search in email, first_name, last_name)
	if params.Keyword != "" {
		searchCondition := fmt.Sprintf(`
			(LOWER(email) LIKE $%d OR
			 LOWER(COALESCE(first_name, '')) LIKE $%d OR
			 LOWER(COALESCE(last_name, '')) LIKE $%d)
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
			case "is_active", "active":
				dbColumn = "is_active"
			case "email":
				dbColumn = "email"
			default:
				continue // Skip unknown fields
			}

			// Build IN clause for multiple values
			placeholders := make([]string, len(values))
			for i, value := range values {
				placeholders[i] = fmt.Sprintf("$%d", argCounter)

				// For boolean fields, convert string to boolean
				if dbColumn == "is_active" {
					boolValue := value == "true" || value == "1"
					args = append(args, boolValue)
				} else {
					args = append(args, strings.ToLower(value))
				}
				argCounter++
			}

			if dbColumn == "is_active" {
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
		return nil, fmt.Errorf("count users: %w", err)
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
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	// Collect queryResults
	queryResults := make([]*UserQueryResult, 0)
	for rows.Next() {
		var usr UserQueryResult
		var firstName, lastName sql.NullString

		err := rows.Scan(
			&usr.ID,
			&usr.PublicID,
			&usr.Email,
			&firstName,
			&lastName,
			&usr.IsActive,
			&usr.EmailVerified,
			&usr.CreatedAt,
			&usr.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan user row: %w", err)
		}

		// Handle nullable fields
		if firstName.Valid {
			usr.FirstName = firstName.String
		}
		if lastName.Valid {
			usr.LastName = lastName.String
		}

		queryResults = append(queryResults, &usr)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate user rows: %w", err)
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
	results := make([]*User, 0)
	for _, v := range queryResults {
		results = append(results, v.ToUser())
	}

	return &query.QueryResult[User]{
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
	case "email":
		dbColumn = "email"
	case "firstName", "first_name":
		dbColumn = "first_name"
	case "lastName", "last_name":
		dbColumn = "last_name"
	case "isActive", "is_active":
		dbColumn = "is_active"
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

	// is_active filter (boolean)
	filters["is_active"] = []string{"true", "false"}

	return filters, nil
}

// Create creates a new user in the database
func (r *Repo) Create(ctx context.Context, input *CreateUserInput) (*CreateUserResult, error) {
	// Generate public ID
	publicID, _ := nanoid.GeneratePublicID()

	// Email is already lowercased by service layer, but ensure it here too
	email := strings.ToLower(input.Email)

	// Determine is_active value: use input.IsActive if provided, otherwise default to true
	isActive := true
	if input.IsActive != nil {
		isActive = *input.IsActive
	}

	// Insert query - NO project_id
	insertQuery := `
		INSERT INTO altalune_users (
			public_id,
			email,
			first_name,
			last_name,
			is_active,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, public_id, email, first_name, last_name, is_active, email_verified, created_at, updated_at
	`

	now := time.Now()
	var result CreateUserResult
	var firstName, lastName sql.NullString

	err := r.db.QueryRowContext(
		ctx,
		insertQuery,
		publicID,
		email,
		input.FirstName,
		input.LastName,
		isActive,
		now,
		now,
	).Scan(
		&result.ID,
		&result.PublicID,
		&result.Email,
		&firstName,
		&lastName,
		&result.IsActive,
		&result.EmailVerified,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		// Check for unique constraint violation
		if postgres.IsUniqueViolation(err) {
			// Check if it's the email constraint
			if strings.Contains(err.Error(), "ux_altalune_users_email") {
				return nil, ErrUserAlreadyExists
			}
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	// Handle nullable fields
	if firstName.Valid {
		result.FirstName = firstName.String
	}
	if lastName.Valid {
		result.LastName = lastName.String
	}

	return &result, nil
}

// GetByEmail retrieves a user by email (case-insensitive)
func (r *Repo) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT
			public_id,
			email,
			first_name,
			last_name,
			is_active,
			email_verified,
			created_at,
			updated_at
		FROM altalune_users
		WHERE LOWER(email) = LOWER($1)
		LIMIT 1
	`

	var usr User
	var firstName, lastName sql.NullString

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&usr.ID,
		&usr.Email,
		&firstName,
		&lastName,
		&usr.IsActive,
		&usr.EmailVerified,
		&usr.CreatedAt,
		&usr.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	// Handle nullable fields
	if firstName.Valid {
		usr.FirstName = firstName.String
	}
	if lastName.Valid {
		usr.LastName = lastName.String
	}

	return &usr, nil
}

// GetByID retrieves a user by its public ID
func (r *Repo) GetByID(ctx context.Context, publicID string) (*User, error) {
	sqlQuery := `
		SELECT
			public_id,
			email,
			first_name,
			last_name,
			is_active,
			email_verified,
			created_at,
			updated_at
		FROM altalune_users
		WHERE public_id = $1
	`

	var usr User
	var firstName, lastName sql.NullString

	err := r.db.QueryRowContext(ctx, sqlQuery, publicID).Scan(
		&usr.ID,
		&usr.Email,
		&firstName,
		&lastName,
		&usr.IsActive,
		&usr.EmailVerified,
		&usr.CreatedAt,
		&usr.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Handle nullable fields
	if firstName.Valid {
		usr.FirstName = firstName.String
	}
	if lastName.Valid {
		usr.LastName = lastName.String
	}

	return &usr, nil
}

// GetByInternalID retrieves a user by internal database ID
func (r *Repo) GetByInternalID(ctx context.Context, internalID int64) (*User, error) {
	sqlQuery := `
		SELECT
			public_id,
			email,
			first_name,
			last_name,
			is_active,
			email_verified,
			created_at,
			updated_at
		FROM altalune_users
		WHERE id = $1
	`

	var usr User
	var firstName, lastName sql.NullString

	err := r.db.QueryRowContext(ctx, sqlQuery, internalID).Scan(
		&usr.ID,
		&usr.Email,
		&firstName,
		&lastName,
		&usr.IsActive,
		&usr.EmailVerified,
		&usr.CreatedAt,
		&usr.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by internal id: %w", err)
	}

	if firstName.Valid {
		usr.FirstName = firstName.String
	}
	if lastName.Valid {
		usr.LastName = lastName.String
	}

	return &usr, nil
}

// Update updates a user in the database
func (r *Repo) Update(ctx context.Context, input *UpdateUserInput) (*UpdateUserResult, error) {
	// Email is already lowercased by service layer, but ensure it here too
	email := strings.ToLower(input.Email)

	// Check email uniqueness (exclude current user)
	existing, err := r.GetByEmail(ctx, email)
	if err == nil && existing.ID != input.PublicID {
		return nil, ErrUserAlreadyExists
	}

	sqlQuery := `
		UPDATE altalune_users
		SET email = $1, first_name = $2, last_name = $3, updated_at = CURRENT_TIMESTAMP
		WHERE public_id = $4
		RETURNING id, public_id, email, first_name, last_name, is_active, email_verified,
		          created_at, updated_at
	`

	var result UpdateUserResult
	var firstName, lastName sql.NullString

	err = r.db.QueryRowContext(ctx, sqlQuery,
		email,
		input.FirstName,
		input.LastName,
		input.PublicID,
	).Scan(
		&result.ID,
		&result.PublicID,
		&result.Email,
		&firstName,
		&lastName,
		&result.IsActive,
		&result.EmailVerified,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		if postgres.IsUniqueViolation(err) {
			return nil, ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Handle nullable fields
	if firstName.Valid {
		result.FirstName = firstName.String
	}
	if lastName.Valid {
		result.LastName = lastName.String
	}

	return &result, nil
}

// Delete deletes a user from the database
func (r *Repo) Delete(ctx context.Context, publicID string) error {
	sqlQuery := `DELETE FROM altalune_users WHERE public_id = $1`

	result, err := r.db.ExecContext(ctx, sqlQuery, publicID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// Activate activates a user
func (r *Repo) Activate(ctx context.Context, publicID string) (*User, error) {
	// First check current state
	current, err := r.GetByID(ctx, publicID)
	if err != nil {
		return nil, err
	}

	if current.IsActive {
		return nil, ErrUserAlreadyActive
	}

	sqlQuery := `
		UPDATE altalune_users
		SET is_active = true, updated_at = CURRENT_TIMESTAMP
		WHERE public_id = $1
		RETURNING public_id, email, first_name, last_name, is_active, email_verified, created_at, updated_at
	`

	var usr User
	var firstName, lastName sql.NullString

	err = r.db.QueryRowContext(ctx, sqlQuery, publicID).Scan(
		&usr.ID,
		&usr.Email,
		&firstName,
		&lastName,
		&usr.IsActive,
		&usr.EmailVerified,
		&usr.CreatedAt,
		&usr.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to activate user: %w", err)
	}

	// Handle nullable fields
	if firstName.Valid {
		usr.FirstName = firstName.String
	}
	if lastName.Valid {
		usr.LastName = lastName.String
	}

	return &usr, nil
}

// Deactivate deactivates a user
func (r *Repo) Deactivate(ctx context.Context, publicID string) (*User, error) {
	// First check current state
	current, err := r.GetByID(ctx, publicID)
	if err != nil {
		return nil, err
	}

	if !current.IsActive {
		return nil, ErrUserAlreadyInactive
	}

	sqlQuery := `
		UPDATE altalune_users
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE public_id = $1
		RETURNING public_id, email, first_name, last_name, is_active, email_verified, created_at, updated_at
	`

	var usr User
	var firstName, lastName sql.NullString

	err = r.db.QueryRowContext(ctx, sqlQuery, publicID).Scan(
		&usr.ID,
		&usr.Email,
		&firstName,
		&lastName,
		&usr.IsActive,
		&usr.EmailVerified,
		&usr.CreatedAt,
		&usr.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to deactivate user: %w", err)
	}

	// Handle nullable fields
	if firstName.Valid {
		usr.FirstName = firstName.String
	}
	if lastName.Valid {
		usr.LastName = lastName.String
	}

	return &usr, nil
}
