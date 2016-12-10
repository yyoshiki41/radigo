package radigo

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/mitchellh/cli"
	"github.com/yyoshiki41/go-radiko"
)

type browseCommand struct {
	ui cli.Ui
}

func (c *browseCommand) Run(args []string) int {
	var stationID, start string

	f := flag.NewFlagSet("browse", flag.ContinueOnError)
	f.StringVar(&stationID, "id", "", "id")
	f.StringVar(&start, "start", "", "start")
	f.StringVar(&start, "s", "", "start")
	f.Usage = func() { c.ui.Error(c.Help()) }
	if err := f.Parse(args); err != nil {
		return 1
	}

	if stationID == "" {
		c.ui.Error("StationID is empty.")
		return 1
	}
	startTime, err := time.ParseInLocation(datetimeLayout, start, location)
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Invalid start time format: %s", start))
		return 1
	}

	url := radiko.GetTimeshiftURL(stationID, startTime)
	cmd := exec.Command("open", url)
	if err := cmd.Run(); err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to open browser: %s", err))
	}

	return 0
}

func (c *browseCommand) Synopsis() string {
	return "Browse radiko.jp"
}

func (c *browseCommand) Help() string {
	return strings.TrimSpace(`
Usage: radigo browse [options]
  Browse radiko.jp
Options:
  -id=name                 Station id
  -start,s=201610101000    Start time
`)
}
