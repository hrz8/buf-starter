package chatbot_node

import "errors"

var (
	// ErrNodeNotFound is returned when a chatbot node is not found
	ErrNodeNotFound = errors.New("chatbot node not found")

	// ErrInvalidNodeName is returned when the node name format is invalid
	ErrInvalidNodeName = errors.New("invalid node name format")

	// ErrInvalidLanguage is returned when the language code is invalid
	ErrInvalidLanguage = errors.New("invalid language code")

	// ErrDuplicateNameLang is returned when a node with the same name+lang already exists
	ErrDuplicateNameLang = errors.New("node with this name and language already exists")

	// ErrAtLeastOneTriggerRequired is returned when no triggers are provided
	ErrAtLeastOneTriggerRequired = errors.New("at least one trigger is required")

	// ErrAtLeastOneMessageRequired is returned when no messages are provided
	ErrAtLeastOneMessageRequired = errors.New("at least one message is required")

	// ErrInvalidTriggerType is returned when a trigger type is invalid
	ErrInvalidTriggerType = errors.New("invalid trigger type")

	// ErrInvalidRegexPattern is returned when a regex pattern is invalid
	ErrInvalidRegexPattern = errors.New("invalid regex pattern")
)
