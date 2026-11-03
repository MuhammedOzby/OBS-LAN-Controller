package web

import "embed"

// Assets contains the templates and static files bundled with the application.
//
//go:embed templates/*.tmpl static
var Assets embed.FS
