package radigo

import "path"

const (
	version    = "0.0.1"
	radigoPath = "/tmp/radigo"
)

var (
	aacPath   = path.Join(radigoPath, "aac")
	cachePath = path.Join(radigoPath, ".cache")
)

// Version returns the app version.
func Version() string {
	return version
}
