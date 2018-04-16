package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/mitchellh/cli"
	"github.com/yyoshiki41/radigo"
)

func main() {
	v := fmt.Sprintf("%s built with %s",
		radigo.Version(), runtime.Version())
	c := cli.NewCLI("radigo", v)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"area":        radigo.AreaCommandFactory,
		"rec":         radigo.RecCommandFactory,
		"rec-live":    radigo.RecLiveCommandFactory,
		"browse":      radigo.BrowseCommandFactory,
		"browse-live": radigo.BrowseLiveCommandFactory,
	}

	exitCode, err := c.Run()
	if err != nil {
		log.Printf("Error executing CLI: %s", err.Error())
	}

	os.Exit(exitCode)
}
