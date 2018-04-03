package radiko

import (
	"compress/zlib"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	playerURL = "http://radiko.jp/apps/js/flash/myplayer-release.swf"

	// swfextract
	targetID     = 12 // swfextract -b "12"
	targetCode   = 87 // swfextract "-b" 12
	headerCWS    = 8
	headerRect   = 5
	rectNum      = 4
	headerRest   = 2 + 2
	binaryOffset = 6
)

// DownloadPlayer downloads a swf player file.
func DownloadPlayer(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	resp, err := http.Get(playerURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(f, resp.Body)
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	return err
}

func downloadBinary() ([]byte, error) {
	resp, err := http.Get(playerURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return swfExtract(resp.Body)
}

func swfExtract(body io.Reader) ([]byte, error) {
	io.CopyN(ioutil.Discard, body, headerCWS)
	zf, err := zlib.NewReader(body)
	if err != nil {
		return nil, err
	}
	buf, err := ioutil.ReadAll(zf)
	if err != nil {
		return nil, err
	}

	offset := 0

	// Skip Rect
	rectSize := int(buf[offset] >> 3)
	rectOffset := (headerRect + rectNum*rectSize + 7) / 8

	offset += rectOffset

	// Skip the rest header
	offset += headerRest

	// Read tags
	for i := 0; ; i++ {
		// tag code
		code := int(buf[offset+1])<<2 + int(buf[offset])>>6

		// tag length
		len := int(buf[offset] & 0x3f)

		// Skip tag header
		offset += 2

		// tag length (if long version)
		if len == 0x3f {
			len = int(buf[offset])
			len += int(buf[offset+1]) << 8
			len += int(buf[offset+2]) << 16
			len += int(buf[offset+3]) << 24

			// skip tag lentgh header
			offset += 4
		}

		// Not found...
		if code == 0 {
			return nil, errors.New("swf extract failed")
		}
		// tag ID
		id := int(buf[offset]) + int(buf[offset+1])<<8

		// Found?
		if code == targetCode && id == targetID {
			return buf[offset+binaryOffset : offset+len], nil
		}

		offset += len
	}
}
