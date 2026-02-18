// Package assets provides embedded icon assets for diagram nodes.
//
// Icons are embedded at compile time using Go's embed package.
// The ReadFile function provides the same interface as the previous
// fileb0x-based implementation for backward compatibility.
package assets

import "embed"

//go:embed assets
var embeddedAssets embed.FS

// ReadFile reads an embedded asset file by path.
// Paths should be in the format "assets/<provider>/<category>/<icon>.png",
// e.g. "assets/generic/network/router.png".
func ReadFile(path string) ([]byte, error) {
	return embeddedAssets.ReadFile(path)
}
