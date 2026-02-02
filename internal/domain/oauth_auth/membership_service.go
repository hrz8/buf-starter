package oauth_auth

import (
	"context"

	"github.com/hrz8/altalune/internal/domain/iam_mapper"
)

// UserMembershipProvider defines the interface for fetching user project memberships.
type UserMembershipProvider interface {
	GetUserMemberships(ctx context.Context, userID int64) (map[string]string, error)
}

// MembershipRepositor defines the interface for fetching user project memberships from repository.
type MembershipRepositor interface {
	GetUserProjects(ctx context.Context, userID int64) ([]*iam_mapper.UserProjectMembership, error)
}

// MembershipService adapts the IAM mapper repository to the MembershipProvider interface.
type MembershipService struct {
	repo MembershipRepositor
}

// NewMembershipService creates a new membership fetcher adapter.
func NewMembershipService(repo MembershipRepositor) *MembershipService {
	return &MembershipService{repo: repo}
}

// GetUserMemberships fetches all project memberships for a user and returns map of project_public_id -> role.
func (s *MembershipService) GetUserMemberships(ctx context.Context, userID int64) (map[string]string, error) {
	projects, err := s.repo.GetUserProjects(ctx, userID)
	if err != nil {
		return nil, err
	}

	memberships := make(map[string]string, len(projects))
	for _, p := range projects {
		memberships[p.ProjectID] = p.Role
	}
	return memberships, nil
}
