package m3u8

import (
	"bufio"
	"os"
	"path/filepath"
	"testing"
)

func readTestData(fileName string) *os.File {
	const testDir = "github.com/yyoshiki41/go-radiko/testdata"

	GOPATH := os.Getenv("GOPATH")
	f, err := os.Open(filepath.Join(GOPATH, "src", testDir, fileName))
	if err != nil {
		panic(err)
	}
	return f
}

func TestGetURI(t *testing.T) {
	expected := "https://radiko.jp/v2/api/ts/chunklist/NejwTOkX.m3u8"

	input := bufio.NewReader(readTestData("uri.m3u8"))
	u, err := GetURI(input)
	if err != nil {
		t.Error(err)
	}
	if u != expected {
		t.Errorf("expected %s, but %s", expected, u)
	}
}

func TestGetChunklist(t *testing.T) {
	input := bufio.NewReader(readTestData("chunklist.m3u8"))
	chunklist, err := GetChunklist(input)
	if err != nil {
		t.Error(err)
	}
	if len(chunklist) == 0 {
		t.Error("chunklist is empty.")
	}
}
