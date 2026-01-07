package oauth_provider

import (
	"context"
	"strings"

	"buf.build/go/protovalidate"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"github.com/hrz8/altalune/internal/shared/query"
)

type Service struct {
	altalunev1.UnimplementedOAuthProviderServiceServer
	validator protovalidate.Validator
	logger    altalune.Logger
	repo      Repository
}

func NewService(validator protovalidate.Validator, logger altalune.Logger, repo Repository) *Service {
	return &Service{
		validator: validator,
		logger:    logger,
		repo:      repo,
	}
}

func (s *Service) QueryOAuthProviders(ctx context.Context, req *altalunev1.QueryOAuthProvidersRequest) (*altalunev1.QueryOAuthProvidersResponse, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Set default query if not provided
	if req.Query == nil {
		req.Query = &altalunev1.QueryRequest{
			Pagination: &altalunev1.Pagination{
				Page:     1,
				PageSize: 10,
			},
		}
	}

	// Convert proto request to domain query params
	queryParams := query.DefaultQueryParams(req.Query)

	// Query OAuth providers from repository
	result, err := s.repo.Query(ctx, queryParams)
	if err != nil {
		s.logger.Error("failed to query OAuth providers",
			"error", err,
			"keyword", queryParams.Keyword,
		)
		return nil, altalune.NewUnexpectedError("failed to query OAuth providers: %w", err)
	}

	// Convert domain result to proto response
	if result == nil {
		return &altalunev1.QueryOAuthProvidersResponse{
			Data: []*altalunev1.OAuthProvider{},
			Meta: &altalunev1.QueryMetaResponse{
				RowCount:  0,
				PageCount: 0,
				Filters:   make(map[string]*altalunev1.FilterValues),
			},
		}, nil
	}

	return &altalunev1.QueryOAuthProvidersResponse{
		Data: mapOAuthProvidersToProto(result.Data),
		Meta: &altalunev1.QueryMetaResponse{
			RowCount:  result.TotalRows,
			PageCount: result.TotalPages,
			Filters:   mapFiltersToProto(result.Filters),
		},
	}, nil
}

func (s *Service) CreateOAuthProvider(ctx context.Context, req *altalunev1.CreateOAuthProviderRequest) (*altalunev1.CreateOAuthProviderResponse, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Convert proto ProviderType to domain
	providerType := ProviderTypeFromProto(req.ProviderType)

	// Check for duplicate provider_type
	existingProvider, err := s.repo.GetByProviderType(ctx, providerType)
	if err != nil && err != ErrOAuthProviderNotFound {
		s.logger.Error("failed to check existing OAuth provider",
			"error", err,
			"provider_type", providerType,
		)
		return nil, altalune.NewUnexpectedError("failed to check existing OAuth provider: %w", err)
	}

	if existingProvider != nil {
		return nil, altalune.NewOAuthProviderDuplicateTypeError(string(providerType))
	}

	// Trim whitespace from inputs
	result, err := s.repo.Create(ctx, &CreateOAuthProviderInput{
		ProviderType: providerType,
		ClientID:     strings.TrimSpace(req.ClientId),
		ClientSecret: strings.TrimSpace(req.ClientSecret), // Plaintext, repo encrypts it
		RedirectURL:  strings.TrimSpace(req.RedirectUrl),
		Scopes:       strings.TrimSpace(req.Scopes),
		Enabled:      req.Enabled,
	})
	if err != nil {
		if err == ErrDuplicateProviderType {
			return nil, altalune.NewOAuthProviderDuplicateTypeError(string(providerType))
		}
		if err == ErrEncryptionFailed {
			return nil, altalune.NewOAuthProviderEncryptionError("")
		}
		s.logger.Error("failed to create OAuth provider",
			"error", err,
			"provider_type", providerType,
		)
		return nil, altalune.NewUnexpectedError("failed to create OAuth provider: %w", err)
	}

	s.logger.Info("OAuth provider created successfully",
		"provider_id", result.PublicID,
		"provider_type", providerType,
	)

	return &altalunev1.CreateOAuthProviderResponse{
		Provider: result.ToOAuthProvider().ToOAuthProviderProto(),
		Message:  "OAuth provider created successfully",
	}, nil
}

