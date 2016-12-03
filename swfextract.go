package radigo

import (
	"os"
	"os/exec"
	"path"

	"github.com/yyoshiki41/go-radiko"
)

var (
	swfPlayer = path.Join(radigoPath, "myplayer.swf")
)

func downloadSwfPlayer(flagForce bool) error {
	_, err := os.Stat(swfPlayer)
	if flagForce && os.IsExist(err) {
		os.Remove(swfPlayer)
	}

	if flagForce || os.IsNotExist(err) {
		err := radiko.DownloadPlayer(swfPlayer)
		if err != nil {
			return err
		}
	}
	return nil
}

func extractPngFile() (string, error) {
	pngPath := path.Join(cachePath, "authkey.png")
	if err := swfExtract(swfPlayer, pngPath); err != nil {
		return "", err
	}

	return pngPath, nil
}

func swfExtract(swfPath, output string) error {
	cmdPath, err := exec.LookPath("swfextract")
	if err != nil {
		return err
	}

	swfExtract := exec.Command(cmdPath, "-b", "12", swfPath, "-o", output)
	return swfExtract.Run()
}
