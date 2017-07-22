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
)

type recLiveCommand struct {
	ui cli.Ui
}

func (c *recLiveCommand) Run(args []string) int {
	var (
		stationID, duration, areaID string
		flagForce                   bool
	)

	f := flag.NewFlagSet("rec-live", flag.ContinueOnError)
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

	// if expiration date has passed, remove cache files
	if flagForce || isExpiredCache() {
		removeTokenCache()
		c.ui.Info("Delete token cache.")
	}

	c.ui.Output("Now downloading.. ")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Station ID", "Duration(sec)"})
	table.Append([]string{stationID, duration})
	table.Render()

	spin := spinner.New(spinner.CharSets[9], time.Second)
	spin.Start()
	defer spin.Stop()

	if err := downloadSwfPlayer(flagForce); err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to download player.swf: %s", err))
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
		token, err := client.AuthorizeToken(ctx)
		if err != nil {
			c.ui.Error(fmt.Sprintf(
				"Failed to get auth_token: %s", err))
			return 1
		}
		if err := saveToken(token); err != nil {
			c.ui.Info("Save token cache.")
		}
	}

	items, err := radiko.GetStreamMultiURL(stationID)
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
				streamURL = i.Item
				break
			}
			continue
		}
		// Normal user
		if !i.Areafree {
			streamURL = i.Item
			break
		}
	}

	rtmpdumpCmd, err := newRtmpdump(ctx, streamURL, client.AuthToken(), duration)
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to construct rtmpdump command: %s", err))
		return 1
	}

	ffmpegCmd, err := newFfmpeg(ctx, "-")
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

	outputFile := path.Join(radigoPath, "output",
		fmt.Sprintf("%s-%s.mp3",
			time.Now().In(location).Format(datetimeLayout), stationID,
		))

	err = ffmpegCmd.start(outputFile)
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to start ffmpeg command: %s", err))
		return 1
	}

	go func() {
		err := rtmpdumpCmd.Run()
		if err != nil {
			ctxCancel()
			c.ui.Error(fmt.Sprintf(
				"Failed to execute rtmpdump command: %s", err))
		}
	}()

	err = ffmpegCmd.Wait()
	if err != nil {
		c.ui.Error(fmt.Sprintf(
			"Failed to execute ffmpeg command: %s", err))
		return 1
	}

	c.ui.Output(fmt.Sprintf("Completed!\n%s", outputFile))

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
