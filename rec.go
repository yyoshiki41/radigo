package radigo

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/mitchellh/cli"
	"github.com/yyoshiki41/go-radiko"
)

type recCommand struct {
	ui cli.Ui
}

func (c *recCommand) Run(args []string) int {
	var stationID, start string
	var flagForce bool

	f := flag.NewFlagSet("rec", flag.ContinueOnError)
	f.StringVar(&stationID, "id", "", "id")
	f.StringVar(&start, "start", "", "start")
	f.StringVar(&start, "s", "", "start")
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
			"invalid start time format: %s", start))
		return 1
	}

	myPlayerPath := path.Join(radigoPath, "myplayer.swf")
	_, err = os.Stat(myPlayerPath)

	var dlErr error
	switch {
	case flagForce:
		if os.IsExist(err) {
			os.Remove(myPlayerPath)
		}
		dlErr = radiko.DownloadPlayer(myPlayerPath)
	case os.IsNotExist(err):
		dlErr = radiko.DownloadPlayer(myPlayerPath)
	}
	if dlErr != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to download player.swf: %s", dlErr))
		return 1
	}

	pngPath := path.Join(cachePath, "authkey.png")
	if err := swfExtract(myPlayerPath, pngPath); err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to execute swfextract: %s", err))
		return 1
	}

	client, err := radiko.New("")
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to construct a radiko Client.: %s", err))
		return 1
	}

	_, err = client.AuthorizeToken(context.Background(), pngPath)
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to get auth_token: %s", err))
		return 1
	}

	uri, err := client.TimeshiftPlaylistM3U8(context.Background(), stationID, startTime)
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

	err = bulkDownload(chunklist)
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to download aac files: %s", err))
		return 1
	}

	err = createConcatedAACFile()
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to create concat aac file: %s", err))
		return 1
	}

	err = convertAACToMP3()
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to convert aac to mp3: %s", err))
		return 1
	}

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
  -force,-f                Force flag (Do not use cache)
`)
}
