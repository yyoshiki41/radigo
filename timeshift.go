package radigo

import (
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/grafov/m3u8"
)

const chunkSeconds = 300

type streamURLs struct {
	URLs []streamURL `xml:"url"`
}

type streamURL struct {
	Areafree          int    `xml:"areafree,attr"`
	Timefree          int    `xml:"timefree,attr"`
	PlaylistCreateURL string `xml:"playlist_create_url"`
}

func generateLsid() string {
	b := make([]byte, 16)
	for i := range b {
		b[i] = byte(rand.Intn(256))
	}
	return fmt.Sprintf("%x", md5.Sum(b))
}

func getTimefreeStreamURL(stationID string, areafree bool) (string, error) {
	apiURL := fmt.Sprintf("https://radiko.jp/v3/station/stream/pc_html5/%s.xml", stationID)
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get stream XML: status %d", resp.StatusCode)
	}

	var urls streamURLs
	if err := xml.NewDecoder(resp.Body).Decode(&urls); err != nil {
		return "", fmt.Errorf("failed to decode stream XML: %w", err)
	}

	wantAreafree := 0
	if areafree {
		wantAreafree = 1
	}

	for _, u := range urls.URLs {
		if u.Timefree == 1 && u.Areafree == wantAreafree {
			return u.PlaylistCreateURL, nil
		}
	}
	for _, u := range urls.URLs {
		if u.Timefree == 1 {
			return u.PlaylistCreateURL, nil
		}
	}
	return "", fmt.Errorf("no timefree stream URL found for station %s", stationID)
}

func isAbsoluteURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.Scheme != ""
}

func resolveReference(baseURL, ref string) (string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	r, err := url.Parse(ref)
	if err != nil {
		return "", err
	}
	return base.ResolveReference(r).String(), nil
}

// fetchMedialistURL fetches the master playlist and returns the medialist URI.
func fetchMedialistURL(playlistBaseURL, stationID, startAt, ft, to, authToken string, l int) (string, error) {
	u, err := url.Parse(playlistBaseURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("station_id", stationID)
	q.Set("start_at", startAt)
	q.Set("ft", ft)
	q.Set("end_at", to)
	q.Set("to", to)
	q.Set("l", fmt.Sprintf("%d", l))
	q.Set("lsid", generateLsid())
	q.Set("type", "b")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("X-Radiko-AuthToken", authToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get timeshift playlist: status %d", resp.StatusCode)
	}

	baseForResolve := resp.Request.URL.String()

	playlist, listType, err := m3u8.DecodeFrom(resp.Body, true)
	if err != nil {
		return "", fmt.Errorf("failed to decode m3u8: %w", err)
	}
	if listType != m3u8.MASTER {
		return "", fmt.Errorf("unexpected playlist type: %d", listType)
	}
	p := playlist.(*m3u8.MasterPlaylist)
	if p == nil || len(p.Variants) == 0 || p.Variants[0] == nil {
		return "", fmt.Errorf("invalid m3u8 format: no variants")
	}

	variantURI := p.Variants[0].URI
	if !isAbsoluteURL(variantURI) {
		variantURI, err = resolveReference(baseForResolve, variantURI)
		if err != nil {
			return "", fmt.Errorf("failed to resolve variant URI: %w", err)
		}
	}
	return variantURI, nil
}

// fetchChunklist fetches a medialist and returns segment URLs.
func fetchChunklist(uri, authToken string) ([]string, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Radiko-AuthToken", authToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get chunklist: status %d", resp.StatusCode)
	}

	baseForResolve := resp.Request.URL.String()

	playlist, listType, err := m3u8.DecodeFrom(resp.Body, true)
	if err != nil {
		return nil, fmt.Errorf("failed to decode m3u8: %w", err)
	}
	if listType != m3u8.MEDIA {
		return nil, fmt.Errorf("unexpected playlist type: %d", listType)
	}
	p := playlist.(*m3u8.MediaPlaylist)

	var chunklist []string
	for _, v := range p.Segments {
		if v != nil {
			chunkURI := v.URI
			if !isAbsoluteURL(chunkURI) {
				chunkURI, err = resolveReference(baseForResolve, chunkURI)
				if err != nil {
					return nil, fmt.Errorf("failed to resolve chunk URI: %w", err)
				}
			}
			chunklist = append(chunklist, chunkURI)
		}
	}
	return chunklist, nil
}

// getTimeshiftChunklist fetches all segment URLs for a timefree program
// by splitting the program into chunks (the API limits each request).
func getTimeshiftChunklist(stationID, ft, to, authToken string, areafree bool) ([]string, error) {
	playlistBaseURL, err := getTimefreeStreamURL(stationID, areafree)
	if err != nil {
		return nil, err
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, err
	}
	ftTime, err := time.ParseInLocation(datetimeLayout, ft, loc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ft: %w", err)
	}
	toTime, err := time.ParseInLocation(datetimeLayout, to, loc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse to: %w", err)
	}
	totalDuration := int(toTime.Sub(ftTime).Seconds())

	var allChunks []string
	for offset := 0; offset < totalDuration; offset += chunkSeconds {
		chunkStart := ftTime.Add(time.Duration(offset) * time.Second)
		l := chunkSeconds
		if offset+l > totalDuration {
			l = totalDuration - offset
		}

		startAt := chunkStart.Format(datetimeLayout)
		mediaURL, err := fetchMedialistURL(playlistBaseURL, stationID, startAt, ft, to, authToken, l)
		if err != nil {
			return nil, fmt.Errorf("failed to get medialist (offset=%d): %w", offset, err)
		}

		chunks, err := fetchChunklist(mediaURL, authToken)
		if err != nil {
			return nil, fmt.Errorf("failed to get chunklist (offset=%d): %w", offset, err)
		}
		allChunks = append(allChunks, chunks...)
	}
	return allChunks, nil
}
