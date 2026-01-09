package api_key

import (
	"context"
	"time"

	"buf.build/go/protovalidate"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	project_domain "github.com/hrz8/altalune/internal/domain/project"
	"github.com/hrz8/altalune/internal/shared/query"
)

type Service struct {
	altalunev1.UnimplementedApiKeyServiceServer
	validator   protovalidate.Validator
	log         altalune.Logger
	projectRepo project_domain.Repositor
	apiKeyRepo  Repositor
}

func NewService(v protovalidate.Validator, log altalune.Logger, projectRepo project_domain.Repositor, apiKeyRepo Repositor) *Service {
	return &Service{
		validator:   v,
		log:         log,
		projectRepo: projectRepo,
		apiKeyRepo:  apiKeyRepo,
	}
}

func (s *Service) QueryApiKeys(ctx context.Context, req *altalunev1.QueryApiKeysRequest) (*altalunev1.QueryApiKeysResponse, error) {
	// Add artificial delay to match TypeScript version
	time.Sleep(600 * time.Millisecond)

	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Validate that query is provided
	if req.Query == nil {
		return nil, altalune.NewInvalidPayloadError("query is required")
	}

	// Extract project ID
	projectID, err := s.projectRepo.GetIDByPublicID(ctx, req.ProjectId)
	if err != nil {
		if err == project_domain.ErrProjectNotFound {
			return nil, altalune.NewProjectNotFound(req.ProjectId)
		}
		return nil, altalune.NewInvalidPayloadError("invalid project_id")
	}

	// Convert proto request to domain query params
	queryParams := query.DefaultQueryParams(req.Query)

	// Query API keys from repository
	result, err := s.apiKeyRepo.Query(ctx, projectID, queryParams)
	if err != nil {
		s.log.Error("failed to query api keys",
			"error", err,
			"project_id", projectID,
			"keyword", queryParams.Keyword,
		)
		return nil, altalune.NewUnexpectedError("failed to query api keys: %w", err)
	}

	// Convert domain result to proto response
	if result == nil {
		return &altalunev1.QueryApiKeysResponse{
			Data: []*altalunev1.ApiKey{},
			Meta: &altalunev1.QueryMetaResponse{
				RowCount:  0,
				PageCount: 0,
				Filters:   make(map[string]*altalunev1.FilterValues),
			},
		}, nil
	}

	return &altalunev1.QueryApiKeysResponse{
		Data: mapApiKeysToProto(result.Data),
		Meta: &altalunev1.QueryMetaResponse{
			RowCount:  result.TotalRows,
			PageCount: result.TotalPages,
			Filters:   mapFiltersToProto(result.Filters),
		},
	}, nil
}

func (s *Service) CreateApiKey(ctx context.Context, req *altalunev1.CreateApiKeyRequest) (*altalunev1.CreateApiKeyResponse, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Extract and validate project ID
	projectID, err := s.projectRepo.GetIDByPublicID(ctx, req.ProjectId)
	if err != nil {
		if err == project_domain.ErrProjectNotFound {
			return nil, altalune.NewProjectNotFound(req.ProjectId)
		}
		return nil, altalune.NewInvalidPayloadError("invalid project_id")
	}

	// Validate expiration date
	expiration := req.Expiration.AsTime()
	now := time.Now()
	if expiration.Before(now) {
		return nil, altalune.NewInvalidPayloadError("expiration must be in the future")
	}
	maxExpiration := now.AddDate(2, 0, 0) // 2 years from now
	if expiration.After(maxExpiration) {
		return nil, altalune.NewInvalidPayloadError("expiration cannot be more than 2 years in the future")
	}

	// Prepare input for repository
	input := &CreateApiKeyInput{
		ProjectID:  projectID,
		Name:       req.Name,
		Expiration: expiration,
	}

	// Create API key
	result, err := s.apiKeyRepo.Create(ctx, input)
	if err != nil {
		if err == ErrApiKeyAlreadyExists {
			return nil, altalune.NewApiKeyAlreadyExistsError(req.Name)
		}
		s.log.Error("failed to create api key",
			"error", err,
			"project_id", projectID,
			"name", req.Name,
		)
		return nil, altalune.NewUnexpectedError("failed to create api key: %w", err)
	}

	// Log successful creation for audit purposes
	s.log.Info("api key created",
		"project_id", projectID,
		"api_key_id", result.PublicID,
		"name", result.Name,
		"expiration", result.Expiration,
	)

	// Convert to domain model (without the actual key for security)
	apiKey := &ApiKey{
		ID:         result.PublicID,
		Name:       result.Name,
		Expiration: result.Expiration,
		Active:     true, // New API keys are active by default
		CreatedAt:  result.CreatedAt,
		UpdatedAt:  result.UpdatedAt,
	}

	return &altalunev1.CreateApiKeyResponse{
		ApiKey:   apiKey.ToApiKeyProto(),
		KeyValue: result.Key, // Only returned once during creation
		Message:  "API key created successfully",
	}, nil
}

