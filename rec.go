package radigo

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/mitchellh/cli"
)

type recCommand struct {
	ui cli.Ui
}

func (c *recCommand) Run(args []string) int {
	var stationID, start, end string

	f := flag.NewFlagSet("rec", flag.ContinueOnError)
	f.StringVar(&stationID, "id", "", "id")
	f.StringVar(&start, "start", "", "start")
	f.StringVar(&start, "s", "", "start")
	f.StringVar(&end, "end", "", "end")
	f.StringVar(&end, "e", "", "end")
	f.Usage = func() { c.ui.Error(c.Help()) }
	if err := f.Parse(args); err != nil {
		return 1
	}

	myPlayerPath := path.Join(radigoPath, "myplayer.swf")
	if _, err := os.Stat(myPlayerPath); err != nil {
		/* TODO: force option
		if os.IsExist(err) {
			os.Remove(myPlayerPath)
		}
		*/
		if err := downloadPlayer(myPlayerPath); err != nil {
			c.ui.Error(fmt.Sprintf(
				"Failed to download player.swf: %s", err))
			return 1
		}
	}

	r := newRadiko(stationID)
	authToken, partialKey, err := r.auth1_fms(myPlayerPath)
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to get auth token and key: %s", err))
		return 1
	}

	_, err = r.auth2_fms(authToken, partialKey)
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to auth2_fms: %s", err))
		return 1
	}

	fmt.Println(authToken)
	res, err := r.playlistM3U8(authToken, start, end)
	fmt.Println(res)

	return 0
}

func (c *recCommand) Synopsis() string {
	return "Record a radiko program"
}

func (c *recCommand) Help() string {
	return strings.TrimSpace(`
Usage: radiko rec [options]
  Record a radiko program.
Options:
  -id=name                 Station id
  -start,s=201610101000    Start time
  -end,e=201610101200      End time
`)
}
