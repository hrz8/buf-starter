package project

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/hrz8/altalune/internal/postgres"
	"github.com/hrz8/altalune/internal/query"
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
		results = append(results, &Project{
			ID:          v.PublicID,
			Name:        v.Name,
			Description: v.Description,
			Timezone:    v.Timezone,
			Environment: v.Environment,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		})
	}

	return &query.QueryResult[Project]{
		Data:       results,
		TotalRows:  totalRows,
		TotalPages: totalPages,
		Filters:    filters,
	}, nil
}

func (r *Repo) buildOrderClause(sorting *query.SortingParams) string {
	fmt.Println("KADIEE")
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
