package chatbot_node

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hrz8/altalune/internal/postgres"
	"github.com/hrz8/altalune/internal/shared/nanoid"
	"github.com/lib/pq"
)

type Repo struct {
	db postgres.DB
}

func NewRepo(db postgres.DB) *Repo {
	return &Repo{db: db}
}

// ListByProjectID retrieves all nodes for a project, sorted alphabetically by name_lang
func (r *Repo) ListByProjectID(ctx context.Context, projectID int64) ([]*ChatbotNode, error) {
	query := `
		SELECT
			id,
			public_id,
			project_id,
			name,
			lang,
			tags,
			enabled,
			triggers,
			messages,
			created_at,
			updated_at
		FROM altalune_chatbot_nodes
		WHERE project_id = $1
		ORDER BY name || '_' || lang ASC
	`

	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("list nodes: %w", err)
	}
	defer rows.Close()

	var nodes []*ChatbotNode
	for rows.Next() {
		var result ChatbotNodeQueryResult
		err := rows.Scan(
			&result.ID,
			&result.PublicID,
			&result.ProjectID,
			&result.Name,
			&result.Lang,
			pq.Array(&result.Tags),
			&result.Enabled,
			&result.Triggers,
			&result.Messages,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan node row: %w", err)
		}

		node, err := result.ToChatbotNode()
		if err != nil {
			return nil, fmt.Errorf("convert query result: %w", err)
		}
		nodes = append(nodes, node)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate nodes: %w", err)
	}

	return nodes, nil
}

// Create creates a new chatbot node with default empty triggers and messages
// Note: Partition is created by migration (for existing projects) or project/repo.go (for new projects)
func (r *Repo) Create(ctx context.Context, input *CreateNodeInput) (*ChatbotNode, error) {
	// Generate public_id for the new node
	publicID, err := nanoid.GeneratePublicID()
	if err != nil {
		return nil, fmt.Errorf("generate public_id: %w", err)
	}

	// Default empty triggers and messages
	defaultTriggers := "[]"
	defaultMessages := "[]"

	// Ensure tags is not nil
	tags := input.Tags
	if tags == nil {
		tags = []string{}
	}

	insertQuery := `
		INSERT INTO altalune_chatbot_nodes (
			public_id, project_id, name, lang, tags, enabled, triggers, messages, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, true, $6::jsonb, $7::jsonb, $8, $9)
		RETURNING id, public_id, project_id, name, lang, tags, enabled, triggers, messages, created_at, updated_at
	`

	now := time.Now()
	var result ChatbotNodeQueryResult
	err = r.db.QueryRowContext(
		ctx,
		insertQuery,
		publicID,
		input.ProjectID,
		input.Name,
		input.Lang,
		pq.Array(tags),
		defaultTriggers,
		defaultMessages,
		now,
		now,
	).Scan(
		&result.ID,
		&result.PublicID,
		&result.ProjectID,
		&result.Name,
		&result.Lang,
		pq.Array(&result.Tags),
		&result.Enabled,
		&result.Triggers,
		&result.Messages,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("create node: %w", err)
	}

	return result.ToChatbotNode()
}

// GetByID retrieves a single node by project ID and node public ID
func (r *Repo) GetByID(ctx context.Context, projectID int64, nodeID string) (*ChatbotNode, error) {
	query := `
		SELECT
			id,
			public_id,
			project_id,
			name,
			lang,
			tags,
			enabled,
			triggers,
			messages,
			created_at,
			updated_at
		FROM altalune_chatbot_nodes
		WHERE project_id = $1 AND public_id = $2
	`

	var result ChatbotNodeQueryResult
	err := r.db.QueryRowContext(ctx, query, projectID, nodeID).Scan(
		&result.ID,
		&result.PublicID,
		&result.ProjectID,
		&result.Name,
		&result.Lang,
		pq.Array(&result.Tags),
		&result.Enabled,
		&result.Triggers,
		&result.Messages,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNodeNotFound
		}
		return nil, fmt.Errorf("get node: %w", err)
	}

	return result.ToChatbotNode()
}

