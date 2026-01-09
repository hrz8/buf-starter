package oauth_client

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"buf.build/go/protovalidate"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	project_domain "github.com/hrz8/altalune/internal/domain/project"
	"github.com/hrz8/altalune/internal/shared/query"
)

type Service struct {
	altalunev1.UnimplementedOAuthClientServiceServer
	validator       protovalidate.Validator
	log             altalune.Logger
	projectRepo     project_domain.Repositor
	oauthClientRepo Repositor
}

func NewService(v protovalidate.Validator, log altalune.Logger, projectRepo project_domain.Repositor, oauthClientRepo Repositor) *Service {
	return &Service{
		validator:       v,
		log:             log,
		projectRepo:     projectRepo,
		oauthClientRepo: oauthClientRepo,
	}
}

// CreateOAuthClient creates a new OAuth client with generated client_id and secret
func (s *Service) CreateOAuthClient(ctx context.Context, req *altalunev1.CreateOAuthClientRequest) (*altalunev1.CreateOAuthClientResponse, error) {
	// 1. Validate request with protovalidate
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// 2. Validate redirect URIs
	if len(req.RedirectUris) == 0 {
		return nil, altalune.NewInvalidPayloadError("at least one redirect URI required")
	}
	for _, uri := range req.RedirectUris {
		if !isValidRedirectURI(uri) {
			return nil, altalune.NewInvalidPayloadError(fmt.Sprintf("invalid redirect URI: %s", uri))
		}
	}

	// 3. Extract and validate project ID
	projectID, err := s.projectRepo.GetIDByPublicID(ctx, req.ProjectId)
	if err != nil {
		if err == project_domain.ErrProjectNotFound {
			return nil, altalune.NewProjectNotFound(req.ProjectId)
		}
		return nil, altalune.NewInvalidPayloadError("invalid project_id")
	}

	// 4. Create OAuth client with Argon2 hashed secret
	input := &CreateOAuthClientInput{
		ProjectID:       projectID,
		ProjectPublicID: req.ProjectId,
		Name:            strings.TrimSpace(req.Name),
		RedirectURIs:    req.RedirectUris,
		PKCERequired:    req.PkceRequired,
		AllowedScopes:   req.AllowedScopes,
	}

	result, err := s.oauthClientRepo.Create(ctx, input)
	if err != nil {
		if err == ErrOAuthClientAlreadyExists {
			return nil, altalune.NewOAuthClientAlreadyExistsError(req.Name)
		}
		s.log.Error("failed to create oauth client",
			"error", err,
			"project_id", projectID,
			"name", req.Name,
		)
		return nil, altalune.NewUnexpectedError("failed to create oauth client: %w", err)
	}

	// 5. Log successful creation
	s.log.Info("oauth client created",
		"project_id", projectID,
		"client_public_id", result.Client.ID,
		"client_id", result.Client.ClientID.String(),
		"name", result.Client.Name,
	)

	// 6. Return client with PLAINTEXT secret (ONLY time it's returned)
	return &altalunev1.CreateOAuthClientResponse{
		Client:       result.Client.ToOAuthClientProto(),
		ClientSecret: result.ClientSecret,
		Message:      "OAuth client created successfully",
	}, nil
}

// QueryOAuthClients returns a paginated list of OAuth clients
func (s *Service) QueryOAuthClients(ctx context.Context, req *altalunev1.QueryOAuthClientsRequest) (*altalunev1.QueryOAuthClientsResponse, error) {
	// 1. Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// 2. Validate that query is provided
	if req.Query == nil {
		return nil, altalune.NewInvalidPayloadError("query is required")
	}

	// 3. Extract and validate project ID
	projectID, err := s.projectRepo.GetIDByPublicID(ctx, req.ProjectId)
	if err != nil {
		if err == project_domain.ErrProjectNotFound {
			return nil, altalune.NewProjectNotFound(req.ProjectId)
		}
		return nil, altalune.NewInvalidPayloadError("invalid project_id")
	}

	// 4. Convert proto request to domain query params
	queryParams := query.DefaultQueryParams(req.Query)

	// 5. Query OAuth clients from repository
	result, err := s.oauthClientRepo.Query(ctx, projectID, queryParams)
	if err != nil {
		s.log.Error("failed to query oauth clients",
			"error", err,
			"project_id", projectID,
			"keyword", queryParams.Keyword,
		)
		return nil, altalune.NewUnexpectedError("failed to query oauth clients: %w", err)
	}

	// 6. Handle empty results
	if result == nil {
		return &altalunev1.QueryOAuthClientsResponse{
			Clients: []*altalunev1.OAuthClient{},
			Meta: &altalunev1.QueryMetaResponse{
				RowCount:  0,
				PageCount: 0,
				Filters:   make(map[string]*altalunev1.FilterValues),
			},
			Message: "No OAuth clients found",
		}, nil
	}

	// 7. Convert domain models to proto messages
	clients := make([]*altalunev1.OAuthClient, 0, len(result.Data))
	for _, client := range result.Data {
		clients = append(clients, client.ToOAuthClientProto())
	}

	// 8. Convert filters to proto format
	protoFilters := make(map[string]*altalunev1.FilterValues)
	for key, values := range result.Filters {
		protoFilters[key] = &altalunev1.FilterValues{Values: values}
	}

	return &altalunev1.QueryOAuthClientsResponse{
		Clients: clients,
		Meta: &altalunev1.QueryMetaResponse{
			RowCount:  result.TotalRows,
			PageCount: result.TotalPages,
			Filters:   protoFilters,
		},
		Message: fmt.Sprintf("Found %d OAuth clients", result.TotalRows),
	}, nil
}

