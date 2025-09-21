package project

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
		FROM altalune_projects 
		WHERE public_id = $1
	`

	var projectID int64
	err := r.db.QueryRowContext(ctx, query, publicID).Scan(&projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrProjectNotFound
		}
		return 0, fmt.Errorf("get project ID by public ID: %w", err)
	}

	return projectID, nil
}

func (r *Repo) Query(ctx context.Context, params *query.QueryParams) (*query.QueryResult[Project], error) {
	// Build the base query
	baseQuery := `
		SELECT 
			id,
			public_id,
			name,
			description,
			timezone,
			environment,
			created_at,
			updated_at
		FROM altalune_projects
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
			case "environment", "environments":
				dbColumn = "environment"
			case "timezone", "timezones":
				dbColumn = "timezone"
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
		return nil, fmt.Errorf("count projects: %w", err)
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
		return nil, fmt.Errorf("query projects: %w", err)
	}
	defer rows.Close()

	// Collect queryResults
	queryResults := make([]*ProjectQueryResult, 0)
	for rows.Next() {
		var prj ProjectQueryResult
		var description sql.NullString
		var environment string

		err := rows.Scan(
			&prj.ID,
			&prj.PublicID,
			&prj.Name,
			&description,
			&prj.Timezone,
			&environment,
			&prj.CreatedAt,
			&prj.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan project row: %w", err)
		}

		// Handle nullable description
		if description.Valid {
			prj.Description = description.String
		}

		// Map environment string to domain enum
		switch environment {
		case "live":
			prj.Environment = EnvironmentStatusLive
		case "sandbox":
			prj.Environment = EnvironmentStatusSandbox
		default:
			prj.Environment = EnvironmentStatusSandbox // Default to sandbox if unknown
		}

		queryResults = append(queryResults, &prj)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate project rows: %w", err)
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
	results := make([]*Project, 0)
	for _, v := range queryResults {
		results = append(results, v.ToProject())
	}

	return &query.QueryResult[Project]{
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
	case "environment":
		dbColumn = "environment"
	case "timezone":
		dbColumn = "timezone"
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

	// Get distinct environments (these are constants)
	filters["environments"] = []string{"live", "sandbox"}

	// Get distinct timezones
	tzQuery := `
		SELECT DISTINCT timezone 
		FROM altalune_projects 
		WHERE timezone IS NOT NULL
		ORDER BY timezone
	`
	timezones, err := r.queryDistinctValues(ctx, tzQuery)
	if err != nil {
		return nil, fmt.Errorf("get distinct timezones: %w", err)
	}
	filters["timezones"] = timezones

	return filters, nil
}

func (r *Repo) queryDistinctValues(ctx context.Context, query string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	values := make([]string, 0)
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return values, nil
}

// Create creates a new project in the database
func (r *Repo) Create(ctx context.Context, input *CreateProjectInput) (*CreateProjectResult, error) {
	// Generate public ID
	publicID, _ := nanoid.GeneratePublicID()

	// Map domain environment to database string
	var environmentStr string
	switch input.Environment {
	case EnvironmentStatusLive:
		environmentStr = "live"
	case EnvironmentStatusSandbox:
		environmentStr = "sandbox"
	default:
		environmentStr = "sandbox"
	}

	// Insert query
	insertQuery := `
		INSERT INTO altalune_projects (
			public_id,
			name,
			description,
			timezone,
			environment,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, public_id, name, description, timezone, environment, created_at, updated_at
	`

	now := time.Now()
	var result CreateProjectResult
	var description sql.NullString
	var returnedEnvironment string

	err := r.db.QueryRowContext(
		ctx,
		insertQuery,
		publicID,
		input.Name,
		input.Description,
		input.Timezone,
		environmentStr,
		now,
		now,
	).Scan(
		&result.ID,
		&result.PublicID,
		&result.Name,
		&description,
		&result.Timezone,
		&returnedEnvironment,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		// Check for unique constraint violation
		if postgres.IsUniqueViolation(err) {
			// Check if it's the name constraint
			if strings.Contains(err.Error(), "ux_altalune_projects_name") {
				return nil, ErrProjectAlreadyExists
			}
		}
		return nil, fmt.Errorf("create project: %w", err)
	}

	// Handle nullable description
	if description.Valid {
		result.Description = description.String
	}

	// Map environment back to domain enum
	switch returnedEnvironment {
	case "live":
		result.Environment = EnvironmentStatusLive
	case "sandbox":
		result.Environment = EnvironmentStatusSandbox
	default:
		result.Environment = EnvironmentStatusSandbox
	}

	// Create partitions for this new project
	if err := r.createPartitionsForProject(ctx, result.ID); err != nil {
		// Log error but don't fail the project creation
		// Partitions can be created manually if needed
		fmt.Printf("Warning: failed to create partitions for project %d: %v\n", result.ID, err)
	}

	return &result, nil
}

// GetByName retrieves a project by name
func (r *Repo) GetByName(ctx context.Context, name string) (*Project, error) {
	query := `
		SELECT
			public_id,
			name,
			description,
			timezone,
			environment,
			created_at,
			updated_at
		FROM altalune_projects
		WHERE LOWER(name) = LOWER($1)
		LIMIT 1
	`

	var prj Project
	var description sql.NullString
	var environment string

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&prj.ID,
		&prj.Name,
		&description,
		&prj.Timezone,
		&environment,
		&prj.CreatedAt,
		&prj.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProjectNotFound
		}
		return nil, fmt.Errorf("get project by name: %w", err)
	}

	// Handle nullable description
	if description.Valid {
		prj.Description = description.String
	}

	// Map environment string to domain enum
	switch environment {
	case "live":
		prj.Environment = EnvironmentStatusLive
	case "sandbox":
		prj.Environment = EnvironmentStatusSandbox
	default:
		prj.Environment = EnvironmentStatusSandbox
	}

	return &prj, nil
}

// partitionedTables defines all tables that need partitions created for new projects
// Add new tables here when they require partitioning by project_id
var partitionedTables = []string{
	"altalune_example_employees",
	"altalune_project_api_keys",
	// Add future partitioned tables here:
	// "altalune_project_logs",
	// "altalune_project_metrics",
	// etc.
}

// createPartitionsForProject creates the necessary database partitions for a new project
func (r *Repo) createPartitionsForProject(ctx context.Context, projectID int64) error {
	// Create partitions for all configured tables
	for _, tableName := range partitionedTables {
		partitionName := fmt.Sprintf("%s_p%d", tableName, projectID)
		query := fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s
			PARTITION OF %s FOR VALUES IN (%d)
		`, partitionName, tableName, projectID)

		if _, err := r.db.ExecContext(ctx, query); err != nil {
			// Log the specific table that failed but continue with others
			fmt.Printf("Warning: failed to create partition %s: %v\n", partitionName, err)
			// Don't return error here - we want to try creating all partitions
		}
	}

	return nil
}
