package radigo

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/yyoshiki41/go-radiko"
)

type browseLiveCommand struct {
	ui cli.Ui
}

func (c *browseLiveCommand) Run(args []string) int {
	var stationID string

	f := flag.NewFlagSet("browse-live", flag.ContinueOnError)
	f.StringVar(&stationID, "id", "", "id")
	f.Usage = func() { c.ui.Error(c.Help()) }
	if err := f.Parse(args); err != nil {
		return 1
	}

	if stationID == "" {
		c.ui.Error("StationID is empty.")
		return 1
	}

	url := radiko.GetLiveURL(stationID)
	cmd := exec.Command("open", url)
	if err := cmd.Run(); err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to open browser: %s", err))
	}

	return 0
}

func (c *browseLiveCommand) Synopsis() string {
	return "Browse radiko.jp live"
}

func (c *browseLiveCommand) Help() string {
	return strings.TrimSpace(`
Usage: radigo browse-live [options]
  Browse radiko.jp live
Options:
  -id=name                 Station id
`)
}
