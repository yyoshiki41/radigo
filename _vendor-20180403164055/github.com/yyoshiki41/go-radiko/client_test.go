package radiko

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	_, err := New("")
	if err != nil {
		t.Fatalf("Failed to construct client: %s", err)
	}
}

func TestNew_EmptyHTTPClient(t *testing.T) {
	var c *http.Client

	SetHTTPClient(c)
	defer teardownHTTPClient()

	client, err := New("")
	if err == nil {
		t.Errorf(
			"Should detect that HTTPClient is nil.\nclient: %v", client)
	}
}

func TestNewRequest(t *testing.T) {
	client, err := New("")
	if err != nil {
		t.Fatalf("Failed to construct client: %s", err)
	}

	ctx := context.Background()
	_, err = client.newRequest(ctx, "GET", "", &Params{})
	if err != nil {
		t.Error(err)
	}
}

func TestNewRequest_WithAuthToken(t *testing.T) {
	const expected = "auth_token"

	client, err := New(expected)
	if err != nil {
		t.Errorf("Failed to construct client: %s", err)
	}

	req, err := client.newRequest(context.Background(), "GET", "", &Params{
		setAuthToken: true,
	})
	if err != nil {
		t.Error(err)
	}
	if actual := req.Header.Get(radikoAuthTokenHeader); actual != expected {
		t.Errorf("expected %s, but %s.", expected, actual)
	}
}

func TestNewRequest_WithContext(t *testing.T) {
	client, err := New("")
	if err != nil {
		t.Fatalf("Failed to construct client: %s", err)
	}

	timeout := 100 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err = client.newRequest(ctx, "GET", "", &Params{})
	if err != nil {
		t.Error(err)
	}

	select {
	case <-time.After(3 * time.Second):
		t.Fatalf("context: %v", ctx)
	case <-ctx.Done():
	}

	if ctx.Err() == nil {
		t.Errorf("Shoud detect the context deadline exceeded.\n%v", ctx)
	}
}

func TestNewRequest_WithEmptyContext(t *testing.T) {
	client, err := New("")
	if err != nil {
		t.Fatalf("Failed to construct client: %s", err)
	}

	var ctx context.Context
	_, err = client.newRequest(ctx, "GET", "", &Params{})
	if err == nil {
		t.Error("Should detect an empty context.")
	}
}

func TestClient_AreaID(t *testing.T) {
	client, err := New("")
	if err != nil {
		t.Errorf("Failed to construct client: %s", err)
	}

	areaID := client.AreaID()
	if areaID == "" {
		t.Error("httpClient.AreaID is empty.")
	}
	if !(strings.HasPrefix(areaID, "JP") || areaID == "OUT") {
		t.Errorf("Invalid area id.\nAreaID: %s", areaID)
	}
}

func TestClient_SetAreaID(t *testing.T) {
	client, err := New("")
	if err != nil {
		t.Errorf("Failed to construct client: %s", err)
	}

	const expected = "JP13"

	client.SetAreaID(expected)
	if actual := client.AreaID(); expected != actual {
		t.Errorf("expected %v, but %v.", expected, actual)
	}
}

func TestClient_SetJar(t *testing.T) {
	client, err := New("")
	if err != nil {
		t.Errorf("Failed to construct client: %s", err)
	}

	expected, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		t.Fatal(err)
	}
	client.SetJar(expected)
	defer teardownHTTPClient()

	if actual := client.Jar(); expected != actual {
		t.Errorf("expected %v, but %v.", expected, actual)
	}
}

func TestClient_SetAuthTokenHeader(t *testing.T) {
	client, err := New("")
	if err != nil {
		t.Errorf("Failed to construct client: %s", err)
	}

	const expected = "test_token"
	client.setAuthTokenHeader(expected)
	if actual := client.AuthToken(); expected != actual {
		t.Errorf("expected %s, but %s", expected, actual)
	}
}

func TestDo(t *testing.T) {
	client, err := New("")
	if err != nil {
		t.Errorf("Failed to construct client: %s", err)
	}

	ctx := context.Background()
	req, err := client.newRequest(ctx, "GET", "", &Params{})
	if err != nil {
		t.Error(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	const expected = 200
	if actual := resp.StatusCode; actual != expected {
		t.Errorf("expected %d, but StatusCode is %d.", expected, actual)
	}
}

func TestSetHTTPClient(t *testing.T) {
	const expected = 1 * time.Second

	SetHTTPClient(&http.Client{Timeout: expected})
	defer teardownHTTPClient()

	client, err := New("")
	if err != nil {
		t.Errorf("Failed to construct client: %s", err)
	}
	if client.httpClient.Timeout != expected {
		t.Errorf("expected %d, but %d", expected, client.httpClient.Timeout)
	}
}

func TestSetUserAgent(t *testing.T) {
	const expected = "test-user-agent"
	SetUserAgent(expected)
	if expected != userAgent {
		t.Errorf("expected %s, but %s", expected, userAgent)
	}
}

func TestAPIPath(t *testing.T) {
	const path = "test"
	var apiEndpoint string

	apiEndpoint = apiPath(apiV2, path)
	if !(strings.HasPrefix(apiEndpoint, apiV2+"/") && strings.HasSuffix(apiEndpoint, "/"+path)) {
		t.Errorf("invalid apiEndpoint: %s", apiEndpoint)
	}

	apiEndpoint = apiPath(apiV3, path)
	if !(strings.HasPrefix(apiEndpoint, apiV3+"/") && strings.HasSuffix(apiEndpoint, "/"+path)) {
		t.Errorf("invalid apiEndpoint: %s", apiEndpoint)
	}
}
