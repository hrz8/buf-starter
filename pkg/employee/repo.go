package employee

import (
	"context"
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

func (r *Repo) Query(ctx context.Context, projectID int64, params *query.QueryParams) (*query.QueryResult[Employee], error) {
	// Build the base query
	baseQuery := `
		SELECT 
			id,
			public_id,
			name,
			email,
			role,
			department,
			status,
			created_at,
			updated_at
		FROM altalune_example_employees
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
			(LOWER(name) LIKE $%d OR 
			 LOWER(email) LIKE $%d OR 
			 LOWER(role) LIKE $%d OR 
			 LOWER(department) LIKE $%d OR 
			 LOWER(status) LIKE $%d)
		`, argCounter, argCounter, argCounter, argCounter, argCounter)
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
			case "role", "roles":
				dbColumn = "role"
			case "department", "departments":
				dbColumn = "department"
			case "status", "statuses":
				dbColumn = "status"
			case "name":
				dbColumn = "name"
			case "email":
				dbColumn = "email"
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
		return nil, fmt.Errorf("count employees: %w", err)
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
		return nil, fmt.Errorf("query employees: %w", err)
	}
	defer rows.Close()

	// Collect queryResults
	queryResults := make([]*EmployeeQueryResult, 0) // Initialize as empty slice, not nil
	for rows.Next() {
		var emp EmployeeQueryResult
		var status string
		err := rows.Scan(
			&emp.ID,
			&emp.PublicID,
			&emp.Name,
			&emp.Email,
			&emp.Role,
			&emp.Department,
			&status,
			&emp.CreatedAt,
			&emp.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan employee row: %w", err)
		}

		// Map status string to domain enum
		switch status {
		case "active":
			emp.Status = EmployeeStatusActive
		case "inactive":
			emp.Status = EmployeeStatusInactive
		default:
			emp.Status = EmployeeStatusActive // Default to active if unknown
		}

		queryResults = append(queryResults, &emp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate employee rows: %w", err)
	}

	// Calculate total pages (handle case where totalRows is 0)
	var totalPages int32
	if totalRows > 0 {
		totalPages = (totalRows + pageSize - 1) / pageSize
	}

	// Get distinct values for filters
	filters, err := r.getDistinctValues(ctx, projectID)
	if err != nil {
		// Don't fail the entire query if we can't get filters
		// Just log and return empty filters
		filters = make(map[string][]string)
	}

	results := make([]*Employee, 0)
	for _, v := range queryResults {
		results = append(results, &Employee{
			ID:         v.PublicID,
			Name:       v.Name,
			Email:      v.Email,
			Role:       v.Role,
			Department: v.Department,
			Status:     v.Status,
			CreatedAt:  v.CreatedAt,
			UpdatedAt:  v.UpdatedAt,
		})
	}

	return &query.QueryResult[Employee]{
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
	case "name":
		dbColumn = "name"
	case "email":
		dbColumn = "email"
	case "role":
		dbColumn = "role"
	case "department":
		dbColumn = "department"
	case "status":
		dbColumn = "status"
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

	// Get distinct roles
	rolesQuery := `
		SELECT DISTINCT role 
		FROM altalune_example_employees 
		WHERE project_id = $1 AND role IS NOT NULL
		ORDER BY role
	`
	roles, err := r.queryDistinctValues(ctx, rolesQuery, projectID)
	if err != nil {
		return nil, fmt.Errorf("get distinct roles: %w", err)
	}
	filters["roles"] = roles

	// Get distinct departments
	deptsQuery := `
		SELECT DISTINCT department 
		FROM altalune_example_employees 
		WHERE project_id = $1 AND department IS NOT NULL
		ORDER BY department
	`
	departments, err := r.queryDistinctValues(ctx, deptsQuery, projectID)
	if err != nil {
		return nil, fmt.Errorf("get distinct departments: %w", err)
	}
	filters["departments"] = departments

	// Get distinct statuses (these are constants)
	filters["statuses"] = []string{"active", "inactive"}

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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Empty results are valid
	return values, nil
}
