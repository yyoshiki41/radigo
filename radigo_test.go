package radigo

import (
	"strings"
	"testing"
)

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

func TestVersion(t *testing.T) {
	if Version() == "" {
		t.Errorf("Version is empty: %s", version)
	}
}