// GetOAuthClient retrieves a single OAuth client by public ID
func (s *Service) GetOAuthClient(ctx context.Context, req *altalunev1.GetOAuthClientRequest) (*altalunev1.GetOAuthClientResponse, error) {
	// 1. Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// 2. Extract and validate project ID
	projectID, err := s.projectRepo.GetIDByPublicID(ctx, req.ProjectId)
	if err != nil {
		if err == project_domain.ErrProjectNotFound {
			return nil, altalune.NewProjectNotFound(req.ProjectId)
		}
		return nil, altalune.NewInvalidPayloadError("invalid project_id")
	}

	// 3. Get OAuth client from repository
	client, err := s.oauthClientRepo.GetByPublicID(ctx, projectID, req.Id)
	if err != nil {
		if err == ErrOAuthClientNotFound {
			return nil, altalune.NewOAuthClientNotFoundError(req.Id)
		}
		s.log.Error("failed to get oauth client",
			"error", err,
			"project_id", projectID,
			"client_public_id", req.Id,
		)
		return nil, altalune.NewUnexpectedError("failed to get oauth client: %w", err)
	}

	return &altalunev1.GetOAuthClientResponse{
		Client:  client.ToOAuthClientProto(),
		Message: "OAuth client retrieved successfully",
	}, nil
}

// UpdateOAuthClient updates an existing OAuth client
func (s *Service) UpdateOAuthClient(ctx context.Context, req *altalunev1.UpdateOAuthClientRequest) (*altalunev1.UpdateOAuthClientResponse, error) {
	// 1. Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// 2. Extract and validate project ID
	projectID, err := s.projectRepo.GetIDByPublicID(ctx, req.ProjectId)
	if err != nil {
		if err == project_domain.ErrProjectNotFound {
			return nil, altalune.NewProjectNotFound(req.ProjectId)
		}
		return nil, altalune.NewInvalidPayloadError("invalid project_id")
	}

	// 3. Get existing client to check if default
	existingClient, err := s.oauthClientRepo.GetByPublicID(ctx, projectID, req.Id)
	if err != nil {
		if err == ErrOAuthClientNotFound {
			return nil, altalune.NewOAuthClientNotFoundError(req.Id)
		}
		return nil, altalune.NewUnexpectedError("failed to get oauth client: %w", err)
	}

	// 4. Validate redirect URIs if provided
	if len(req.RedirectUris) > 0 {
		for _, uri := range req.RedirectUris {
			if !isValidRedirectURI(uri) {
				return nil, altalune.NewInvalidPayloadError(fmt.Sprintf("invalid redirect URI: %s", uri))
			}
		}
	}

	// 5. Protect default client from disabling PKCE
	if existingClient.IsDefault && req.PkceRequired != nil && !*req.PkceRequired {
		return nil, altalune.NewInvalidPayloadError("PKCE cannot be disabled for default client")
	}

	// 6. Build update input
	input := &UpdateOAuthClientInput{
		PublicID:      req.Id,
		ProjectID:     projectID,
		Name:          req.Name,
		PKCERequired:  req.PkceRequired,
		AllowedScopes: req.AllowedScopes,
	}

	if len(req.RedirectUris) > 0 {
		input.RedirectURIs = req.RedirectUris
	}

	// 7. Update OAuth client
	updatedClient, err := s.oauthClientRepo.Update(ctx, input)
	if err != nil {
		if err == ErrOAuthClientNotFound {
			return nil, altalune.NewOAuthClientNotFoundError(req.Id)
		}
		if err == ErrOAuthClientAlreadyExists {
			return nil, altalune.NewOAuthClientAlreadyExistsError(*req.Name)
		}
		s.log.Error("failed to update oauth client",
			"error", err,
			"project_id", projectID,
			"client_public_id", req.Id,
		)
		return nil, altalune.NewUnexpectedError("failed to update oauth client: %w", err)
	}

	// 8. Log successful update
	s.log.Info("oauth client updated",
		"project_id", projectID,
		"client_public_id", updatedClient.ID,
		"name", updatedClient.Name,
	)

	return &altalunev1.UpdateOAuthClientResponse{
		Client:  updatedClient.ToOAuthClientProto(),
		Message: "OAuth client updated successfully",
	}, nil
}