func (s *Service) GetApiKey(ctx context.Context, req *altalunev1.GetApiKeyRequest) (*altalunev1.GetApiKeyResponse, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Extract and validate project ID
	projectID, err := s.projectRepo.GetIDByPublicID(ctx, req.ProjectId)
	if err != nil {
		if err == project_domain.ErrProjectNotFound {
			return nil, altalune.NewProjectNotFound(req.ProjectId)
		}
		return nil, altalune.NewInvalidPayloadError("invalid project_id")
	}

	// Get API key from repository
	apiKey, err := s.apiKeyRepo.GetByID(ctx, projectID, req.ApiKeyId)
	if err != nil {
		if err == ErrApiKeyNotFound {
			return nil, altalune.NewApiKeyNotFoundError(req.ApiKeyId)
		}
		s.log.Error("failed to get api key",
			"error", err,
			"project_id", projectID,
			"api_key_id", req.ApiKeyId,
		)
		return nil, altalune.NewUnexpectedError("failed to get api key: %w", err)
	}

	return &altalunev1.GetApiKeyResponse{
		ApiKey: apiKey.ToApiKeyProto(),
	}, nil
}

func (s *Service) UpdateApiKey(ctx context.Context, req *altalunev1.UpdateApiKeyRequest) (*altalunev1.UpdateApiKeyResponse, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Extract and validate project ID
	projectID, err := s.projectRepo.GetIDByPublicID(ctx, req.ProjectId)
	if err != nil {
		if err == project_domain.ErrProjectNotFound {
			return nil, altalune.NewProjectNotFound(req.ProjectId)
		}
		return nil, altalune.NewInvalidPayloadError("invalid project_id")
	}

	// Validate expiration date
	expiration := req.Expiration.AsTime()
	now := time.Now()
	if expiration.Before(now) {
		return nil, altalune.NewInvalidPayloadError("expiration must be in the future")
	}
	maxExpiration := now.AddDate(2, 0, 0) // 2 years from now
	if expiration.After(maxExpiration) {
		return nil, altalune.NewInvalidPayloadError("expiration cannot be more than 2 years in the future")
	}

	// Prepare input for repository
	input := &UpdateApiKeyInput{
		ProjectID:  projectID,
		PublicID:   req.ApiKeyId,
		Name:       req.Name,
		Expiration: expiration,
	}

	// Update API key
	result, err := s.apiKeyRepo.Update(ctx, input)
	if err != nil {
		if err == ErrApiKeyNotFound {
			return nil, altalune.NewApiKeyNotFoundError(req.ApiKeyId)
		}
		if err == ErrApiKeyAlreadyExists {
			return nil, altalune.NewApiKeyAlreadyExistsError(req.Name)
		}
		s.log.Error("failed to update api key",
			"error", err,
			"project_id", projectID,
			"api_key_id", req.ApiKeyId,
		)
		return nil, altalune.NewUnexpectedError("failed to update api key: %w", err)
	}

	// Log successful update for audit purposes
	s.log.Info("api key updated",
		"project_id", projectID,
		"api_key_id", result.PublicID,
		"name", result.Name,
		"expiration", result.Expiration,
	)

	// Convert to domain model
	apiKey := &ApiKey{
		ID:         result.PublicID,
		Name:       result.Name,
		Expiration: result.Expiration,
		Active:     result.Active,
		CreatedAt:  result.CreatedAt,
		UpdatedAt:  result.UpdatedAt,
	}

	return &altalunev1.UpdateApiKeyResponse{
		ApiKey:  apiKey.ToApiKeyProto(),
		Message: "API key updated successfully",
	}, nil
}

