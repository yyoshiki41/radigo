package radigo

import "testing"

func TestVersion(t *testing.T) {
	if Version() == "" {
		t.Error("Version is empty: %s", version)
	}
}
