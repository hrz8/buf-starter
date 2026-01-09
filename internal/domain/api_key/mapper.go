package api_key

import (
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
)

// mapApiKeysToProto converts slice of domain ApiKeys to proto ApiKeys
func mapApiKeysToProto(apiKeys []*ApiKey) []*altalunev1.ApiKey {
	if apiKeys == nil {
		return make([]*altalunev1.ApiKey, 0)
	}

	result := make([]*altalunev1.ApiKey, 0, len(apiKeys))
	for _, key := range apiKeys {
		result = append(result, key.ToApiKeyProto())
	}
	return result
}

// mapFiltersToProto converts domain filters to proto FilterValues map
func mapFiltersToProto(filters map[string][]string) map[string]*altalunev1.FilterValues {
	if filters == nil {
		filters = make(map[string][]string)
	}

	result := make(map[string]*altalunev1.FilterValues)

	// Map names
	if names, ok := filters["names"]; ok && names != nil {
		result["names"] = &altalunev1.FilterValues{Values: names}
	} else {
		result["names"] = &altalunev1.FilterValues{Values: []string{}}
	}

	// Map combined statuses
	if statuses, ok := filters["statuses"]; ok && statuses != nil {
		result["statuses"] = &altalunev1.FilterValues{Values: statuses}
	} else {
		result["statuses"] = &altalunev1.FilterValues{Values: []string{"active", "inactive", "expired", "expiring_soon"}}
	}

	return result
}
