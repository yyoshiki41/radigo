package radiko

import (
	"net/http"

	"github.com/yyoshiki41/go-radiko/internal/m3u8"
)

// GetChunklistFromM3U8 returns a slice of url.
func GetChunklistFromM3U8(uri string) ([]string, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return m3u8.GetChunklist(resp.Body)
}
