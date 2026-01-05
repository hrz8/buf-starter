package role

import (
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
)

// mapRolesToProto converts slice of domain Roles to proto Roles
func mapRolesToProto(roles []*Role) []*altalunev1.Role {
	if roles == nil {
		return make([]*altalunev1.Role, 0)
	}

	result := make([]*altalunev1.Role, 0, len(roles))
	for _, r := range roles {
		result = append(result, r.ToRoleProto())
	}
	return result
}

// mapFiltersToProto converts domain filters to proto FilterValues map
func mapFiltersToProto(filters map[string][]string) map[string]*altalunev1.FilterValues {
	if filters == nil {
		filters = make(map[string][]string)
	}

	result := make(map[string]*altalunev1.FilterValues)

	// Currently no specific filters for roles beyond keyword search
	// Can be extended in the future

	return result
}
