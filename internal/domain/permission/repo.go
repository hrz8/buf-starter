package permission

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
		FROM altalune_permissions
		WHERE public_id = $1
	`

	var permissionID int64
	err := r.db.QueryRowContext(ctx, query, publicID).Scan(&permissionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrPermissionNotFound
		}
		return 0, fmt.Errorf("get permission ID by public ID: %w", err)
	}

	return permissionID, nil
}

func (r *Repo) Query(ctx context.Context, params *query.QueryParams) (*query.QueryResult[Permission], error) {
	// Build the base query - NO project_id filtering
	baseQuery := `
		SELECT
			id,
			public_id,
			name,
			effect,
			description,
			created_at,
			updated_at
		FROM altalune_permissions
		WHERE 1=1
	`

	// Build WHERE conditions for filters and search
	var whereConditions []string
	var args []interface{}
	argCounter := 1

	// Handle keyword search (search in name and description)
	if params.Keyword != "" {
		searchCondition := fmt.Sprintf(`
			(LOWER(name) LIKE $%d OR
			 LOWER(COALESCE(description, '')) LIKE $%d)
		`, argCounter, argCounter)
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
			case "effect":
				dbColumn = "effect"
			case "name":
				dbColumn = "name"
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
		return nil, fmt.Errorf("count permissions: %w", err)
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
		return nil, fmt.Errorf("query permissions: %w", err)
	}
	defer rows.Close()

	// Collect queryResults
	queryResults := make([]*PermissionQueryResult, 0)
	for rows.Next() {
		var perm PermissionQueryResult
		var description sql.NullString

		err := rows.Scan(
			&perm.ID,
			&perm.PublicID,
			&perm.Name,
			&perm.Effect,
			&description,
			&perm.CreatedAt,
			&perm.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan permission row: %w", err)
		}

		// Handle nullable description
		if description.Valid {
			perm.Description = description.String
		}

		queryResults = append(queryResults, &perm)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate permission rows: %w", err)
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
	results := make([]*Permission, 0)
	for _, v := range queryResults {
		results = append(results, v.ToPermission())
	}

	return &query.QueryResult[Permission]{
		Data:       results,
		TotalRows:  totalRows,
		TotalPages: totalPages,
		Filters:    filters,
	}, nil
}

func (r *Repo) buildOrderClause(sorting *query.SortingParams) string {
	if sorting == nil || sorting.Field == "" {
		return " ORDER BY name ASC" // Default sorting
	}

	// Map field to database column
	var dbColumn string
	switch sorting.Field {
	case "name":
		dbColumn = "name"
	case "effect":
		dbColumn = "effect"
	case "description":
		dbColumn = "description"
	case "createdAt", "created_at":
		dbColumn = "created_at"
	case "updatedAt", "updated_at":
		dbColumn = "updated_at"
	case "id":
		dbColumn = "id"
	default:
		dbColumn = "name" // Fallback to default
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

	// effect filter
	filters["effect"] = []string{"allow", "deny"}

	return filters, nil
}

// Create creates a new permission in the database
func (r *Repo) Create(ctx context.Context, input *CreatePermissionInput) (*CreatePermissionResult, error) {
	// Generate public ID
	publicID, _ := nanoid.GeneratePublicID()

	// Insert query - NO project_id
	insertQuery := `
		INSERT INTO altalune_permissions (
			public_id,
			name,
			effect,
			description,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, public_id, name, effect, description, created_at, updated_at
	`

	now := time.Now()
	var result CreatePermissionResult
	var description sql.NullString

	err := r.db.QueryRowContext(
		ctx,
		insertQuery,
		publicID,
		input.Name,
		input.Effect,
		input.Description,
		now,
		now,
	).Scan(
		&result.ID,
		&result.PublicID,
		&result.Name,
		&result.Effect,
		&description,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		// Check for unique constraint violation
		if postgres.IsUniqueViolation(err) {
			// Check if it's the name constraint
			if strings.Contains(err.Error(), "ux_altalune_permissions_name") {
				return nil, ErrPermissionAlreadyExists
			}
		}
		return nil, fmt.Errorf("create permission: %w", err)
	}

	// Handle nullable description
	if description.Valid {
		result.Description = description.String
	}

	return &result, nil
}

// GetByName retrieves a permission by name
func (r *Repo) GetByName(ctx context.Context, name string) (*Permission, error) {
	query := `
		SELECT
			public_id,
			name,
			effect,
			description,
			created_at,
			updated_at
		FROM altalune_permissions
		WHERE LOWER(name) = LOWER($1)
		LIMIT 1
	`

	var perm Permission
	var description sql.NullString

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&perm.ID,
		&perm.Name,
		&perm.Effect,
		&description,
		&perm.CreatedAt,
		&perm.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPermissionNotFound
		}
		return nil, fmt.Errorf("get permission by name: %w", err)
	}

	// Handle nullable description
	if description.Valid {
		perm.Description = description.String
	}

	return &perm, nil
}

// GetByID retrieves a permission by its public ID
func (r *Repo) GetByID(ctx context.Context, publicID string) (*Permission, error) {
	sqlQuery := `
		SELECT
			public_id,
			name,
			effect,
			description,
			created_at,
			updated_at
		FROM altalune_permissions
		WHERE public_id = $1
	`

	var perm Permission
	var description sql.NullString

	err := r.db.QueryRowContext(ctx, sqlQuery, publicID).Scan(
		&perm.ID,
		&perm.Name,
		&perm.Effect,
		&description,
		&perm.CreatedAt,
		&perm.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPermissionNotFound
		}
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}

	// Handle nullable description
	if description.Valid {
		perm.Description = description.String
	}

	return &perm, nil
}

// Update updates a permission in the database
func (r *Repo) Update(ctx context.Context, input *UpdatePermissionInput) (*UpdatePermissionResult, error) {
	// Check name uniqueness (exclude current permission)
	existing, err := r.GetByName(ctx, input.Name)
	if err == nil && existing.ID != input.PublicID {
		return nil, ErrPermissionAlreadyExists
	}

	sqlQuery := `
		UPDATE altalune_permissions
		SET name = $1, effect = $2, description = $3, updated_at = CURRENT_TIMESTAMP
		WHERE public_id = $4
		RETURNING id, public_id, name, effect, description, created_at, updated_at
	`

	var result UpdatePermissionResult
	var description sql.NullString

	err = r.db.QueryRowContext(ctx, sqlQuery,
		input.Name,
		input.Effect,
		input.Description,
		input.PublicID,
	).Scan(
		&result.ID,
		&result.PublicID,
		&result.Name,
		&result.Effect,
		&description,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPermissionNotFound
		}
		if postgres.IsUniqueViolation(err) {
			return nil, ErrPermissionAlreadyExists
		}
		return nil, fmt.Errorf("failed to update permission: %w", err)
	}

	// Handle nullable description
	if description.Valid {
		result.Description = description.String
	}

	return &result, nil
}

// Delete deletes a permission from the database
func (r *Repo) Delete(ctx context.Context, publicID string) error {
	sqlQuery := `DELETE FROM altalune_permissions WHERE public_id = $1`

	result, err := r.db.ExecContext(ctx, sqlQuery, publicID)
	if err != nil {
		// Check for foreign key constraint violation
		if postgres.IsForeignKeyViolation(err) {
			return ErrPermissionInUse
		}
		return fmt.Errorf("failed to delete permission: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrPermissionNotFound
	}

	return nil
}
