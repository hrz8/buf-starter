package user

import (
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
)

// mapUsersToProto converts slice of domain Users to proto Users
func mapUsersToProto(users []*User) []*altalunev1.User {
	if users == nil {
		return make([]*altalunev1.User, 0)
	}

	result := make([]*altalunev1.User, 0, len(users))
	for _, usr := range users {
		result = append(result, usr.ToUserProto())
	}
	return result
}

// mapFiltersToProto converts domain filters to proto FilterValues map
func mapFiltersToProto(filters map[string][]string) map[string]*altalunev1.FilterValues {
	if filters == nil {
		filters = make(map[string][]string)
	}

	result := make(map[string]*altalunev1.FilterValues)

	// Map is_active filter (boolean as strings: "true", "false")
	if activeStatuses, ok := filters["is_active"]; ok && activeStatuses != nil {
		result["is_active"] = &altalunev1.FilterValues{Values: activeStatuses}
	} else {
		result["is_active"] = &altalunev1.FilterValues{Values: []string{"true", "false"}}
	}

	return result
}
