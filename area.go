package radigo

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/olekukonko/tablewriter"
	"github.com/yyoshiki41/go-radiko"
)

type areaCommand struct {
	ui cli.Ui
}

func (c *areaCommand) Run(args []string) int {
	var areaID string

	f := flag.NewFlagSet("area", flag.ContinueOnError)
	f.StringVar(&areaID, "id", "", "id")
	f.Usage = func() { c.ui.Error(c.Help()) }
	if err := f.Parse(args); err != nil {
		return 1
	}

	client, err := radiko.New("")
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to construct a radiko Client.: %s", err))
		return 1
	}

	if areaID != "" {
		client.SetAreaID(areaID)
	}

	stations, err := client.GetNowPrograms(context.Background())
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to get stations: %s", err))
		return 1
	}

	c.ui.Output(fmt.Sprintf(" Area ID: %s", client.AreaID()))
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Station ID"})
	for _, s := range stations {
		table.Append([]string{s.Name, s.ID})
	}
	table.Render()

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
