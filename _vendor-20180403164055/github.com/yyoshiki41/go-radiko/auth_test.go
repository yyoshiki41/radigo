package radiko

import (
	"context"
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

func TestAuthorizeToken(t *testing.T) {
	if isOutsideJP() {
		t.Skip("Skipping test in limited mode.")
	}

	c, err := New("")
	if err != nil {
		t.Fatalf("Failed to construct client: %s", err)
	}

	ctx := context.Background()
	authToken, err := c.AuthorizeToken(ctx)
	if err != nil {
		t.Error(err)
	}
	if len(authToken) == 0 {
		t.Error("AuthToken is empty.")
	}
}

func TestAuth1Fms(t *testing.T) {
	c, err := New("")
	if err != nil {
		t.Fatalf("Failed to construct client: %s", err)
	}

	ctx := context.Background()
	authToken, length, offset, err := c.Auth1Fms(ctx)
	if err != nil {
		t.Error(err)
	}
	if len(authToken) == 0 || length < 0 || offset < 0 {
		t.Errorf("AuthToken: %s, Length: %d, Offset: %d", authToken, length, offset)
	}
}

func TestAuth2Fms(t *testing.T) {
	c, err := New("")
	if err != nil {
		t.Fatalf("Failed to construct client: %s", err)
	}

	ctx := context.Background()
	authToken, length, offset, err := c.Auth1Fms(ctx)
	if err != nil {
		t.Error(err)
	}

	pngPath := filepath.Join(testdataDir, "authkey.png")
	f, err := os.Open(pngPath)
	if err != nil {
		t.Error(err)
	}
	b := make([]byte, length)
	if _, err = f.ReadAt(b, offset); err != nil {
		t.Error(err)
	}
	partialKey := base64.StdEncoding.EncodeToString(b)

	_, err = c.Auth2Fms(ctx, authToken, partialKey)
	if err != nil {
		t.Error(err)
	}
}

func TestVerifyAuth2FmsResponse(t *testing.T) {
	cases := []struct {
		slc         []string
		expectedErr bool
	}{
		{
			slc: []string{`

	JP13`, "東京都", "tokyo", "Japan",
			},
			expectedErr: false,
		},
		{
			slc:         []string{},
			expectedErr: true,
		},
		{
			slc:         []string{"OUT"},
			expectedErr: true,
		},
	}
	for _, c := range cases {
		err := verifyAuth2FmsResponse(c.slc)
		if c.expectedErr {
			if err == nil {
				t.Error("Should detect an error.")
			}
			continue
		}
		if err != nil {
			t.Error(err)
		}
	}
}
