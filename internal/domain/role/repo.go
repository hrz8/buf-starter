package role

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
		FROM altalune_roles
		WHERE public_id = $1
	`

	var roleID int64
	err := r.db.QueryRowContext(ctx, query, publicID).Scan(&roleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrRoleNotFound
		}
		return 0, fmt.Errorf("get role ID by public ID: %w", err)
	}

	return roleID, nil
}

func (r *Repo) Query(ctx context.Context, params *query.QueryParams) (*query.QueryResult[Role], error) {
	// Build the base query - NO project_id filtering
	baseQuery := `
		SELECT
			id,
			public_id,
			name,
			description,
			created_at,
			updated_at
		FROM altalune_roles
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
		return nil, fmt.Errorf("count roles: %w", err)
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
		return nil, fmt.Errorf("query roles: %w", err)
	}
	defer rows.Close()

	// Collect queryResults
	queryResults := make([]*RoleQueryResult, 0)
	for rows.Next() {
		var role RoleQueryResult
		var description sql.NullString

		err := rows.Scan(
			&role.ID,
			&role.PublicID,
			&role.Name,
			&description,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan role row: %w", err)
		}

		// Handle nullable description
		if description.Valid {
			role.Description = description.String
		}

		queryResults = append(queryResults, &role)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate role rows: %w", err)
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
	results := make([]*Role, 0)
	for _, v := range queryResults {
		results = append(results, v.ToRole())
	}

	return &query.QueryResult[Role]{
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

	// Currently no specific filters for roles
	// Can be extended in the future

	return filters, nil
}

// Create creates a new role in the database
func (r *Repo) Create(ctx context.Context, input *CreateRoleInput) (*CreateRoleResult, error) {
	// Generate public ID
	publicID, _ := nanoid.GeneratePublicID()

	// Insert query - NO project_id
	insertQuery := `
		INSERT INTO altalune_roles (
			public_id,
			name,
			description,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id, public_id, name, description, created_at, updated_at
	`

	now := time.Now()
	var result CreateRoleResult
	var description sql.NullString

	err := r.db.QueryRowContext(
		ctx,
		insertQuery,
		publicID,
		input.Name,
		input.Description,
		now,
		now,
	).Scan(
		&result.ID,
		&result.PublicID,
		&result.Name,
		&description,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		// Check for unique constraint violation
		if postgres.IsUniqueViolation(err) {
			// Check if it's the name constraint
			if strings.Contains(err.Error(), "ux_altalune_roles_name") {
				return nil, ErrRoleAlreadyExists
			}
		}
		return nil, fmt.Errorf("create role: %w", err)
	}

	// Handle nullable description
	if description.Valid {
		result.Description = description.String
	}

	return &result, nil
}

// GetByName retrieves a role by name
func (r *Repo) GetByName(ctx context.Context, name string) (*Role, error) {
	query := `
		SELECT
			public_id,
			name,
			description,
			created_at,
			updated_at
		FROM altalune_roles
		WHERE LOWER(name) = LOWER($1)
		LIMIT 1
	`

	var role Role
	var description sql.NullString

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&role.ID,
		&role.Name,
		&description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRoleNotFound
		}
		return nil, fmt.Errorf("get role by name: %w", err)
	}

	// Handle nullable description
	if description.Valid {
		role.Description = description.String
	}

	return &role, nil
}

// GetByID retrieves a role by its public ID
func (r *Repo) GetByID(ctx context.Context, publicID string) (*Role, error) {
	sqlQuery := `
		SELECT
			public_id,
			name,
			description,
			created_at,
			updated_at
		FROM altalune_roles
		WHERE public_id = $1
	`

	var role Role
	var description sql.NullString

	err := r.db.QueryRowContext(ctx, sqlQuery, publicID).Scan(
		&role.ID,
		&role.Name,
		&description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRoleNotFound
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	// Handle nullable description
	if description.Valid {
		role.Description = description.String
	}

	return &role, nil
}

// Update updates a role in the database
func (r *Repo) Update(ctx context.Context, input *UpdateRoleInput) (*UpdateRoleResult, error) {
	// Check name uniqueness (exclude current role)
	existing, err := r.GetByName(ctx, input.Name)
	if err == nil && existing.ID != input.PublicID {
		return nil, ErrRoleAlreadyExists
	}

	sqlQuery := `
		UPDATE altalune_roles
		SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP
		WHERE public_id = $3
		RETURNING id, public_id, name, description, created_at, updated_at
	`

	var result UpdateRoleResult
	var description sql.NullString

	err = r.db.QueryRowContext(ctx, sqlQuery,
		input.Name,
		input.Description,
		input.PublicID,
	).Scan(
		&result.ID,
		&result.PublicID,
		&result.Name,
		&description,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRoleNotFound
		}
		if postgres.IsUniqueViolation(err) {
			return nil, ErrRoleAlreadyExists
		}
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	// Handle nullable description
	if description.Valid {
		result.Description = description.String
	}

	return &result, nil
}

// Delete deletes a role from the database
func (r *Repo) Delete(ctx context.Context, publicID string) error {
	sqlQuery := `DELETE FROM altalune_roles WHERE public_id = $1`

	result, err := r.db.ExecContext(ctx, sqlQuery, publicID)
	if err != nil {
		// Check for foreign key constraint violation
		if postgres.IsForeignKeyViolation(err) {
			return ErrRoleInUse
		}
		return fmt.Errorf("failed to delete role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrRoleNotFound
	}

	return nil
}
