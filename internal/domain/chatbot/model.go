package chatbot

import (
	"encoding/json"
	"time"

	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ChatbotConfig represents a chatbot configuration in the domain
type ChatbotConfig struct {
	ID            string
	ProjectID     int64
	ModulesConfig map[string]interface{}
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// ToProto converts domain ChatbotConfig to proto ChatbotConfig
func (c *ChatbotConfig) ToProto() (*altalunev1.ChatbotConfig, error) {
	modulesStruct, err := structpb.NewStruct(c.ModulesConfig)
	if err != nil {
		return nil, err
	}

	return &altalunev1.ChatbotConfig{
		Id:            c.ID,
		ModulesConfig: modulesStruct,
		CreatedAt:     timestamppb.New(c.CreatedAt),
		UpdatedAt:     timestamppb.New(c.UpdatedAt),
	}, nil
}

// ChatbotConfigQueryResult represents raw database query result
type ChatbotConfigQueryResult struct {
	ID            int64
	PublicID      string
	ProjectID     int64
	ModulesConfig []byte // Raw JSONB from database
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// ToChatbotConfig converts query result to domain model
func (r *ChatbotConfigQueryResult) ToChatbotConfig() (*ChatbotConfig, error) {
	var modulesConfig map[string]interface{}
	if err := json.Unmarshal(r.ModulesConfig, &modulesConfig); err != nil {
		return nil, err
	}

	return &ChatbotConfig{
		ID:            r.PublicID,
		ProjectID:     r.ProjectID,
		ModulesConfig: modulesConfig,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}, nil
}

// UpdateModuleConfigInput represents input for updating a module config
type UpdateModuleConfigInput struct {
	ProjectID  int64
	ModuleName string
	Config     map[string]interface{}
}

// ValidModuleNames defines the allowed module names
var ValidModuleNames = map[string]bool{
	"llm":       true,
	"mcpServer": true,
	"widget":    true,
	"prompt":    true,
}

// IsValidModuleName checks if a module name is valid
func IsValidModuleName(name string) bool {
	return ValidModuleNames[name]
}
