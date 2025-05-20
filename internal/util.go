package internal

import (
	"log"
	"os"
)

func CreateTempDir() (string, func()) {
	dir, err := os.MkdirTemp("", "radigo")
	if err != nil {
		log.Fatalf("Failed to create temp dir: %s", err)
	}

	return dir, func() { os.RemoveAll(dir) }
}
