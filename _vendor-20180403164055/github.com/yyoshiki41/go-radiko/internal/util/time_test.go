package util

import (
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	if location == nil {
		t.Fatal("location is nil.")
	}
}

func TestDate(t *testing.T) {
	s := Date(time.Now())
	if len(s) != len(dateLayout) {
		t.Errorf("date: %s", s)
	}
}

func TestDatetime(t *testing.T) {
	s := Datetime(time.Now())
	if len(s) != len(datetimeLayout) {
		t.Errorf("datetime: %s", s)
	}
}

func TestProgramsDate(t *testing.T) {
	y, m, d := time.Now().Date()

	date := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	pDate := ProgramsDate(date)
	if expected := Date(date); expected != pDate {
		t.Errorf("expected %s, but %s", expected, pDate)
	}

	date = date.Add(16 * time.Hour)
	pDate = ProgramsDate(date)
	if expected := Date(date.Add(-24 * time.Hour)); expected != pDate {
		t.Errorf("expected %s, but %s", expected, pDate)
	}
}
