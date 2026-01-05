package iam_mapper

import (
	"context"
	"fmt"
	"strings"

	"github.com/hrz8/altalune/internal/domain/permission"
	"github.com/hrz8/altalune/internal/domain/role"
	"github.com/hrz8/altalune/internal/postgres"
	"github.com/lib/pq"
)

type Repo struct {
	db postgres.DB
}

func NewRepo(db postgres.DB) *Repo {
	return &Repo{db: db}
}

// ==================== User-Role Mappings ====================

func (r *Repo) AssignUserRoles(ctx context.Context, userID int64, roleIDs []int64) error {
	if len(roleIDs) == 0 {
		return nil
	}

	// Build multi-row INSERT with ON CONFLICT DO NOTHING for idempotency
	query := `INSERT INTO altalune_users_roles (user_id, role_id) VALUES `

	args := []interface{}{userID}
	placeholders := []string{}

	for i, roleID := range roleIDs {
		placeholders = append(placeholders, fmt.Sprintf("($1, $%d)", i+2))
		args = append(args, roleID)
	}

	query += strings.Join(placeholders, ", ") + " ON CONFLICT (user_id, role_id) DO NOTHING"

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("assign user roles: %w", err)
	}

	return nil
}

func (r *Repo) RemoveUserRoles(ctx context.Context, userID int64, roleIDs []int64) error {
	if len(roleIDs) == 0 {
		return nil
	}

	query := `
		DELETE FROM altalune_users_roles
		WHERE user_id = $1 AND role_id = ANY($2)
	`

	_, err := r.db.ExecContext(ctx, query, userID, pq.Array(roleIDs))
	if err != nil {
		return fmt.Errorf("remove user roles: %w", err)
	}

	return nil
}

func (r *Repo) GetUserRoles(ctx context.Context, userID int64) ([]*role.Role, error) {
	query := `
		SELECT r.id, r.public_id, r.name, r.description, r.created_at, r.updated_at
		FROM altalune_roles r
		INNER JOIN altalune_users_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1
		ORDER BY r.name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get user roles: %w", err)
	}
	defer rows.Close()

	var roles []*role.Role
	for rows.Next() {
		var result RoleQueryResult
		err := rows.Scan(
			&result.ID,
			&result.PublicID,
			&result.Name,
			&result.Description,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan role: %w", err)
		}
		roles = append(roles, result.ToRole())
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return roles, nil
}

// ==================== Role-Permission Mappings ====================

func (r *Repo) AssignRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error {
	if len(permissionIDs) == 0 {
		return nil
	}

	// Build multi-row INSERT with ON CONFLICT DO NOTHING for idempotency
	query := `INSERT INTO altalune_roles_permissions (role_id, permission_id) VALUES `

	args := []interface{}{roleID}
	placeholders := []string{}

	for i, permissionID := range permissionIDs {
		placeholders = append(placeholders, fmt.Sprintf("($1, $%d)", i+2))
		args = append(args, permissionID)
	}

	query += strings.Join(placeholders, ", ") + " ON CONFLICT (role_id, permission_id) DO NOTHING"

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("assign role permissions: %w", err)
	}

	return nil
}

func (r *Repo) RemoveRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error {
	if len(permissionIDs) == 0 {
		return nil
	}

	query := `
		DELETE FROM altalune_roles_permissions
		WHERE role_id = $1 AND permission_id = ANY($2)
	`

	_, err := r.db.ExecContext(ctx, query, roleID, pq.Array(permissionIDs))
	if err != nil {
		return fmt.Errorf("remove role permissions: %w", err)
	}

	return nil
}

func (r *Repo) GetRolePermissions(ctx context.Context, roleID int64) ([]*permission.Permission, error) {
	query := `
		SELECT p.id, p.public_id, p.name, p.effect, p.description, p.created_at, p.updated_at
		FROM altalune_permissions p
		INNER JOIN altalune_roles_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id = $1
		ORDER BY p.name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, roleID)
	if err != nil {
		return nil, fmt.Errorf("get role permissions: %w", err)
	}
	defer rows.Close()

	var permissions []*permission.Permission
	for rows.Next() {
		var result PermissionQueryResult
		err := rows.Scan(
			&result.ID,
			&result.PublicID,
			&result.Name,
			&result.Effect,
			&result.Description,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan permission: %w", err)
		}
		permissions = append(permissions, result.ToPermission())
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return permissions, nil
}

// ==================== User-Permission Mappings (Direct) ====================

func (r *Repo) AssignUserPermissions(ctx context.Context, userID int64, permissionIDs []int64) error {
	if len(permissionIDs) == 0 {
		return nil
	}

	// Build multi-row INSERT with ON CONFLICT DO NOTHING for idempotency
	query := `INSERT INTO altalune_users_permissions (user_id, permission_id) VALUES `

	args := []interface{}{userID}
	placeholders := []string{}

	for i, permissionID := range permissionIDs {
		placeholders = append(placeholders, fmt.Sprintf("($1, $%d)", i+2))
		args = append(args, permissionID)
	}

	query += strings.Join(placeholders, ", ") + " ON CONFLICT (user_id, permission_id) DO NOTHING"

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("assign user permissions: %w", err)
	}

	return nil
}

