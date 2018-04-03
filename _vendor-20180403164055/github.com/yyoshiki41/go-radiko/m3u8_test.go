package radiko

import (
	"testing"
)

func TestGetChunklistFromM3U8(t *testing.T) {
	_, err := GetChunklistFromM3U8("")
	if err == nil {
		t.Error("Should detect an error.")
	}
}
