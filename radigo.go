package radigo

import (
	"time"

	"github.com/yyoshiki41/go-radiko"
)

const (
	tz             = "Asia/Tokyo"
	datetimeLayout = "20060102150405"
)

var (
	currentAreaID string
	location      *time.Location
)

func init() {
	var err error

	currentAreaID, err = radiko.AreaID()
	if err != nil {
		panic(err)
	}

	location, err = time.LoadLocation(tz)
	if err != nil {
		panic(err)
	}
}

const version = "v0.11.0"

// Version returns the app version.
func Version() string {
	return version
}
