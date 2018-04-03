package radiko

import (
	"context"
	"testing"
)

func TestLogin_StatusCode400(t *testing.T) {
	// delete http.Client.Jar
	defer teardownHTTPClient()

	client, err := New("")
	if err != nil {
		t.Errorf("Failed to construct client: %s", err)
	}

	resp, err := client.Login(context.Background(),
		"test_mail", "test_pass")

	if err != nil {
		t.Error(err)
	}
	expected := 400
	if actual := resp.StatusCode(); expected != actual {
		t.Errorf("expected %d, but %d.", expected, actual)
	}
}
