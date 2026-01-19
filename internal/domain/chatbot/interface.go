package chatbot

import (
	"context"
)

// Repositor defines the interface for chatbot repository operations
type Repositor interface {
	// GetByProjectID retrieves chatbot config by project ID, creating it if not exists (lazy initialization)
	GetByProjectID(ctx context.Context, projectID int64) (*ChatbotConfig, error)

	// UpdateModuleConfig updates a specific module's config within the chatbot config
	UpdateModuleConfig(ctx context.Context, input *UpdateModuleConfigInput) (*ChatbotConfig, error)
}
