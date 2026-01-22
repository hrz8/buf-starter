package email

import (
	"embed"
)

// Templates embeds the email templates for use at runtime.
// This ensures templates are available regardless of working directory.
//
//go:embed templates/*.html templates/*.txt
var Templates embed.FS
