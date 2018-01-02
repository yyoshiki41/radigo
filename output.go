package radigo

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	// "github.com/yyoshiki41/radigo/internal"
	"github.com/soeyusuke/radigo/internal"
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
	aactempdir, err := ioutil.TempDir(radigoPath, "aac")
	if err != nil {
		return err
	}
	defer os.RemoveAll(aactempdir)

	if err := RemakeAAC(ctx, aacDir, aactempdir); err != nil {
		return err
	}

	name, err := internal.ConcatFileNames(aactempdir)
	if err != nil {
		return err
	}

	f, err := newFfmpeg(ctx, fmt.Sprintf("concat:%s", name))
	if err != nil {
		return err
	}

	f.setDir(aactempdir)
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

func RemakeAAC(ctx context.Context, aacDir, tempdir string) error {
	var (
		wg      sync.WaitGroup
		errChan = make(chan error, 1)
	)

	files, err := ioutil.ReadDir(aacDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		wg.Add(1)
		go func(fname, tempdir string) {
			defer wg.Done()

			buf, err := AACtoByte(fname, aacDir)
			if err != nil {
				errChan <- err
			}

			c := 0
			for i, _ := range buf {
				if fmt.Sprintf("%x", buf[i]) == "5c" && fmt.Sprintf("%x", buf[i+1]) == "ff" {
					c = i + 1
					break
				}
			}

			if err := createAAC(fmt.Sprintf("%s/%s", tempdir, fname), buf[c:]); err != nil {
				errChan <- err
			}
		}(f.Name(), tempdir)
	}
	select {
	case err := <-errChan:
		return err
	default:
	}
	wg.Wait()

	return nil
}

func createAAC(name string, bf []byte) error {
	wf, err := os.Create(name)
	if err != nil {
		return err
	}
	defer wf.Close()

	wf.Write(bf)

	return nil
}

func AACtoByte(fname, aacDir string) ([]byte, error) {
	fpath := filepath.Join(aacDir, fname)
	file, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}

	f, err := file.Stat()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, f.Size())
	_, err = file.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
