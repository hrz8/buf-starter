package oauth_auth

import (
	"context"
)

// PermissionService adapts the IAM mapper repository to the PermissionFetcher interface.
type PermissionService struct {
	repo IAMMapperRepositor
}

// NewPermissionService creates a new permission fetcher adapter.
func NewPermissionService(repo IAMMapperRepositor) *PermissionService {
	return &PermissionService{repo: repo}
}

// GetUserPermissions fetches all permissions for a user and returns their names.
func (f *PermissionService) GetUserPermissions(ctx context.Context, userID int64) ([]string, error) {
	permissions, err := f.repo.GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(permissions))
	for _, p := range permissions {
		names = append(names, p.Name)
	}
	return names, nil
}
