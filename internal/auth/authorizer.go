package auth

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
)

// RootPermission is the superadmin permission that bypasses all checks.
const RootPermission = "root"

// Authorizer provides authorization helper functions.
type Authorizer struct{}

// NewAuthorizer creates a new Authorizer.
func NewAuthorizer() *Authorizer {
	return &Authorizer{}
}

// CheckAuthenticated verifies user is authenticated.
func (a *Authorizer) CheckAuthenticated(ctx context.Context) error {
	auth := FromContext(ctx)
	if !auth.IsAuthenticated {
		return connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("missing authorization header"))
	}
	return nil
}

// CheckPermission validates user has required permission.
// Returns error if unauthorized.
func (a *Authorizer) CheckPermission(ctx context.Context, permission string) error {
	auth := FromContext(ctx)

	if !auth.IsAuthenticated {
		return connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("missing authorization header"))
	}

	// Superadmin bypass
	if a.IsSuperAdmin(ctx) {
		return nil
	}

	// Check permission
	for _, p := range auth.Permissions {
		if p == permission {
			return nil
		}
	}

	return connect.NewError(connect.CodePermissionDenied, fmt.Errorf("permission denied: requires %s", permission))
}

// CheckProjectAccess validates user has permission AND is project member.
// projectID is the project's public_id.
func (a *Authorizer) CheckProjectAccess(ctx context.Context, permission string, projectID string) error {
	auth := FromContext(ctx)

	if !auth.IsAuthenticated {
		return connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("missing authorization header"))
	}

	// Superadmin bypass
	if a.IsSuperAdmin(ctx) {
		return nil
	}

	// Check permission
	hasPermission := false
	for _, p := range auth.Permissions {
		if p == permission {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		return connect.NewError(connect.CodePermissionDenied, fmt.Errorf("permission denied: requires %s", permission))
	}

	// Check project membership
	if _, ok := auth.Memberships[projectID]; !ok {
		return connect.NewError(connect.CodePermissionDenied, fmt.Errorf("permission denied: not a member of this project"))
	}

	return nil
}

// CheckProjectMembership validates user is member of project.
func (a *Authorizer) CheckProjectMembership(ctx context.Context, projectID string) error {
	auth := FromContext(ctx)

	if !auth.IsAuthenticated {
		return connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("missing authorization header"))
	}

	// Superadmin bypass
	if a.IsSuperAdmin(ctx) {
		return nil
	}

	if _, ok := auth.Memberships[projectID]; !ok {
		return connect.NewError(connect.CodePermissionDenied, fmt.Errorf("permission denied: not a member of this project"))
	}

	return nil
}

// IsSuperAdmin checks if user has root permission.
func (a *Authorizer) IsSuperAdmin(ctx context.Context) bool {
	auth := FromContext(ctx)
	if !auth.IsAuthenticated {
		return false
	}

	for _, p := range auth.Permissions {
		if p == RootPermission {
			return true
		}
	}
	return false
}

// GetUserProjects returns list of project IDs user has access to.
func (a *Authorizer) GetUserProjects(ctx context.Context) []string {
	auth := FromContext(ctx)
	if !auth.IsAuthenticated {
		return nil
	}

	projects := make([]string, 0, len(auth.Memberships))
	for projectID := range auth.Memberships {
		projects = append(projects, projectID)
	}
	return projects
}

// GetProjectRole returns user's role in a specific project.
func (a *Authorizer) GetProjectRole(ctx context.Context, projectID string) string {
	auth := FromContext(ctx)
	if !auth.IsAuthenticated {
		return ""
	}
	return auth.Memberships[projectID]
}

// HasPermission checks if user has a specific permission.
func (a *Authorizer) HasPermission(ctx context.Context, permission string) bool {
	auth := FromContext(ctx)
	if !auth.IsAuthenticated {
		return false
	}

	// Superadmin has all permissions
	if a.IsSuperAdmin(ctx) {
		return true
	}

	for _, p := range auth.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// HasAnyPermission checks if user has any of the specified permissions.
func (a *Authorizer) HasAnyPermission(ctx context.Context, permissions []string) bool {
	auth := FromContext(ctx)
	if !auth.IsAuthenticated {
		return false
	}

	if a.IsSuperAdmin(ctx) {
		return true
	}

	for _, required := range permissions {
		for _, p := range auth.Permissions {
			if p == required {
				return true
			}
		}
	}
	return false
}

// CheckAnyPermission validates user has at least one of the required permissions (OR logic).
// Returns error if unauthorized.
func (a *Authorizer) CheckAnyPermission(ctx context.Context, permissions []string) error {
	auth := FromContext(ctx)

	if !auth.IsAuthenticated {
		return connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("missing authorization header"))
	}

	// Superadmin bypass
	if a.IsSuperAdmin(ctx) {
		return nil
	}

	// Check if user has any of the permissions
	for _, required := range permissions {
		for _, p := range auth.Permissions {
			if p == required {
				return nil
			}
		}
	}

	return connect.NewError(connect.CodePermissionDenied, fmt.Errorf("permission denied: requires one of %v", permissions))
}
