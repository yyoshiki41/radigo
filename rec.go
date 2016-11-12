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

	f := flag.NewFlagSet("rec", flag.ContinueOnError)
	f.StringVar(&stationID, "id", "", "id")
	f.StringVar(&start, "start", "", "start")
	f.StringVar(&start, "s", "", "start")
	f.Usage = func() { c.ui.Error(c.Help()) }
	if err := f.Parse(args); err != nil {
		return 1
	}

	startTime, err := time.ParseInLocation(datetimeLayout, start, location)
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"invalid start time format: %s", start))
		return 1
	}

	myPlayerPath := path.Join(radigoPath, "myplayer.swf")
	if _, err := os.Stat(myPlayerPath); err != nil {
		/* TODO: force option
		if os.IsExist(err) {
			os.Remove(myPlayerPath)
		}
		*/
		if err := radiko.DownloadPlayer(myPlayerPath); err != nil {
			c.ui.Error(fmt.Sprintf(
				"Failed to download player.swf: %s", err))
			return 1
		}
	}

	client, err := radiko.New("")
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to construct a radiko Client.: %s", err))
		return 1
	}

	pngPath := path.Join(cachePath, "authkey.png")
	if err := swfExtract(myPlayerPath, pngPath); err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to execute swfextract: %s", err))
		return 1
	}

	if _, err = client.AuthorizeToken(context.Background(), pngPath); err != nil {
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
`)
}
