package ui

import "embed"

// this looks like a comment but it is a comment directive. instructs Go to store files from ui/html and ui/static in the embed.FS

//go:embed "html" "static"
var Files embed.FS
