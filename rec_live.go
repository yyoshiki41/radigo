package radigo

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/mitchellh/cli"
	radiko "github.com/yyoshiki41/go-radiko"
)

type recLiveCommand struct {
	ui cli.Ui
}

func (c *recLiveCommand) Run(args []string) int {
	var stationID, duration, areaID string
	var flagForce bool

	f := flag.NewFlagSet("rec", flag.ContinueOnError)
	f.StringVar(&stationID, "id", "", "id")
	f.StringVar(&duration, "time", "", "duration")
	f.StringVar(&duration, "t", "", "duration")
	f.StringVar(&areaID, "area", "", "area")
	f.StringVar(&areaID, "a", "", "area")
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
	if duration == "" {
		c.ui.Error("Duration is empty.")
		return 1
	}

	err := downloadSwfPlayer(flagForce)
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to download player.swf: %s", err))
		return 1
	}

	pngPath, err := extractPngFile()
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to execute swfextract: %s", err))
		return 1
	}

	client, err := getClient("", areaID)
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to construct a radiko Client: %s", err))
		return 1
	}
	_, err = client.AuthorizeToken(context.Background(), pngPath)
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to get auth_token: %s", err))
		return 1
	}

	items, err := radiko.GetStreamMultiURL(stationID)
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to get a stream url: %s", err))
		return 1
	}

	var streamURL string
	for _, i := range items {
		if !i.Areafree {
			streamURL = i.Item
			break
		}
	}

	rtmpdumpCmd, err := newRtmpdump(streamURL, client.AuthToken(), duration)
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to construct rtmpdump command: %s", err))
		return 1
	}

	ffmpegCmd, err := newFfmpeg("-")
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to construct ffmpeg command: %s", err))
		return 1
	}
	ffmpegCmd.setArgs(
		"-vn",
		"-acodec", "libmp3lame",
		"-ar", "44100",
		"-ab", "64k",
		"-ac", "2",
	)

	ffmpegCmd.Stdin, err = rtmpdumpCmd.StdoutPipe()
	if err != nil {
		c.ui.Error(fmt.Sprintf("%v", err))
		return 1
	}
	ffmpegCmd.Stdout = os.Stdout

	output := path.Join(radigoPath, "result.mp3")
	err = ffmpegCmd.start(output)
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to start ffmpeg command: %s", err))
		return 1
	}

	// TODO:
	// context 使って、Runが失敗したらffmpegCmdをkillする
	go rtmpdumpCmd.Run()

	err = ffmpegCmd.Wait()
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to execute ffmpeg command: %s", err))
		return 1
	}

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
  -force,-f                Ignore cache and force refresh
`)
}
