package radigo

import (
	"strings"

	"github.com/mitchellh/cli"
)

type recCommand struct {
	ui cli.Ui
}

func (c *recCommand) Run(args []string) int {
	return 0
}

func (c *recCommand) Help() string {
	return strings.TrimSpace(`
Usage: terraform rec [options]
  Record a radiko program.
Options:
  -key=name              Key station
  -start=201610101000    Start time
  -end=201610101200      Start time
`)
}

func (c *recCommand) Synopsis() string {
	return "Record a radiko program"
}
