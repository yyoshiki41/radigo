package radigo

import "testing"

func TestVersion(t *testing.T) {
	if Version() == "" {
		t.Errorf("Version is empty: %s", version)
	}
}
