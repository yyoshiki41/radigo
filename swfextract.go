package radigo

import (
	"os"
	"os/exec"
	"path"
)

var (
	swfPlayer = path.Join(radigoPath, "myplayer.swf")
	pngFile   = path.Join(cachePath, "authkey.png")
)

func extractPngFile(flagForce bool) error {
	_, err := os.Stat(pngFile)
	if flagForce && os.IsExist(err) {
		os.Remove(pngFile)
	}

	if flagForce || os.IsNotExist(err) {
		err := swfExtract(swfPlayer, pngFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func swfExtract(swfPath, output string) error {
	cmdPath, err := exec.LookPath("swfextract")
	if err != nil {
		return err
	}

	swfExtract := exec.Command(cmdPath, "-b", "12", swfPath, "-o", output)
	return swfExtract.Run()
}
