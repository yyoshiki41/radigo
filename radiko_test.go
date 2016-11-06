package radigo

import (
	"os"
	"path"
	"testing"
)

const (
	// EnvTestDir is the environment variable that overrrides the working directory.
	EnvTestDir = "TEST_DIR"
)

func TestDownloadPlayer(t *testing.T) {
	dir := os.Getenv(EnvTestDir)
	if dir == "" {
		dir = radigoPath
	}
	myPlayerPath := path.Join(dir, "myplayer.swf")
	err := downloadPlayer(myPlayerPath)
	if err != nil {
		t.Errorf("Failed to download player.swf: %s", err)
	}
}