func (s *Service) GetOAuthProvider(ctx context.Context, req *altalunev1.GetOAuthProviderRequest) (*altalunev1.GetOAuthProviderResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	provider, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		if err == ErrOAuthProviderNotFound {
			return nil, altalune.NewOAuthProviderNotFoundError(req.Id)
		}
		s.logger.Error("failed to get OAuth provider", "error", err, "provider_id", req.Id)
		return nil, altalune.NewUnexpectedError("failed to get OAuth provider", err)
	}

	return &altalunev1.GetOAuthProviderResponse{
		Provider: provider.ToOAuthProviderProto(),
	}, nil
}

func (s *Service) UpdateOAuthProvider(ctx context.Context, req *altalunev1.UpdateOAuthProviderRequest) (*altalunev1.UpdateOAuthProviderResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Get existing provider to preserve provider_type and createdAt
	existingProvider, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		if err == ErrOAuthProviderNotFound {
			return nil, altalune.NewOAuthProviderNotFoundError(req.Id)
		}
		s.logger.Error("failed to get existing OAuth provider", "error", err, "provider_id", req.Id)
		return nil, altalune.NewUnexpectedError("failed to get existing OAuth provider", err)
	}

	// Trim whitespace from inputs
	input := &UpdateOAuthProviderInput{
		PublicID:     req.Id,
		ClientID:     strings.TrimSpace(req.ClientId),
		ClientSecret: strings.TrimSpace(req.ClientSecret), // Optional - if empty, repo retains existing secret
		RedirectURL:  strings.TrimSpace(req.RedirectUrl),
		Scopes:       strings.TrimSpace(req.Scopes),
		Enabled:      req.Enabled,
	}

	result, err := s.repo.Update(ctx, input)
	if err != nil {
		if err == ErrOAuthProviderNotFound {
			return nil, altalune.NewOAuthProviderNotFoundError(req.Id)
		}
		if err == ErrEncryptionFailed {
			return nil, altalune.NewOAuthProviderEncryptionError(req.Id)
		}
		s.logger.Error("failed to update OAuth provider", "error", err, "provider_id", req.Id)
		return nil, altalune.NewUnexpectedError("failed to update OAuth provider", err)
	}

	s.logger.Info("OAuth provider updated successfully",
		"provider_id", req.Id,
		"provider_type", existingProvider.ProviderType,
	)

	return &altalunev1.UpdateOAuthProviderResponse{
		Provider: result.ToOAuthProvider(existingProvider.ProviderType, existingProvider.CreatedAt).ToOAuthProviderProto(),
		Message:  "OAuth provider updated successfully",
	}, nil
}

func (s *Service) DeleteOAuthProvider(ctx context.Context, req *altalunev1.DeleteOAuthProviderRequest) (*altalunev1.DeleteOAuthProviderResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	err := s.repo.Delete(ctx, &DeleteOAuthProviderInput{
		PublicID: req.Id,
	})
	if err != nil {
		if err == ErrOAuthProviderNotFound {
			return nil, altalune.NewOAuthProviderNotFoundError(req.Id)
		}
		s.logger.Error("failed to delete OAuth provider", "error", err, "provider_id", req.Id)
		return nil, altalune.NewUnexpectedError("failed to delete OAuth provider", err)
	}

	s.logger.Info("OAuth provider deleted successfully", "provider_id", req.Id)

	return &altalunev1.DeleteOAuthProviderResponse{
		Message: "OAuth provider deleted successfully",
	}, nil
}

func (s *Service) RevealClientSecret(ctx context.Context, req *altalunev1.RevealClientSecretRequest) (*altalunev1.RevealClientSecretResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Call repo to decrypt and return plaintext client secret
	clientSecret, err := s.repo.RevealClientSecret(ctx, req.Id)
	if err != nil {
		if err == ErrOAuthProviderNotFound {
			return nil, altalune.NewOAuthProviderNotFoundError(req.Id)
		}
		if err == ErrDecryptionFailed {
			return nil, altalune.NewOAuthProviderDecryptionError(req.Id)
		}
		s.logger.Error("failed to reveal OAuth provider client secret", "error", err, "provider_id", req.Id)
		return nil, altalune.NewUnexpectedError("failed to reveal OAuth provider client secret", err)
	}

	// SECURITY AUDIT LOG: Log this explicit secret reveal action
	s.logger.Info("OAuth provider client secret revealed",
		"provider_id", req.Id,
		"action", "reveal_client_secret",
	)

	return &altalunev1.RevealClientSecretResponse{
		ClientSecret: clientSecret,
	}, nil
}
