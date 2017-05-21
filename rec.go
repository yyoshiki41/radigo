package radigo

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/mitchellh/cli"
	"github.com/olekukonko/tablewriter"
	"github.com/yyoshiki41/go-radiko"
	"github.com/yyoshiki41/radigo/internal"
)

type recCommand struct {
	ui cli.Ui
}

func (c *recCommand) Run(args []string) int {
	var (
		stationID, start, areaID, fileType string
		flagForce                          bool
	)

	f := flag.NewFlagSet("rec", flag.ContinueOnError)
	f.StringVar(&stationID, "id", "", "id")
	f.StringVar(&start, "start", "", "start")
	f.StringVar(&start, "s", "", "start")
	f.StringVar(&areaID, "area", "", "area")
	f.StringVar(&areaID, "a", "", "area")
	f.StringVar(&fileType, "output", "", "output")
	f.StringVar(&fileType, "o", "", "output")
	f.BoolVar(&flagForce, "force", false, "force")
	f.BoolVar(&flagForce, "f", false, "force")
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
	if fileType == "" {
		fileType = "aac"
	}
	if fileType != "mp3" && fileType != "aac" {
		c.ui.Error(fmt.Sprintf(
			"Unsupported file type: %s", fileType))
		return 1
	}

	if flagForce || isExpiredCache() {
		removeTokenCache()
		c.ui.Info("Delete token cache.")
	}

	c.ui.Output("Now downloading.. ")
	spin := spinner.New(spinner.CharSets[9], time.Second)
	spin.Start()
	defer spin.Stop()

	if err := downloadSwfPlayer(flagForce); err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to download player.swf: %s", err))
		return 1
	}

	if err := extractPngFile(flagForce); err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to execute swfextract: %s", err))
		return 1
	}

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	client, err := getClient(ctx, areaID)
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to construct a radiko Client: %s", err))
		return 1
	}
	if client.AuthToken() == "" {
		token, err := client.AuthorizeToken(ctx, pngFile)
		if err != nil {
			c.ui.Error(fmt.Sprintf(
				"Failed to get auth_token: %s", err))
			return 1
		}
		if err := saveToken(token); err != nil {
			c.ui.Info("Save token cache.")
		}
	}

	go func() {
		pg, err := client.GetProgramByStartTime(ctx, stationID, startTime)
		if err != nil {
			ctxCancel()
			c.ui.Error(fmt.Sprintf(
				"Failed to get the program: %s", err))
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"STATION ID", "TITLE"})
		table.Append([]string{stationID, pg.Title})
		fmt.Print("\n")
		table.Render()
	}()

	uri, err := client.TimeshiftPlaylistM3U8(ctx, stationID, startTime)
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to get playlist.m3u8: %s", err))
		return 1
	}

	chunklist, err := radiko.GetChunklistFromM3U8(uri)
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to get chunklist: %s", err))
		return 1
	}

	aacDir, err := initTempAACDir()
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to create the aac dir: %s", err))
		return 1
	}
	defer os.RemoveAll(aacDir)

	if err := internal.BulkDownload(chunklist, aacDir); err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to download aac files: %s", err))
		return 1
	}
	if err := createConcatedAACFile(ctx, aacDir); err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to concat aac files: %s", err))
		return 1
	}

	outputFile := path.Join(radigoPath, "output",
		fmt.Sprintf("%s-%s.%s",
			startTime.In(location).Format(datetimeLayout), stationID, fileType,
		))

	if err := output(ctx, fileType, outputFile); err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to output a result file: %s", err))
		return 1
	}

	c.ui.Output(fmt.Sprintf("Completed!\n%s", outputFile))

	return 0
}

func (c *recCommand) Synopsis() string {
	return "Record a radiko program"
}

func (c *recCommand) Help() string {
	return strings.TrimSpace(`
Usage: radigo rec [options]
  Record a radiko program.
Options:
  -id=name                 Station id
  -start,s=201610101000    Start time
  -output,o=mp3            Output file type (mp3, aac)
  -area,a=name             Area id
  -force,f                 Ignore cache and force refresh
`)
}