// DeleteOAuthClient deletes an OAuth client (with default client protection)
func (s *Service) DeleteOAuthClient(ctx context.Context, req *altalunev1.DeleteOAuthClientRequest) (*altalunev1.DeleteOAuthClientResponse, error) {
	// 1. Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// 2. Extract and validate project ID
	projectID, err := s.projectRepo.GetIDByPublicID(ctx, req.ProjectId)
	if err != nil {
		if err == project_domain.ErrProjectNotFound {
			return nil, altalune.NewProjectNotFound(req.ProjectId)
		}
		return nil, altalune.NewInvalidPayloadError("invalid project_id")
	}

	// 3. Get client to verify it exists and check if default
	client, err := s.oauthClientRepo.GetByPublicID(ctx, projectID, req.Id)
	if err != nil {
		if err == ErrOAuthClientNotFound {
			return nil, altalune.NewOAuthClientNotFoundError(req.Id)
		}
		return nil, altalune.NewUnexpectedError("failed to get oauth client: %w", err)
	}

	// 4. Protect default client from deletion
	if client.IsDefault {
		return nil, altalune.NewInvalidPayloadError("cannot delete default dashboard client")
	}

	// 5. Delete OAuth client
	err = s.oauthClientRepo.Delete(ctx, projectID, req.Id)
	if err != nil {
		if err == ErrOAuthClientNotFound {
			return nil, altalune.NewOAuthClientNotFoundError(req.Id)
		}
		if err == ErrDefaultClientCannotBeDeleted {
			return nil, altalune.NewInvalidPayloadError("cannot delete default dashboard client")
		}
		s.log.Error("failed to delete oauth client",
			"error", err,
			"project_id", projectID,
			"client_public_id", req.Id,
		)
		return nil, altalune.NewUnexpectedError("failed to delete oauth client: %w", err)
	}

	// 6. Log successful deletion
	s.log.Info("oauth client deleted",
		"project_id", projectID,
		"client_public_id", req.Id,
		"name", client.Name,
	)

	return &altalunev1.DeleteOAuthClientResponse{
		Message: "OAuth client deleted successfully",
	}, nil
}

// RevealOAuthClientSecret reveals the hashed client secret (with audit logging)
func (s *Service) RevealOAuthClientSecret(ctx context.Context, req *altalunev1.RevealOAuthClientSecretRequest) (*altalunev1.RevealOAuthClientSecretResponse, error) {
	// 1. Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// 2. Extract and validate project ID
	projectID, err := s.projectRepo.GetIDByPublicID(ctx, req.ProjectId)
	if err != nil {
		if err == project_domain.ErrProjectNotFound {
			return nil, altalune.NewProjectNotFound(req.ProjectId)
		}
		return nil, altalune.NewInvalidPayloadError("invalid project_id")
	}

	// 3. Reveal client secret (hashed) from repository
	hashedSecret, err := s.oauthClientRepo.RevealClientSecret(ctx, projectID, req.Id)
	if err != nil {
		if err == ErrOAuthClientNotFound {
			return nil, altalune.NewOAuthClientNotFoundError(req.Id)
		}
		s.log.Error("failed to reveal client secret",
			"error", err,
			"project_id", projectID,
			"client_public_id", req.Id,
		)
		return nil, altalune.NewUnexpectedError("failed to reveal client secret: %w", err)
	}

	// 4. Log audit event (CRITICAL for security)
	s.log.Warn("oauth_client_secret_revealed",
		"project_id", projectID,
		"client_public_id", req.Id,
		// TODO: Add user_id from context when authentication is implemented
	)

	// 5. Return hashed secret (Argon2id PHC string format)
	// NOTE: This is the HASHED secret, not plaintext
	// The plaintext secret was only shown during creation
	return &altalunev1.RevealOAuthClientSecretResponse{
		ClientSecret: hashedSecret,
		Message:      "Client secret revealed. This action has been logged for audit purposes.",
	}, nil
}

// isValidRedirectURI validates a redirect URI for OAuth 2.0 compliance
func isValidRedirectURI(uri string) bool {
	// Parse URI
	parsed, err := url.Parse(uri)
	if err != nil {
		return false
	}

	// Must be http or https scheme
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false
	}

	// No wildcards or query parameters allowed
	if strings.Contains(uri, "*") || strings.Contains(uri, "?") {
		return false
	}

	// Must have a host
	if parsed.Host == "" {
		return false
	}

	return true
}
