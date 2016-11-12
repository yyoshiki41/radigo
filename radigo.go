package radigo

import (
	"path"
	"time"
)

const (
	version    = "0.0.1"
	radigoPath = "/tmp/radigo"

	tz             = "Asia/Tokyo"
	datetimeLayout = "20060102150405"
)

var (
	aacPath   = path.Join(radigoPath, "aac")
	cachePath = path.Join(radigoPath, ".cache")

	location *time.Location
)

func init() {
	var err error

	location, err = time.LoadLocation(tz)
	if err != nil {
		panic(err)
	}
}

// Version returns the app version.
func Version() string {
	return version
}
