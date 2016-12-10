package radigo

import (
	"os"

	"github.com/mitchellh/cli"
)

var Ui cli.Ui

const (
	outputPrefix = "OUTPUT: "
	infoPrefix   = "INFO: "
	errorPrefix  = "ERROR: "
	warnPrefix   = "WARN: "
)

func init() {
	Ui = &cli.PrefixedUi{
		OutputPrefix: outputPrefix,
		InfoPrefix:   infoPrefix,
		ErrorPrefix:  errorPrefix,
		WarnPrefix:   warnPrefix,
		Ui: &cli.BasicUi{
			Writer: os.Stdout,
		},
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
