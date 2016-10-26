package radigo

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
)

const (
	playerURL  = "http://radiko.jp/apps/js/flash/myplayer-release.swf"
	baseURL    = "https://radiko.jp"
	apiVersion = "v2"
)

func downloadPlayer(path string) error {
	resp, err := http.Get(playerURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, resp.Body)
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	return err
}

type radiko struct {
	client    *http.Client
	basePath  string
	stationID string
}

func newRadiko(stationID string) *radiko {
	return &radiko{
		client:    &http.Client{},
		basePath:  baseURL,
		stationID: stationID,
	}
}

func (r *radiko) buildEndpoint(p string) (string, error) {
	u, err := url.Parse(r.basePath)
	if err != nil {
		return "", fmt.Errorf("Parse given URL: %#v", err)
	}

	u.Path = path.Join(u.Path, p)
	return u.String(), err
}

func (r *radiko) buildAPIEndpoint(p string) (string, error) {
	apiPath := path.Join(apiVersion, "api", p)
	return r.buildEndpoint(apiPath)
}

func (r *radiko) auth1_fms(myPlayerPath string) (string, string, error) {
	apiEndpoint, err := r.buildAPIEndpoint("auth1_fms")
	if err != nil {
		return "", "", err
	}

	// TODO: 関数として切り出す
	req, err := http.NewRequest("POST", apiEndpoint, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("X-Radiko-App", "pc_ts")
	req.Header.Set("X-Radiko-App-Version", "4.0.0")
	req.Header.Set("X-Radiko-User", "test-stream")
	req.Header.Set("X-Radiko-Device", "pc")

	resp, err := r.client.Do(req)
	defer resp.Body.Close()

	authToken := resp.Header.Get("X-Radiko-Authtoken")
	keyLength := resp.Header.Get("X-Radiko-Keylength")
	keyOffset := resp.Header.Get("X-Radiko-Keyoffset")

	length, err := strconv.ParseInt(keyLength, 10, 64)
	if err != nil {
		return "", "", err
	}
	offset, err := strconv.ParseInt(keyOffset, 10, 64)
	if err != nil {
		return "", "", err
	}

	pngPath := path.Join(cachePath, "authkey.png")
	p, err := newPngService(myPlayerPath, pngPath)
	if err != nil {
		return "", "", err
	}
	if err := p.Create(); err != nil {
		return "", "", err
	}
	partialKey, err := p.getPartialKey(length, offset)
	if err != nil {
		return "", "", err
	}

	return authToken, partialKey, err
}

func (r *radiko) auth2_fms(authToken, partialKey string) ([]string, error) {
	apiEndpoint, err := r.buildAPIEndpoint("auth2_fms")
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", apiEndpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("X-Radiko-App", "pc_ts")
	req.Header.Set("X-Radiko-App-Version", "4.0.0")
	req.Header.Set("X-Radiko-User", "test-stream")
	req.Header.Set("X-Radiko-Device", "pc")
	req.Header.Set("X-Radiko-Authtoken", authToken)
	req.Header.Set("X-Radiko-Partialkey", partialKey)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	s := strings.Split(string(b), ",")
	return s, nil
}

func (r *radiko) playlistM3U8(authToken, start, end string) (string, error) {
	apiEndpoint, err := r.buildAPIEndpoint("ts/playlist.m3u8")
	u, err := url.Parse(apiEndpoint)
	if err != nil {
		return "", fmt.Errorf("Parse given URL: %#v", err)
	}
	q := u.Query()
	q.Set("station_id", r.stationID)
	q.Set("ft", start)
	q.Set("to", end)
	q.Set("l", "15") // must?
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("X-Radiko-Authtoken", authToken)

	resp, err := r.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
