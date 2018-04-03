package radiko

import (
	"context"
	"path"
	"time"

	"github.com/yyoshiki41/go-radiko/internal/m3u8"
	"github.com/yyoshiki41/go-radiko/internal/util"
)

// TimeshiftPlaylistM3U8 returns uri.
func (c *Client) TimeshiftPlaylistM3U8(ctx context.Context, stationID string, start time.Time) (string, error) {
	prog, err := c.GetProgramByStartTime(ctx, stationID, start)
	if err != nil {
		return "", err
	}

	apiEndpoint := apiPath(apiV2, "ts/playlist.m3u8")
	req, err := c.newRequest(ctx, "POST", apiEndpoint, &Params{
		query: map[string]string{
			"station_id": stationID,
			"ft":         prog.Ft,
			"to":         prog.To,
			"l":          "15", // must?
		},
		setAuthToken: true,
	})

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return m3u8.GetURI(resp.Body)
}

// GetTimeshiftURL returns a timeshift url for web browser.
func GetTimeshiftURL(stationID string, start time.Time) string {
	endpoint := path.Join("#!/ts", stationID, util.Datetime(start))
	return defaultEndpoint + "/" + endpoint
}
