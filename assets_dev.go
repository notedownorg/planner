//go:build dev

package main

import "embed"

// Development build: use empty embed to avoid requiring frontend/dist during development
// This allows Go linting and development without needing to build the frontend first
var assets embed.FS
