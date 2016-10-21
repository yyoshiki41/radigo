package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/yyoshiki41/radigo"
)

var (
	// ErrUsage is returned when a usage message was printed and the process
	// should simply exit with an error.
	ErrUsage = errors.New("usage")

	// ErrUnknownCommand is returned when a CLI command is not specified.
	ErrUnknownCommand = errors.New("unknown command")
)

func main() {
	if len(os.Args) == 0 {
		fmt.Fprintln(os.Stderr, usage())
	}

	err := run(os.Args[1:]...)
	if err == ErrUsage {
		os.Exit(2)
	}
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run(args ...string) error {
	// Print the version.
	if args[0] == "version" || args[0] == "-v" {
		fmt.Fprintln(os.Stderr, radigo.Version())
		return nil
	}

	// Require a command at the beginning.
	if strings.HasPrefix(args[0], "-") {
		fmt.Fprintln(os.Stderr, usage())
		return ErrUsage
	}

	// Execute command.
	switch args[0] {
	case "help":
		fmt.Fprintln(os.Stderr, usage())
		return ErrUsage
	case "rec":
		return nil
	case "search":
		return nil
	default:
		return ErrUnknownCommand
	}
}

func usage() string {
	return strings.TrimLeft(`
radigo is a tool for recording radiko.
Usage:
    radigo command [arguments]

The commands are:
    help        Print this screen.
    version     Print the version.
    rec         Record a radiko program.
    search      Search a radiko program.

Use "radigo [command] -h" for more information about a command.
`, "\n")
}
