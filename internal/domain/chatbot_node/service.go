package chatbot_node

import (
	"context"
	"regexp"

	"buf.build/go/protovalidate"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	project_domain "github.com/hrz8/altalune/internal/domain/project"
)

// nodeNamePattern validates lowercase_snake_case names
var nodeNamePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

type Service struct {
	altalunev1.UnimplementedChatbotNodeServiceServer
	validator   protovalidate.Validator
	log         altalune.Logger
	projectRepo project_domain.Repositor
	nodeRepo    Repositor
}

func NewService(v protovalidate.Validator, log altalune.Logger, projectRepo project_domain.Repositor, nodeRepo Repositor) *Service {
	return &Service{
		validator:   v,
		log:         log,
		projectRepo: projectRepo,
		nodeRepo:    nodeRepo,
	}
}

// ListNodes retrieves all nodes for a project (for sidebar display)
func (s *Service) ListNodes(ctx context.Context, req *altalunev1.ListNodesRequest) (*altalunev1.ListNodesResponse, error) {
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

	// Get all nodes for the project
	nodes, err := s.nodeRepo.ListByProjectID(ctx, projectID)
	if err != nil {
		s.log.Error("failed to list nodes",
			"error", err,
			"project_id", projectID,
		)
		return nil, altalune.NewUnexpectedError("failed to list nodes: %w", err)
	}

	return &altalunev1.ListNodesResponse{
		Nodes: MapNodesToProto(nodes),
	}, nil
}

// CreateNode creates a new chatbot node
func (s *Service) CreateNode(ctx context.Context, req *altalunev1.CreateNodeRequest) (*altalunev1.CreateNodeResponse, error) {
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

	// Validate node name format
	if !nodeNamePattern.MatchString(req.Name) {
		return nil, altalune.NewInvalidPayloadError("invalid node name format: must be lowercase_snake_case starting with a letter")
	}

	// Validate language
	if !IsValidLanguage(req.Lang) {
		return nil, altalune.NewInvalidPayloadError("invalid language: must be one of en-US, id-ID")
	}

	// Check for duplicate name+lang
	exists, err := s.nodeRepo.ExistsByNameLang(ctx, projectID, req.Name, req.Lang, nil)
	if err != nil {
		s.log.Error("failed to check name_lang exists",
			"error", err,
			"project_id", projectID,
			"name", req.Name,
			"lang", req.Lang,
		)
		return nil, altalune.NewUnexpectedError("failed to check duplicate: %w", err)
	}
	if exists {
		return nil, altalune.NewInvalidPayloadError("node with this name and language already exists")
	}

	// Create the node
	input := &CreateNodeInput{
		ProjectID: projectID,
		Name:      req.Name,
		Lang:      req.Lang,
		Tags:      req.Tags,
	}

	node, err := s.nodeRepo.Create(ctx, input)
	if err != nil {
		s.log.Error("failed to create node",
			"error", err,
			"project_id", projectID,
			"name", req.Name,
		)
		return nil, altalune.NewUnexpectedError("failed to create node: %w", err)
	}

	s.log.Info("node created",
		"project_id", projectID,
		"node_id", node.ID,
		"name", node.Name,
		"lang", node.Lang,
	)

	return &altalunev1.CreateNodeResponse{
		Node: MapNodeToProto(node),
	}, nil
}

// GetNode retrieves a single node with full data (triggers, messages)
func (s *Service) GetNode(ctx context.Context, req *altalunev1.GetNodeRequest) (*altalunev1.GetNodeResponse, error) {
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

	// Get the node
	node, err := s.nodeRepo.GetByID(ctx, projectID, req.NodeId)
	if err != nil {
		if err == ErrNodeNotFound {
			return nil, altalune.NewChatbotNodeNotFoundError(req.NodeId)
		}
		s.log.Error("failed to get node",
			"error", err,
			"project_id", projectID,
			"node_id", req.NodeId,
		)
		return nil, altalune.NewUnexpectedError("failed to get node: %w", err)
	}

	return &altalunev1.GetNodeResponse{
		Node: MapNodeToProto(node),
	}, nil
}

