package radigo

import (
	"os"

	"github.com/mitchellh/cli"
)

var Ui cli.Ui

func init() {
	Ui = &cli.BasicUi{
		Writer: os.Stdout,
	}
}

func RecCommandFactory() (cli.Command, error) {
	return &recCommand{
		ui: Ui,
	}, nil
}
