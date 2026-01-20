package chatbot_node

import (
	"time"
)

// ChatbotNode represents a chatbot response node in the domain.
// Nodes contain triggers (conditions to activate) and messages (responses to send).
type ChatbotNode struct {
	ID        string
	ProjectID int64
	Name      string
	Lang      string
	Tags      []string
	Enabled   bool
	Triggers  []Trigger
	Messages  []Message
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Trigger defines a condition that activates a node.
type Trigger struct {
	Type  string // keyword, contains, regex, equals
	Value string
}

// Message defines a response message in OpenAI-compatible format.
type Message struct {
	Role    string // assistant
	Content string
}

// ChatbotNodeQueryResult represents raw database query result
type ChatbotNodeQueryResult struct {
	ID        int64
	PublicID  string
	ProjectID int64
	Name      string
	Lang      string
	Tags      []string
	Enabled   bool
	Triggers  []byte // Raw JSONB from database
	Messages  []byte // Raw JSONB from database
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateNodeInput represents input for creating a node
type CreateNodeInput struct {
	ProjectID int64
	Name      string
	Lang      string
	Tags      []string
}

// UpdateNodeInput represents input for updating a node
type UpdateNodeInput struct {
	ProjectID int64
	NodeID    string
	Name      *string
	Tags      []string
	Enabled   *bool
	Triggers  []Trigger
	Messages  []Message
}

// ValidTriggerTypes defines the allowed trigger types
var ValidTriggerTypes = map[string]bool{
	"keyword":  true,
	"contains": true,
	"regex":    true,
	"equals":   true,
}

// ValidLanguages defines the allowed language codes
var ValidLanguages = map[string]bool{
	"en-US": true,
	"id-ID": true,
}

// IsValidTriggerType checks if a trigger type is valid
func IsValidTriggerType(t string) bool {
	return ValidTriggerTypes[t]
}

// IsValidLanguage checks if a language code is valid
func IsValidLanguage(lang string) bool {
	return ValidLanguages[lang]
}
