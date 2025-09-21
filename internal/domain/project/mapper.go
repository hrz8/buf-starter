package project

import (
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
)

// mapProjectsToProto converts slice of domain Projects to proto Projects
func mapProjectsToProto(projects []*Project) []*altalunev1.Project {
	if projects == nil {
		return make([]*altalunev1.Project, 0)
	}

	result := make([]*altalunev1.Project, 0, len(projects))
	for _, prj := range projects {
		result = append(result, prj.ToProjectProto())
	}
	return result
}

// mapFiltersToProto converts domain filters to proto FilterValues map
func mapFiltersToProto(filters map[string][]string) map[string]*altalunev1.FilterValues {
	if filters == nil {
		filters = make(map[string][]string)
	}

	result := make(map[string]*altalunev1.FilterValues)

	// Map environments (always include as they're constants)
	if environments, ok := filters["environments"]; ok && environments != nil {
		result["environments"] = &altalunev1.FilterValues{Values: environments}
	} else {
		result["environments"] = &altalunev1.FilterValues{Values: []string{"live", "sandbox"}}
	}

	// Map timezones
	if timezones, ok := filters["timezones"]; ok && timezones != nil {
		result["timezones"] = &altalunev1.FilterValues{Values: timezones}
	} else {
		result["timezones"] = &altalunev1.FilterValues{Values: []string{}}
	}

	return result
}