// UpdateNode updates an existing chatbot node
func (s *Service) UpdateNode(ctx context.Context, req *altalunev1.UpdateNodeRequest) (*altalunev1.UpdateNodeResponse, error) {
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

	// Get the existing node to validate it exists
	existingNode, err := s.nodeRepo.GetByID(ctx, projectID, req.NodeId)
	if err != nil {
		if err == ErrNodeNotFound {
			return nil, altalune.NewChatbotNodeNotFoundError(req.NodeId)
		}
		s.log.Error("failed to get node for update",
			"error", err,
			"project_id", projectID,
			"node_id", req.NodeId,
		)
		return nil, altalune.NewUnexpectedError("failed to get node: %w", err)
	}

	// Build update input
	input := &UpdateNodeInput{
		ProjectID: projectID,
		NodeID:    req.NodeId,
	}

	// Handle optional name update
	if req.Name != nil {
		name := *req.Name
		// Validate node name format
		if !nodeNamePattern.MatchString(name) {
			return nil, altalune.NewInvalidPayloadError("invalid node name format: must be lowercase_snake_case starting with a letter")
		}

		// Check for duplicate name+lang (excluding current node)
		exists, err := s.nodeRepo.ExistsByNameLang(ctx, projectID, name, existingNode.Lang, &req.NodeId)
		if err != nil {
			s.log.Error("failed to check name_lang exists",
				"error", err,
				"project_id", projectID,
				"name", name,
				"lang", existingNode.Lang,
			)
			return nil, altalune.NewUnexpectedError("failed to check duplicate: %w", err)
		}
		if exists {
			return nil, altalune.NewInvalidPayloadError("node with this name and language already exists")
		}

		input.Name = &name
	}

	// Handle tags - use empty slice if provided but empty
	if req.Tags != nil {
		input.Tags = req.Tags
	}

	// Handle optional enabled update
	if req.Enabled != nil {
		input.Enabled = req.Enabled
	}

	// Handle triggers update
	if req.Triggers != nil {
		triggers := MapProtoTriggersToModel(req.Triggers)

		// Validate triggers
		for _, t := range triggers {
			if !IsValidTriggerType(t.Type) {
				return nil, altalune.NewInvalidPayloadError("invalid trigger type: must be one of keyword, contains, regex, equals")
			}

			// Validate regex patterns
			if t.Type == "regex" {
				if _, err := regexp.Compile(t.Value); err != nil {
					return nil, altalune.NewInvalidPayloadError("invalid regex pattern: " + err.Error())
				}
			}
		}

		input.Triggers = triggers
	}

	// Handle messages update
	if req.Messages != nil {
		input.Messages = MapProtoMessagesToModel(req.Messages)
	}

	// Update the node
	node, err := s.nodeRepo.Update(ctx, input)
	if err != nil {
		if err == ErrNodeNotFound {
			return nil, altalune.NewChatbotNodeNotFoundError(req.NodeId)
		}
		s.log.Error("failed to update node",
			"error", err,
			"project_id", projectID,
			"node_id", req.NodeId,
		)
		return nil, altalune.NewUnexpectedError("failed to update node: %w", err)
	}

	s.log.Info("node updated",
		"project_id", projectID,
		"node_id", node.ID,
		"name", node.Name,
	)

	return &altalunev1.UpdateNodeResponse{
		Node:    MapNodeToProto(node),
		Message: "Node updated successfully",
	}, nil
}

// DeleteNode permanently removes a chatbot node
func (s *Service) DeleteNode(ctx context.Context, req *altalunev1.DeleteNodeRequest) (*altalunev1.DeleteNodeResponse, error) {
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

	// Delete the node
	err = s.nodeRepo.Delete(ctx, projectID, req.NodeId)
	if err != nil {
		if err == ErrNodeNotFound {
			return nil, altalune.NewChatbotNodeNotFoundError(req.NodeId)
		}
		s.log.Error("failed to delete node",
			"error", err,
			"project_id", projectID,
			"node_id", req.NodeId,
		)
		return nil, altalune.NewUnexpectedError("failed to delete node: %w", err)
	}

	s.log.Info("node deleted",
		"project_id", projectID,
		"node_id", req.NodeId,
	)

	return &altalunev1.DeleteNodeResponse{
		Message: "Node deleted successfully",
	}, nil
}
