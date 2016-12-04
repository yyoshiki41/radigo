package radigo

import (
	"os"
	"path"
	"time"

	"github.com/yyoshiki41/go-radiko"
)

const (
	version = "0.1.0"

	tz             = "Asia/Tokyo"
	datetimeLayout = "20060102150405"
)

var (
	radigoPath = "/tmp/radigo"
	cachePath  = path.Join(radigoPath, ".cache")

	currentAreaID string
	location      *time.Location
)

func init() {
	var err error

	// If the environment variable RADIGO_HOME is set,
	// override working directory path.
	if e := os.Getenv("RADIGO_HOME"); e != "" {
		radigoPath = e
	}

	currentAreaID, err = radiko.AreaID()
	if err != nil {
		panic(err)
	}

	location, err = time.LoadLocation(tz)
	if err != nil {
		panic(err)
	}
}

// Version returns the app version.
func Version() string {
	return version
}
