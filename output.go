package radigo

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/yyoshiki41/radigo/internal"
)

var aacResultFile string

func initTempAACDir() (string, error) {
	aacDir, err := ioutil.TempDir(radigoPath, "aac")
	if err != nil {
		return "", err
	}

	aacResultFile = filepath.Join(aacDir, "result.aac")
	return aacDir, nil
}

func createConcatedAACFile(ctx context.Context, aacDir string) error {
	name, err := internal.ConcatFileNames(aacDir)
	if err != nil {
		return err
	}

	f, err := newFfmpeg(ctx, fmt.Sprintf("concat:%s", name))
	if err != nil {
		return err
	}

	f.setDir(aacDir)
	f.setArgs("-c", "copy")
	// TODO: Run 結果の標準出力を拾う
	return f.run(aacResultFile)
}

func output(ctx context.Context, fileType, outputFile string) error {
	switch fileType {
	case "mp3":
		return outputMP3(ctx, outputFile)
	case "aac":
		return outputAAC(outputFile)
	}
	return fmt.Errorf("Unsupported file type: %s", fileType)
}

func outputAAC(outputFile string) error {
	if err := os.Rename(aacResultFile, outputFile); err != nil {
		return err
	}
	return nil
}

func outputMP3(ctx context.Context, outputFile string) error {
	f, err := newFfmpeg(ctx, aacResultFile)
	if err != nil {
		return err
	}

	f.setDir(radigoPath)
	f.setArgs(
		"-c:a", "libmp3lame",
		"-ac", "2",
		"-q:a", "2",
		"-y", // overwrite the output file without asking
	)
	// TODO: Run 結果の標準出力を拾う
	return f.run(outputFile)
}
