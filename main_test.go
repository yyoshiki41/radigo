package radigo

import (
	"flag"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()

	code := m.Run()
	os.Exit(code)
}

func TestAreaID(t *testing.T) {
	if currentAreaID == "" {
		t.Errorf("currentAreaID is empty")
	}

	if !(strings.HasPrefix(currentAreaID, "JP") || currentAreaID == "OUT") {
		t.Errorf("Invalid currentAreaID: %s", currentAreaID)
	}
}

func TestLocation(t *testing.T) {
	if actual := location.String(); actual != tz {
		t.Errorf("location is expected %s, but %s", tz, actual)
	}
}
