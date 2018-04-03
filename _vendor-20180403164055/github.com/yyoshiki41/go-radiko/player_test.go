package radiko

import (
	"path/filepath"
	"testing"
)

func TestDownloadPlayer(t *testing.T) {
	dir, removeDir := createTestTempDir(t)
	defer removeDir() // clean up

	playerPath := filepath.Join(dir, "myplayer.swf")
	err := DownloadPlayer(playerPath)
	if err != nil {
		t.Errorf("Failed to download player.swf: %s", err)
	}
}

func TestDownloadBinary(t *testing.T) {
	_, err := downloadBinary()
	if err != nil {
		t.Errorf("Failed to download binary: %s", err)
	}
}
