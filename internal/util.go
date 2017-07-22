package internal

import (
	"io/ioutil"
	"log"
	"os"
)

func CreateTempDir() (string, func()) {
	dir, err := ioutil.TempDir("", "radigo")
	if err != nil {
		log.Fatalf("Failed to create temp dir: %s", err)
	}

	return dir, func() { os.RemoveAll(dir) }
}