// Update updates an existing chatbot node
func (r *Repo) Update(ctx context.Context, input *UpdateNodeInput) (*ChatbotNode, error) {
	// First verify the node exists
	_, err := r.GetByID(ctx, input.ProjectID, input.NodeID)
	if err != nil {
		return nil, err
	}

	// Build dynamic update query
	setClauses := []string{"updated_at = $1"}
	args := []interface{}{time.Now()}
	argIndex := 2

	if input.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *input.Name)
		argIndex++
	}

	if input.Tags != nil {
		setClauses = append(setClauses, fmt.Sprintf("tags = $%d", argIndex))
		args = append(args, pq.Array(input.Tags))
		argIndex++
	}

	if input.Enabled != nil {
		setClauses = append(setClauses, fmt.Sprintf("enabled = $%d", argIndex))
		args = append(args, *input.Enabled)
		argIndex++
	}

	if input.Triggers != nil {
		triggersJSON, err := json.Marshal(input.Triggers)
		if err != nil {
			return nil, fmt.Errorf("marshal triggers: %w", err)
		}
		setClauses = append(setClauses, fmt.Sprintf("triggers = $%d::jsonb", argIndex))
		args = append(args, string(triggersJSON))
		argIndex++
	}

	if input.Messages != nil {
		messagesJSON, err := json.Marshal(input.Messages)
		if err != nil {
			return nil, fmt.Errorf("marshal messages: %w", err)
		}
		setClauses = append(setClauses, fmt.Sprintf("messages = $%d::jsonb", argIndex))
		args = append(args, string(messagesJSON))
		argIndex++
	}

	// Add WHERE clause parameters
	args = append(args, input.ProjectID, input.NodeID)

	query := fmt.Sprintf(`
		UPDATE altalune_chatbot_nodes
		SET %s
		WHERE project_id = $%d AND public_id = $%d
		RETURNING id, public_id, project_id, name, lang, tags, enabled, triggers, messages, created_at, updated_at
	`, joinStrings(setClauses, ", "), argIndex, argIndex+1)

	var result ChatbotNodeQueryResult
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&result.ID,
		&result.PublicID,
		&result.ProjectID,
		&result.Name,
		&result.Lang,
		pq.Array(&result.Tags),
		&result.Enabled,
		&result.Triggers,
		&result.Messages,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNodeNotFound
		}
		return nil, fmt.Errorf("update node: %w", err)
	}

	return result.ToChatbotNode()
}

// Delete permanently removes a chatbot node
func (r *Repo) Delete(ctx context.Context, projectID int64, nodeID string) error {
	query := `
		DELETE FROM altalune_chatbot_nodes
		WHERE project_id = $1 AND public_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, projectID, nodeID)
	if err != nil {
		return fmt.Errorf("delete node: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNodeNotFound
	}

	return nil
}

// ExistsByNameLang checks if a node with the given name and lang already exists
func (r *Repo) ExistsByNameLang(ctx context.Context, projectID int64, name, lang string, excludeNodeID *string) (bool, error) {
	var query string
	var args []interface{}

	if excludeNodeID != nil {
		query = `
			SELECT EXISTS(
				SELECT 1 FROM altalune_chatbot_nodes
				WHERE project_id = $1 AND name = $2 AND lang = $3 AND public_id != $4
			)
		`
		args = []interface{}{projectID, name, lang, *excludeNodeID}
	} else {
		query = `
			SELECT EXISTS(
				SELECT 1 FROM altalune_chatbot_nodes
				WHERE project_id = $1 AND name = $2 AND lang = $3
			)
		`
		args = []interface{}{projectID, name, lang}
	}

	var exists bool
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check name_lang exists: %w", err)
	}

	return exists, nil
}

// joinStrings joins strings with a separator
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
