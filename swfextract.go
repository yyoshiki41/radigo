package radigo

import "os/exec"

func swfExtract(swfPath, output string) error {
	cmdPath, err := exec.LookPath("swfextract")
	if err != nil {
		return err
	}

	swfExtract := exec.Command(cmdPath, "-b", "12", swfPath, "-o", output)
	return swfExtract.Run()
}
