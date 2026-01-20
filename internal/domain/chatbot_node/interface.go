package chatbot_node

import (
	"context"
)

// Repositor defines the interface for chatbot node repository operations
type Repositor interface {
	// ListByProjectID retrieves all nodes for a project, sorted alphabetically by name_lang
	ListByProjectID(ctx context.Context, projectID int64) ([]*ChatbotNode, error)

	// Create creates a new chatbot node
	Create(ctx context.Context, input *CreateNodeInput) (*ChatbotNode, error)

	// GetByID retrieves a single node by project ID and node public ID
	GetByID(ctx context.Context, projectID int64, nodeID string) (*ChatbotNode, error)

	// Update updates an existing chatbot node
	Update(ctx context.Context, input *UpdateNodeInput) (*ChatbotNode, error)

	// Delete permanently removes a chatbot node
	Delete(ctx context.Context, projectID int64, nodeID string) error

	// ExistsByNameLang checks if a node with the given name and lang already exists
	ExistsByNameLang(ctx context.Context, projectID int64, name, lang string, excludeNodeID *string) (bool, error)
}
