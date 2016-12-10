package radigo

import (
	"context"
	"os/exec"
)

type ffmpeg struct {
	*exec.Cmd
}

func newFfmpeg(ctx context.Context, input string) (*ffmpeg, error) {
	cmdPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return nil, err
	}

	return &ffmpeg{exec.CommandContext(
		ctx,
		cmdPath,
		"-i", input,
	)}, nil
}

func (f *ffmpeg) setDir(dir string) {
	f.Dir = dir
}

func (f *ffmpeg) setArgs(args ...string) {
	f.Args = append(f.Args, args...)
}

func (f *ffmpeg) run(output string) error {
	f.Args = append(f.Args, output)
	return f.Run()
}

func (f *ffmpeg) start(output string) error {
	f.Args = append(f.Args, output)
	return f.Start()
}
