package radigo

import (
	"flag"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()

	code := m.Run()
	os.Exit(code)
}
