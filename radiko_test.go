package radigo

import (
	"os"
	"path"
	"testing"
)

const (
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
		t.Error("Failed to download player.swf: ", err)
	}
}
