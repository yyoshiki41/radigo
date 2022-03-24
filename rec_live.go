package radigo

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/mitchellh/cli"
	"github.com/olekukonko/tablewriter"
	"github.com/yyoshiki41/go-radiko"
)

type recLiveCommand struct {
	ui cli.Ui
}

func (c *recLiveCommand) Run(args []string) int {
	var stationID, duration, areaID, fileType string
	var verbose bool

	f := flag.NewFlagSet("rec-live", flag.ContinueOnError)
	f.StringVar(&stationID, "id", "", "id")
	f.StringVar(&duration, "time", "", "duration")
	f.StringVar(&duration, "t", "", "duration")
	f.StringVar(&areaID, "area", "", "area")
	f.StringVar(&areaID, "a", "", "area")
	f.StringVar(&fileType, "output", AudioFormatAAC, "output")
	f.StringVar(&fileType, "o", AudioFormatAAC, "output")
	f.BoolVar(&verbose, "verbose", false, "verbose")
	f.BoolVar(&verbose, "v", false, "verbose")
	f.Usage = func() { c.ui.Error(c.Help()) }
	if err := f.Parse(args); err != nil {
		return 1
	}

	if stationID == "" {
		c.ui.Error("StationID is empty.")
		return 1
	}
	if duration == "" {
		c.ui.Error("Duration is empty.")
		return 1
	}

	output, err := NewOutputConfig(
		fmt.Sprintf("%s-%s", time.Now().In(location).Format(datetimeLayout), stationID),
		fileType)
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to configure output: %s", err))
		return 1
	}
	if err := output.SetupDir(); err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to setup the output dir: %s", err))
		return 1
	}
	if output.IsExist() {
		c.ui.Error(fmt.Sprintf(
			"the output file already exists: %s", output.AbsPath()))
		return 1
	}

	c.ui.Output("Now downloading.. ")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Station ID", "Duration(sec)"})
	table.Append([]string{stationID, duration})
	table.Render()

	spin := spinner.New(spinner.CharSets[9], time.Second)
	spin.Start()
	defer spin.Stop()

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	client, err := getClient(ctx, areaID)
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to construct a radiko Client: %s", err))
		return 1
	}

	_, err = client.AuthorizeToken(ctx)
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to get auth_token: %s", err))
		return 1
	}

	items, err := radiko.GetStreamSmhMultiURL(stationID)
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to get a stream url: %s", err))
		return 1
	}

	var streamURL string
	for _, i := range items {
		// Premium user
		if areaID != "" && areaID != currentAreaID {
			if i.Areafree {
				streamURL = i.PlaylistCreateURL
				break
			}
			continue
		}
		// Normal user
		if !i.Areafree {
			streamURL = i.PlaylistCreateURL
			break
		}
	}

	ffmpegCmd, err := newFfmpeg(ctx)
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to construct ffmpeg command: %s", err))
		return 1
	}

	ffmpegArgs := []string{
		"-loglevel", "quiet",
		"-fflags", "+discardcorrupt",
		"-headers", "X-Radiko-Authtoken: " + client.AuthToken(),
		"-t", duration,
		"-i", streamURL,
		"-vn",
		"-acodec",
	}
	switch fileType {
	case AudioFormatAAC:
		ffmpegArgs = append(ffmpegArgs, "copy")
	case AudioFormatMP3:
		ffmpegArgs = append(ffmpegArgs,
			[]string{"libmp3lame",
				"-ar", "44100",
				"-ab", "64k",
				"-ac", "2"}...)
	}
	ffmpegArgs = append(ffmpegArgs, output.AbsPath())
	ffmpegCmd.setArgs(ffmpegArgs...)

	err = ffmpegCmd.Run()
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to execute ffmpeg command: %s", err))
		return 1
	}

	c.ui.Output(fmt.Sprintf("Completed!\n%s", output.AbsPath()))

	return 0
}

func (c *recLiveCommand) Synopsis() string {
	return "Record a live program"
}

func (c *recLiveCommand) Help() string {
	return strings.TrimSpace(`
Usage: radigo rec-live [options]
  Record a live program.
Options:
  -id=name                 Station id
  -time,t=3600             Time duration (sec)
  -area,a=name             Area id
  -output,o=aac            Output file type (aac, mp3)
  -verbose,v               Verbose mode
`)
}
