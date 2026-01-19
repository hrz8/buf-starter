package chatbot

import (
	"context"

	"buf.build/go/protovalidate"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	project_domain "github.com/hrz8/altalune/internal/domain/project"
)

type Service struct {
	altalunev1.UnimplementedChatbotServiceServer
	validator   protovalidate.Validator
	log         altalune.Logger
	projectRepo project_domain.Repositor
	chatbotRepo Repositor
}

func NewService(v protovalidate.Validator, log altalune.Logger, projectRepo project_domain.Repositor, chatbotRepo Repositor) *Service {
	return &Service{
		validator:   v,
		log:         log,
		projectRepo: projectRepo,
		chatbotRepo: chatbotRepo,
	}
}

func (s *Service) GetChatbotConfig(ctx context.Context, req *altalunev1.GetChatbotConfigRequest) (*altalunev1.GetChatbotConfigResponse, error) {
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

	// Get chatbot config (lazy initialization if not exists)
	config, err := s.chatbotRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		s.log.Error("failed to get chatbot config",
			"error", err,
			"project_id", projectID,
		)
		return nil, altalune.NewUnexpectedError("failed to get chatbot config: %w", err)
	}

	// Convert to proto
	protoConfig, err := mapChatbotConfigToProto(config)
	if err != nil {
		s.log.Error("failed to convert chatbot config to proto",
			"error", err,
			"project_id", projectID,
		)
		return nil, altalune.NewUnexpectedError("failed to convert chatbot config: %w", err)
	}

	return &altalunev1.GetChatbotConfigResponse{
		ChatbotConfig: protoConfig,
	}, nil
}

func (s *Service) UpdateModuleConfig(ctx context.Context, req *altalunev1.UpdateModuleConfigRequest) (*altalunev1.UpdateModuleConfigResponse, error) {
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

	// Validate module name
	if !IsValidModuleName(req.ModuleName) {
		return nil, altalune.NewInvalidPayloadError("invalid module name: must be one of llm, mcpServer, widget, prompt")
	}

	// Convert proto Struct to map
	moduleConfig := req.Config.AsMap()
	if moduleConfig == nil {
		moduleConfig = make(map[string]interface{})
	}

	// Prepare input for repository
	input := &UpdateModuleConfigInput{
		ProjectID:  projectID,
		ModuleName: req.ModuleName,
		Config:     moduleConfig,
	}

	// Update module config
	config, err := s.chatbotRepo.UpdateModuleConfig(ctx, input)
	if err != nil {
		if err == ErrInvalidModuleName {
			return nil, altalune.NewInvalidPayloadError("invalid module name")
		}
		if err == ErrInvalidModuleConfig {
			return nil, altalune.NewInvalidPayloadError("invalid module configuration")
		}
		s.log.Error("failed to update module config",
			"error", err,
			"project_id", projectID,
			"module_name", req.ModuleName,
		)
		return nil, altalune.NewUnexpectedError("failed to update module config: %w", err)
	}

	// Log successful update
	s.log.Info("module config updated",
		"project_id", projectID,
		"module_name", req.ModuleName,
	)

	// Convert to proto
	protoConfig, err := mapChatbotConfigToProto(config)
	if err != nil {
		s.log.Error("failed to convert chatbot config to proto",
			"error", err,
			"project_id", projectID,
		)
		return nil, altalune.NewUnexpectedError("failed to convert chatbot config: %w", err)
	}

	return &altalunev1.UpdateModuleConfigResponse{
		ChatbotConfig: protoConfig,
		Message:       "Module configuration updated successfully",
	}, nil
}
