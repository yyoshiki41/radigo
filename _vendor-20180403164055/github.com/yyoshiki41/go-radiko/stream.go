package radiko

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
)

// URLItem represents a stream url.
type URLItem struct {
	Areafree bool   `xml:"areafree,attr"`
	Item     string `xml:",chardata"`
}

// GetStreamMultiURL returns a slice of the stream url.
func GetStreamMultiURL(stationID string) ([]URLItem, error) {
	endpoint := path.Join(apiV2, "station/stream_multi",
		fmt.Sprintf("%s.xml", stationID))

	resp, err := http.Get("http://radiko.jp/" + endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	urlData := streamURLData{}
	if err = xml.Unmarshal(b, &urlData); err != nil {
		return nil, err
	}
	return urlData.Items, err
}

type streamURLData struct {
	XMLName xml.Name  `xml:"url"`
	Items   []URLItem `xml:"item"`
}

// GetLiveURL returns a live url for web browser.
func GetLiveURL(stationID string) string {
	endpoint := path.Join("#!/live", stationID)
	return defaultEndpoint + "/" + endpoint
}
