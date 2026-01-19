package chatbot

import "errors"

var (
	ErrChatbotConfigNotFound = errors.New("chatbot config not found")
	ErrInvalidModuleName     = errors.New("invalid module name")
	ErrInvalidModuleConfig   = errors.New("invalid module configuration")
)
