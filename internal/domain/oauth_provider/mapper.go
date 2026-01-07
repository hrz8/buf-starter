package oauth_provider

import (
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
)

// ProviderTypeToProto converts domain ProviderType to proto ProviderType
func ProviderTypeToProto(pt ProviderType) altalunev1.ProviderType {
	switch pt {
	case ProviderTypeGoogle:
		return altalunev1.ProviderType_PROVIDER_TYPE_GOOGLE
	case ProviderTypeGithub:
		return altalunev1.ProviderType_PROVIDER_TYPE_GITHUB
	case ProviderTypeMicrosoft:
		return altalunev1.ProviderType_PROVIDER_TYPE_MICROSOFT
	case ProviderTypeApple:
		return altalunev1.ProviderType_PROVIDER_TYPE_APPLE
	default:
		return altalunev1.ProviderType_PROVIDER_TYPE_UNSPECIFIED
	}
}

// ProviderTypeFromProto converts proto ProviderType to domain ProviderType
func ProviderTypeFromProto(pt altalunev1.ProviderType) ProviderType {
	switch pt {
	case altalunev1.ProviderType_PROVIDER_TYPE_GOOGLE:
		return ProviderTypeGoogle
	case altalunev1.ProviderType_PROVIDER_TYPE_GITHUB:
		return ProviderTypeGithub
	case altalunev1.ProviderType_PROVIDER_TYPE_MICROSOFT:
		return ProviderTypeMicrosoft
	case altalunev1.ProviderType_PROVIDER_TYPE_APPLE:
		return ProviderTypeApple
	default:
		return "" // Empty string for unspecified
	}
}

// mapOAuthProvidersToProto converts slice of domain OAuthProviders to proto OAuthProviders
func mapOAuthProvidersToProto(providers []*OAuthProvider) []*altalunev1.OAuthProvider {
	if providers == nil {
		return make([]*altalunev1.OAuthProvider, 0)
	}

	result := make([]*altalunev1.OAuthProvider, 0, len(providers))
	for _, provider := range providers {
		result = append(result, provider.ToOAuthProviderProto())
	}
	return result
}

// mapFiltersToProto converts domain filters to proto FilterValues map
func mapFiltersToProto(filters map[string][]string) map[string]*altalunev1.FilterValues {
	if filters == nil {
		filters = make(map[string][]string)
	}

	result := make(map[string]*altalunev1.FilterValues)

	// Map provider_type filter
	if providerTypes, ok := filters["provider_type"]; ok && providerTypes != nil {
		result["provider_type"] = &altalunev1.FilterValues{Values: providerTypes}
	} else {
		result["provider_type"] = &altalunev1.FilterValues{
			Values: []string{"google", "github", "microsoft", "apple"},
		}
	}

	// Map enabled filter (boolean as strings: "true", "false")
	if enabledStatuses, ok := filters["enabled"]; ok && enabledStatuses != nil {
		result["enabled"] = &altalunev1.FilterValues{Values: enabledStatuses}
	} else {
		result["enabled"] = &altalunev1.FilterValues{Values: []string{"true", "false"}}
	}

	return result
}
