package radiko

import (
	"context"
	"testing"
	"time"
)

func TestTimeshiftPlaylistM3U8(t *testing.T) {
	if isOutsideJP() {
		t.Skip("Skipping test in limited mode.")
	}
}

func TestTimeshiftPlaylistM3U8_EmptyStationID(t *testing.T) {
	client, err := New("")
	if err != nil {
		t.Fatalf("Failed to construct client: %s", err)
	}

	_, err = client.TimeshiftPlaylistM3U8(context.Background(), "", time.Now())
	if err == nil {
		t.Error("Should detect an error.")
	}
}

func TestGetTimeshiftURL(t *testing.T) {
	stationID := "LFR"
	url := GetTimeshiftURL(stationID, time.Now())
	if len(url) == 0 {
		t.Error("A timeshift url is empty.")
	}
}