func (r *Repo) RemoveUserPermissions(ctx context.Context, userID int64, permissionIDs []int64) error {
	if len(permissionIDs) == 0 {
		return nil
	}

	query := `
		DELETE FROM altalune_users_permissions
		WHERE user_id = $1 AND permission_id = ANY($2)
	`

	_, err := r.db.ExecContext(ctx, query, userID, pq.Array(permissionIDs))
	if err != nil {
		return fmt.Errorf("remove user permissions: %w", err)
	}

	return nil
}

func (r *Repo) GetUserPermissions(ctx context.Context, userID int64) ([]*permission.Permission, error) {
	query := `
		SELECT DISTINCT p.id, p.public_id, p.name, p.effect, p.description, p.created_at, p.updated_at
		FROM altalune_permissions p
		INNER JOIN altalune_users_permissions up ON up.permission_id = p.id
		WHERE up.user_id = $1

		UNION

		SELECT DISTINCT p.id, p.public_id, p.name, p.effect, p.description, p.created_at, p.updated_at
		FROM altalune_permissions p
		INNER JOIN altalune_roles_permissions rp ON rp.permission_id = p.id
		INNER JOIN altalune_users_roles ur ON ur.role_id = rp.role_id
		WHERE ur.user_id = $1

		ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get user permissions: %w", err)
	}
	defer rows.Close()

	var permissions []*permission.Permission
	for rows.Next() {
		var result PermissionQueryResult
		err := rows.Scan(
			&result.ID,
			&result.PublicID,
			&result.Name,
			&result.Effect,
			&result.Description,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan permission: %w", err)
		}
		permissions = append(permissions, result.ToPermission())
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return permissions, nil
}

// ==================== Project Members ====================

func (r *Repo) AssignProjectMembers(ctx context.Context, projectID int64, members []ProjectMemberInput) error {
	if len(members) == 0 {
		return nil
	}

	// Build multi-row INSERT with ON CONFLICT to update role
	query := `
		INSERT INTO altalune_project_members (project_id, user_id, role)
		VALUES `

	args := []interface{}{projectID}
	placeholders := []string{}
	argCounter := 2

	for _, member := range members {
		placeholders = append(placeholders, fmt.Sprintf("($1, $%d, $%d)", argCounter, argCounter+1))
		args = append(args, member.UserID, member.Role)
		argCounter += 2
	}

	query += strings.Join(placeholders, ", ") + `
		ON CONFLICT (project_id, user_id)
		DO UPDATE SET role = EXCLUDED.role, updated_at = NOW()
	`

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("assign project members: %w", err)
	}

	return nil
}

func (r *Repo) RemoveProjectMembers(ctx context.Context, projectID int64, userIDs []int64) error {
	if len(userIDs) == 0 {
		return nil
	}

	// First, check if any of the users to remove are owners
	// and ensure we're not removing the last owner
	checkQuery := `
		SELECT COUNT(*) as owner_count
		FROM altalune_project_members
		WHERE project_id = $1 AND role = 'owner' AND user_id = ANY($2)
	`

	var ownersToRemove int
	err := r.db.QueryRowContext(ctx, checkQuery, projectID, pq.Array(userIDs)).Scan(&ownersToRemove)
	if err != nil {
		return fmt.Errorf("check owners to remove: %w", err)
	}

	// If we're removing owners, check total owner count
	if ownersToRemove > 0 {
		countQuery := `
			SELECT COUNT(*) as total_owners
			FROM altalune_project_members
			WHERE project_id = $1 AND role = 'owner'
		`

		var totalOwners int
		err = r.db.QueryRowContext(ctx, countQuery, projectID).Scan(&totalOwners)
		if err != nil {
			return fmt.Errorf("count total owners: %w", err)
		}

		// Prevent removing all owners
		if totalOwners-ownersToRemove < 1 {
			return ErrCannotRemoveLastOwner
		}
	}

	// Proceed with removal
	query := `
		DELETE FROM altalune_project_members
		WHERE project_id = $1 AND user_id = ANY($2)
	`

	_, err = r.db.ExecContext(ctx, query, projectID, pq.Array(userIDs))
	if err != nil {
		return fmt.Errorf("remove project members: %w", err)
	}

	return nil
}

func (r *Repo) GetProjectMembers(ctx context.Context, projectID int64) ([]*ProjectMemberWithUser, error) {
	query := `
		SELECT
			u.id as user_id,
			u.public_id as user_public_id,
			u.email,
			u.first_name,
			u.last_name,
			u.is_active,
			u.created_at as user_created_at,
			u.updated_at as user_updated_at,
			pm.role,
			pm.created_at
		FROM altalune_project_members pm
		INNER JOIN altalune_users u ON u.id = pm.user_id
		WHERE pm.project_id = $1
		ORDER BY pm.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("get project members: %w", err)
	}
	defer rows.Close()

	var members []*ProjectMemberWithUser
	for rows.Next() {
		var result ProjectMemberQueryResult
		err := rows.Scan(
			&result.UserID,
			&result.UserPublicID,
			&result.Email,
			&result.FirstName,
			&result.LastName,
			&result.IsActive,
			&result.UserCreatedAt,
			&result.UserUpdatedAt,
			&result.Role,
			&result.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan project member: %w", err)
		}
		members = append(members, result.ToProjectMemberWithUser())
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return members, nil
}
