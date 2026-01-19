package chatbot

import (
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
)

// mapChatbotConfigToProto converts domain ChatbotConfig to proto ChatbotConfig
func mapChatbotConfigToProto(config *ChatbotConfig) (*altalunev1.ChatbotConfig, error) {
	if config == nil {
		return nil, nil
	}
	return config.ToProto()
}
