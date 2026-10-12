// Package assets embeds the generator's template files,
// compiled into the binary via go:embed.
package assets

import "embed"

// FS is the embedded filesystem containing all project
// templates. Paths are relative to the templates/ root.
//
//go:embed all:templates
var FS embed.FS
