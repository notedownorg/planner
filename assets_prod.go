//go:build !dev

package main

import "embed"

// Production build: embed the built frontend assets from frontend/dist
// This requires the frontend to be built first with `npm run build`
//
//go:embed all:frontend/dist
var assets embed.FS
