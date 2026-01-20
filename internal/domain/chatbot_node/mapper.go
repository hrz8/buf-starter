package chatbot_node

import (
	"encoding/json"

	nodesv1 "github.com/hrz8/altalune/gen/chatbot/nodes/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToProto converts domain ChatbotNode to proto ChatbotNode
func (n *ChatbotNode) ToProto() *nodesv1.ChatbotNode {
	if n == nil {
		return nil
	}

	triggers := make([]*nodesv1.ChatbotNodeTrigger, len(n.Triggers))
	for i, t := range n.Triggers {
		triggers[i] = &nodesv1.ChatbotNodeTrigger{
			Type:  t.Type,
			Value: t.Value,
		}
	}

	messages := make([]*nodesv1.ChatbotNodeMessage, len(n.Messages))
	for i, m := range n.Messages {
		messages[i] = &nodesv1.ChatbotNodeMessage{
			Role:    m.Role,
			Content: m.Content,
		}
	}

	return &nodesv1.ChatbotNode{
		Id:        n.ID,
		Name:      n.Name,
		Lang:      n.Lang,
		Tags:      n.Tags,
		Enabled:   n.Enabled,
		Triggers:  triggers,
		Messages:  messages,
		CreatedAt: timestamppb.New(n.CreatedAt),
		UpdatedAt: timestamppb.New(n.UpdatedAt),
	}
}

// ToChatbotNode converts query result to domain model
func (r *ChatbotNodeQueryResult) ToChatbotNode() (*ChatbotNode, error) {
	var triggers []Trigger
	if len(r.Triggers) > 0 {
		if err := json.Unmarshal(r.Triggers, &triggers); err != nil {
			return nil, err
		}
	}

	var messages []Message
	if len(r.Messages) > 0 {
		if err := json.Unmarshal(r.Messages, &messages); err != nil {
			return nil, err
		}
	}

	// Ensure tags is never nil (empty array instead)
	tags := r.Tags
	if tags == nil {
		tags = []string{}
	}

	return &ChatbotNode{
		ID:        r.PublicID,
		ProjectID: r.ProjectID,
		Name:      r.Name,
		Lang:      r.Lang,
		Tags:      tags,
		Enabled:   r.Enabled,
		Triggers:  triggers,
		Messages:  messages,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}, nil
}

// MapNodeToProto is a helper function for mapping domain node to proto
func MapNodeToProto(node *ChatbotNode) *nodesv1.ChatbotNode {
	if node == nil {
		return nil
	}
	return node.ToProto()
}

// MapNodesToProto maps a slice of domain nodes to proto nodes
func MapNodesToProto(nodes []*ChatbotNode) []*nodesv1.ChatbotNode {
	result := make([]*nodesv1.ChatbotNode, len(nodes))
	for i, n := range nodes {
		result[i] = n.ToProto()
	}
	return result
}

// MapProtoTriggersToModel converts proto triggers to domain triggers
func MapProtoTriggersToModel(protoTriggers []*nodesv1.ChatbotNodeTrigger) []Trigger {
	if protoTriggers == nil {
		return nil
	}
	triggers := make([]Trigger, len(protoTriggers))
	for i, t := range protoTriggers {
		triggers[i] = Trigger{
			Type:  t.Type,
			Value: t.Value,
		}
	}
	return triggers
}

// MapProtoMessagesToModel converts proto messages to domain messages
func MapProtoMessagesToModel(protoMessages []*nodesv1.ChatbotNodeMessage) []Message {
	if protoMessages == nil {
		return nil
	}
	messages := make([]Message, len(protoMessages))
	for i, m := range protoMessages {
		messages[i] = Message{
			Role:    m.Role,
			Content: m.Content,
		}
	}
	return messages
}
