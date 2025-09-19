package altalune

import "embed"

//go:embed all:frontend/.output/public
var FrontendEmbeddedFiles embed.FS
