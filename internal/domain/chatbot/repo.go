package chatbot

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hrz8/altalune/internal/postgres"
	"github.com/hrz8/altalune/internal/shared/nanoid"
)

// defaultChatbotModulesConfig is the default JSONB configuration for new chatbot configs
// This is the single source of truth for default module configuration
const defaultChatbotModulesConfig = `{
	"llm": {
		"enabled": false,
		"model": "",
		"temperature": 0.7,
		"maxToolCalls": 5
	},
	"mcpServer": {
		"enabled": false,
		"urls": [],
		"structuredOutputs": []
	},
	"widget": {
		"enabled": false,
		"cors": {
			"allowedOrigins": [],
			"allowedHeaders": ["Content-Type"],
			"credentials": false
		}
	},
	"prompt": {
		"enabled": true,
		"systemPrompt": "You are a helpful assistant."
	}
}`

type Repo struct {
	db postgres.DB
}

func NewRepo(db postgres.DB) *Repo {
	return &Repo{db: db}
}

// GetByProjectID retrieves chatbot config by project ID, creating it if not exists (lazy initialization)
func (r *Repo) GetByProjectID(ctx context.Context, projectID int64) (*ChatbotConfig, error) {
	// First try to get existing config
	config, err := r.getExisting(ctx, projectID)
	if err == nil {
		return config, nil
	}
	if !errors.Is(err, ErrChatbotConfigNotFound) {
		return nil, err
	}

	// Config doesn't exist - create default config
	// Note: Partition is created by migration (for existing projects) or project/repo.go (for new projects)
	return r.createDefault(ctx, projectID)
}

// getExisting retrieves an existing chatbot config from the database
func (r *Repo) getExisting(ctx context.Context, projectID int64) (*ChatbotConfig, error) {
	query := `
		SELECT
			id,
			public_id,
			project_id,
			modules_config,
			created_at,
			updated_at
		FROM altalune_chatbot_configs
		WHERE project_id = $1
	`

	var result ChatbotConfigQueryResult
	err := r.db.QueryRowContext(ctx, query, projectID).Scan(
		&result.ID,
		&result.PublicID,
		&result.ProjectID,
		&result.ModulesConfig,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrChatbotConfigNotFound
		}
		return nil, fmt.Errorf("get chatbot config: %w", err)
	}

	return result.ToChatbotConfig()
}

// createDefault creates a default chatbot config for a project
func (r *Repo) createDefault(ctx context.Context, projectID int64) (*ChatbotConfig, error) {
	// Check if config already exists (idempotent)
	existing, err := r.getExisting(ctx, projectID)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, ErrChatbotConfigNotFound) {
		return nil, err
	}

	// Generate public_id for the new config
	publicID, err := nanoid.GeneratePublicID()
	if err != nil {
		return nil, fmt.Errorf("generate public_id: %w", err)
	}

	// Insert default config
	insertQuery := `
		INSERT INTO altalune_chatbot_configs (
			public_id, project_id, modules_config, created_at, updated_at
		) VALUES ($1, $2, $3::jsonb, $4, $5)
		RETURNING id, public_id, project_id, modules_config, created_at, updated_at
	`

	now := time.Now()
	var result ChatbotConfigQueryResult
	err = r.db.QueryRowContext(
		ctx,
		insertQuery,
		publicID,
		projectID,
		defaultChatbotModulesConfig,
		now,
		now,
	).Scan(
		&result.ID,
		&result.PublicID,
		&result.ProjectID,
		&result.ModulesConfig,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("create default chatbot config: %w", err)
	}

	fmt.Printf("Info: Created default chatbot config for project %d\n", projectID)
	return result.ToChatbotConfig()
}

// UpdateModuleConfig updates a specific module's config within the chatbot config
func (r *Repo) UpdateModuleConfig(ctx context.Context, input *UpdateModuleConfigInput) (*ChatbotConfig, error) {
	// Validate module name
	if !IsValidModuleName(input.ModuleName) {
		return nil, ErrInvalidModuleName
	}

	// First ensure config exists (lazy init if needed)
	_, err := r.GetByProjectID(ctx, input.ProjectID)
	if err != nil {
		return nil, err
	}

	// Convert config to JSON
	configJSON, err := json.Marshal(input.Config)
	if err != nil {
		return nil, fmt.Errorf("marshal module config: %w", err)
	}

	// Update the specific module using jsonb_set
	updateQuery := `
		UPDATE altalune_chatbot_configs
		SET
			modules_config = jsonb_set(modules_config, $1::text[], $2::jsonb),
			updated_at = $3
		WHERE project_id = $4
		RETURNING id, public_id, project_id, modules_config, created_at, updated_at
	`

	// PostgreSQL path format: {moduleName}
	jsonPath := fmt.Sprintf("{%s}", input.ModuleName)
	now := time.Now()

	var result ChatbotConfigQueryResult
	err = r.db.QueryRowContext(
		ctx,
		updateQuery,
		jsonPath,
		string(configJSON),
		now,
		input.ProjectID,
	).Scan(
		&result.ID,
		&result.PublicID,
		&result.ProjectID,
		&result.ModulesConfig,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrChatbotConfigNotFound
		}
		return nil, fmt.Errorf("update module config: %w", err)
	}

	return result.ToChatbotConfig()
}
