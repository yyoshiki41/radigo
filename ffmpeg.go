package radigo

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

type ffmpeg struct {
	*exec.Cmd
}

func newFfmpeg(ctx context.Context) (*ffmpeg, error) {
	cmdPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return nil, err
	}

	return &ffmpeg{exec.CommandContext(
		ctx,
		cmdPath,
	)}, nil
}

func (f *ffmpeg) setDir(dir string) {
	f.Dir = dir
}

func (f *ffmpeg) setArgs(args ...string) {
	f.Args = append(f.Args, args...)
}

func (f *ffmpeg) setInput(input string) {
	f.setArgs("-i", input)
}

func (f *ffmpeg) run(output string) error {
	f.setArgs(output)
	return f.Run()
}

func (f *ffmpeg) start(output string) error {
	f.setArgs(output)
	return f.Start()
}

func (f *ffmpeg) wait() error {
	return f.Wait()
}

func (f *ffmpeg) stdinPipe() (io.WriteCloser, error) {
	return f.StdinPipe()
}

func (f *ffmpeg) stderrPipe() (io.ReadCloser, error) {
	return f.StderrPipe()
}

// ConvertAACtoMP3 converts an aac file to a mp3 file.
func ConvertAACtoMP3(ctx context.Context, input, output string) error {
	f, err := newFfmpeg(ctx)
	if err != nil {
		return err
	}

	f.setInput(input)
	f.setArgs(
		"-c:a", "libmp3lame",
		"-ac", "2",
		"-q:a", "2",
		"-y", // overwrite the output file without asking
	)
	// TODO: Collect log
	return f.run(output)
}

// ConcatAACFilesFromList concatenates files from the list of resources.
func ConcatAACFilesFromList(ctx context.Context, resourcesDir string) (string, error) {
	files, err := ioutil.ReadDir(resourcesDir)
	if err != nil {
		return "", err
	}

	listFile, err := ioutil.TempFile(resourcesDir, "aac_resources")
	if err != nil {
		return "", err
	}
	defer os.Remove(listFile.Name())

	for _, f := range files {
		p := fmt.Sprintf("file '%s'\n", filepath.Join(resourcesDir, f.Name()))
		if _, err := listFile.WriteString(p); err != nil {
			return "", err
		}
	}

	concatedFile := filepath.Join(resourcesDir, "concated.aac")
	if err := ConcatAACFiles(ctx, listFile.Name(), concatedFile); err != nil {
		return "", err
	}

	return concatedFile, nil
}

// ConcatAACFiles concatenate files of the same type.
func ConcatAACFiles(ctx context.Context, input, output string) error {
	f, err := newFfmpeg(ctx)
	if err != nil {
		return err
	}

	f.setArgs(
		"-f", "concat",
		"-safe", "0",
	)
	f.setInput(input)
	f.setArgs("-c", "copy")
	// TODO: Collect log
	return f.run(output)
}
