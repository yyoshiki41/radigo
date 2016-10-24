package radigo

import (
	"encoding/base64"
	"os"
	"os/exec"
)

type pngService struct {
	pngPath    string
	swfextract *exec.Cmd
}

func newPngService(swfFilePath, output string) (*pngService, error) {
	cmdPath, err := exec.LookPath("swfextract")
	if err != nil {
		return nil, err
	}

	return &pngService{
		pngPath:    output,
		swfextract: exec.Command(cmdPath, "-b", "12", swfFilePath, "-o", output),
	}, nil
}

func (p *pngService) Create() error {
	return p.swfextract.Run()
}

func (p *pngService) getPartialKey(length, offset int64) (string, error) {
	f, err := os.Open(p.pngPath)
	if err != nil {
		return "", err
	}

	b := make([]byte, length)
	if _, err = f.ReadAt(b, offset); err != nil {
		return "", err
	}
	partialKey := base64.StdEncoding.EncodeToString(b)
	return partialKey, nil
}
