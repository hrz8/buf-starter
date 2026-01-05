package permission

import (
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
)

// mapPermissionsToProto converts slice of domain Permissions to proto Permissions
func mapPermissionsToProto(permissions []*Permission) []*altalunev1.Permission {
	if permissions == nil {
		return make([]*altalunev1.Permission, 0)
	}

	result := make([]*altalunev1.Permission, 0, len(permissions))
	for _, p := range permissions {
		result = append(result, p.ToPermissionProto())
	}
	return result
}

// mapFiltersToProto converts domain filters to proto FilterValues map
func mapFiltersToProto(filters map[string][]string) map[string]*altalunev1.FilterValues {
	if filters == nil {
		filters = make(map[string][]string)
	}

	result := make(map[string]*altalunev1.FilterValues)

	// Map effect filter
	if effects, ok := filters["effect"]; ok && effects != nil {
		result["effect"] = &altalunev1.FilterValues{Values: effects}
	} else {
		result["effect"] = &altalunev1.FilterValues{Values: []string{"allow", "deny"}}
	}

	return result
}
