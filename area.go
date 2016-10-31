package radigo

import (
	"flag"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
)

type areaCommand struct {
	ui cli.Ui
}

func (c *areaCommand) Run(args []string) int {
	var areaID string

	f := flag.NewFlagSet("area", flag.ContinueOnError)
	f.StringVar(&areaID, "id", "", "id")
	f.Usage = func() { c.ui.Error(c.Help()) }
	err := f.Parse(args)
	if err != nil {
		return 1
	}

	if areaID == "" {
		var err error
		areaID, err = getAreaID()
		if err != nil {
			c.ui.Error(fmt.Sprintf(
				"Failed to get area id: %s", err))
			return 1
		}
	}

	// TODO: radiko 構造体
	r := newRadiko("")
	pr, err := r.getProgram(areaID)
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to get programs: %s", err))
		return 1
	}
	for _, s := range pr.Stations.Stations {
		c.ui.Output(fmt.Sprintf(
			"%s\n  - %s", s.Name, s.ID))
	}

	return 0
}

func (c *areaCommand) Synopsis() string {
	return "Get available station ids"
}

func (c *areaCommand) Help() string {
	return strings.TrimSpace(`
Usage: radigo area [options]
  Get available stations ids.
Options:
  -id=name                 Area id
`)
}
