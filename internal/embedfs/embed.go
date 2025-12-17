// Package embedfs embed sql migration and web build dist
package embedfs

import (
	"embed"
)

//go:embed web
var WebFile embed.FS

//go:embed sql
var SQLFile embed.FS
