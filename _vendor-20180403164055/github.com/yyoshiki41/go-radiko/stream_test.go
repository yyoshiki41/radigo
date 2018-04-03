package radiko

import (
	"testing"
)

func TestGetStreamMultiURL(t *testing.T) {
	items, err := GetStreamMultiURL("LFR")
	if err != nil {
		t.Error(err)
	}

	const expected = 4
	if actual := len(items); expected != actual {
		t.Errorf("expected %d, but %d", expected, actual)
	}
}

func TestGetStreamMultiURL_EmptyStationID(t *testing.T) {
	_, err := GetStreamMultiURL("")
	if err == nil {
		t.Error("Should detect an error.")
	}
}

func TestGetLiveURL(t *testing.T) {
	stationID := "LFR"
	url := GetLiveURL(stationID)
	if len(url) == 0 {
		t.Error("A live url is empty.")
	}
}
