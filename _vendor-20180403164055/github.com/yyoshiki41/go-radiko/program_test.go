package radiko

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yyoshiki41/go-radiko/internal/util"
)

func TestGetStations(t *testing.T) {
	if isOutsideJP() {
		t.Skip("Skipping test in limited mode.")
	}

	c, err := New("")
	if err != nil {
		t.Fatalf("Failed to construct client: %s", err)
	}

	c.SetAreaID(areaIDTokyo)
	stations, err := c.GetStations(context.Background(), time.Now())
	if err != nil {
		t.Error(err)
	}
	if len(stations) == 0 {
		t.Error("Stations is nil.")
	}
}

func TestGetNowPrograms(t *testing.T) {
	if isOutsideJP() {
		t.Skip("Skipping test in limited mode.")
	}

	c, err := New("")
	if err != nil {
		t.Fatalf("Failed to construct client: %s", err)
	}

	programs, err := c.GetNowPrograms(context.Background())
	if err != nil {
		t.Error(err)
	}
	if len(programs) == 0 {
		t.Error("Programs is nil.")
	}
}

func TestGetProgramByStartTime(t *testing.T) {
	if isOutsideJP() {
		t.Skip("Skipping test in limited mode.")
	}

	c, err := New("")
	if err != nil {
		t.Fatalf("Failed to construct client: %s", err)
	}

	// Tests in ANN
	stationID := "LFR"
	n := time.Now()
	if n.Weekday() == time.Sunday {
		// If it is Sunday, ANN will not be broadcasted.
		n = n.Add(-24 * time.Hour)
	}
	y, m, d := n.Date()
	// ANN starts at 01:00 AM on Monday to Saturday in JST.
	start := time.Date(y, m, d, 16, 0, 0, 0, time.UTC)
	end := time.Date(y, m, d, 18, 0, 0, 0, time.UTC)

	prog, err := c.GetProgramByStartTime(context.Background(), stationID, start)
	if err != nil {
		t.Error(err)
	}
	expected := util.Datetime(end)
	if expected != prog.To {
		t.Errorf("expected %s, but %s", expected, prog.To)
	}
}

func TestGetProgramByStartTimeEmptyStationID(t *testing.T) {
	c, err := New("")
	if err != nil {
		t.Fatalf("Failed to construct client: %s", err)
	}

	_, err = c.GetProgramByStartTime(context.Background(), "", time.Now())
	if err == nil {
		t.Error("Should detect an error.")
	}
}

func TestGetWeeklyPrograms(t *testing.T) {
	c, err := New("")
	if err != nil {
		t.Fatalf("Failed to construct client: %s", err)
	}

	programs, err := c.GetWeeklyPrograms(context.Background(), "LFR")
	if err != nil {
		t.Error(err)
	}
	if len(programs) == 0 {
		t.Error("Programs is nil.")
	}
}

func TestDecodeStationsData(t *testing.T) {
	file, err := os.Open(filepath.Join(testdataDir, "stations.xml"))
	if err != nil {
		t.Fatal(err)
	}

	var d stationsData
	if err = decodeStationsData(file, &d); err != nil {
		t.Error(err)
	}

	const expected = 2
	if s := d.stations(); expected != len(s) {
		t.Errorf("expected number of stations %d, but %d.", expected, len(s))
	}
}