func (s *Service) DeleteApiKey(ctx context.Context, req *altalunev1.DeleteApiKeyRequest) (*altalunev1.DeleteApiKeyResponse, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Extract and validate project ID
	projectID, err := s.projectRepo.GetIDByPublicID(ctx, req.ProjectId)
	if err != nil {
		if err == project_domain.ErrProjectNotFound {
			return nil, altalune.NewProjectNotFound(req.ProjectId)
		}
		return nil, altalune.NewInvalidPayloadError("invalid project_id")
	}

	// Prepare input for repository
	input := &DeleteApiKeyInput{
		ProjectID: projectID,
		PublicID:  req.ApiKeyId,
	}

	// Delete API key
	err = s.apiKeyRepo.Delete(ctx, input)
	if err != nil {
		if err == ErrApiKeyNotFound {
			return nil, altalune.NewApiKeyNotFoundError(req.ApiKeyId)
		}
		s.log.Error("failed to delete api key",
			"error", err,
			"project_id", projectID,
			"api_key_id", req.ApiKeyId,
		)
		return nil, altalune.NewUnexpectedError("failed to delete api key: %w", err)
	}

	// Log successful deletion for audit purposes
	s.log.Info("api key deleted",
		"project_id", projectID,
		"api_key_id", req.ApiKeyId,
	)

	return &altalunev1.DeleteApiKeyResponse{
		Message: "API key deleted successfully",
	}, nil
}

func (s *Service) ActivateApiKey(ctx context.Context, req *altalunev1.ActivateApiKeyRequest) (*altalunev1.ActivateApiKeyResponse, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Extract and validate project ID
	projectID, err := s.projectRepo.GetIDByPublicID(ctx, req.ProjectId)
	if err != nil {
		if err == project_domain.ErrProjectNotFound {
			return nil, altalune.NewProjectNotFound(req.ProjectId)
		}
		return nil, altalune.NewInvalidPayloadError("invalid project_id")
	}

	// Prepare input for repository
	input := &ActivateApiKeyInput{
		ProjectID: projectID,
		PublicID:  req.ApiKeyId,
	}

	// Activate API key
	result, err := s.apiKeyRepo.Activate(ctx, input)
	if err != nil {
		if err == ErrApiKeyNotFound {
			return nil, altalune.NewApiKeyNotFoundError(req.ApiKeyId)
		}
		s.log.Error("failed to activate api key",
			"error", err,
			"project_id", projectID,
			"api_key_id", req.ApiKeyId,
		)
		return nil, altalune.NewUnexpectedError("failed to activate api key: %w", err)
	}

	// Log successful activation for audit purposes
	s.log.Info("api key activated",
		"project_id", projectID,
		"api_key_id", result.PublicID,
		"name", result.Name,
	)

	// Convert to domain model
	apiKey := &ApiKey{
		ID:         result.PublicID,
		Name:       result.Name,
		Expiration: result.Expiration,
		Active:     result.Active,
		CreatedAt:  result.CreatedAt,
		UpdatedAt:  result.UpdatedAt,
	}

	return &altalunev1.ActivateApiKeyResponse{
		ApiKey:  apiKey.ToApiKeyProto(),
		Message: "API key activated successfully",
	}, nil
}

func (s *Service) DeactivateApiKey(ctx context.Context, req *altalunev1.DeactivateApiKeyRequest) (*altalunev1.DeactivateApiKeyResponse, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Extract and validate project ID
	projectID, err := s.projectRepo.GetIDByPublicID(ctx, req.ProjectId)
	if err != nil {
		if err == project_domain.ErrProjectNotFound {
			return nil, altalune.NewProjectNotFound(req.ProjectId)
		}
		return nil, altalune.NewInvalidPayloadError("invalid project_id")
	}

	// Prepare input for repository
	input := &DeactivateApiKeyInput{
		ProjectID: projectID,
		PublicID:  req.ApiKeyId,
	}

	// Deactivate API key
	result, err := s.apiKeyRepo.Deactivate(ctx, input)
	if err != nil {
		if err == ErrApiKeyNotFound {
			return nil, altalune.NewApiKeyNotFoundError(req.ApiKeyId)
		}
		s.log.Error("failed to deactivate api key",
			"error", err,
			"project_id", projectID,
			"api_key_id", req.ApiKeyId,
		)
		return nil, altalune.NewUnexpectedError("failed to deactivate api key: %w", err)
	}

	// Log successful deactivation for audit purposes
	s.log.Info("api key deactivated",
		"project_id", projectID,
		"api_key_id", result.PublicID,
		"name", result.Name,
	)

	// Convert to domain model
	apiKey := &ApiKey{
		ID:         result.PublicID,
		Name:       result.Name,
		Expiration: result.Expiration,
		Active:     result.Active,
		CreatedAt:  result.CreatedAt,
		UpdatedAt:  result.UpdatedAt,
	}

	return &altalunev1.DeactivateApiKeyResponse{
		ApiKey:  apiKey.ToApiKeyProto(),
		Message: "API key deactivated successfully",
	}, nil
}
