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

func AreaCommandFactory() (cli.Command, error) {
	return &areaCommand{
		ui: Ui,
	}, nil
}

func RecCommandFactory() (cli.Command, error) {
	return &recCommand{
		ui: Ui,
	}, nil
}

func RecLiveCommandFactory() (cli.Command, error) {
	return &recLiveCommand{
		ui: Ui,
	}, nil
}
