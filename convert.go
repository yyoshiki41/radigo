package radigo

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"sort"
)

type sliceFileInfo []os.FileInfo

func (f sliceFileInfo) Len() int           { return len(f) }
func (f sliceFileInfo) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f sliceFileInfo) Less(i, j int) bool { return f[i].Name() < f[j].Name() }

func concatFileNames() (string, error) {
	files, err := ioutil.ReadDir(aacPath)
	if err != nil {
		return "", nil
	}
	sort.Sort(sliceFileInfo(files))

	var res []byte
	for _, f := range files {
		res = append(res, f.Name()...)
		res = append(res, '|')
	}

	// remove the last element "|"
	return string(res[:len(res)-1]), nil
}

type ffmpeg struct {
	*exec.Cmd
}

func newFfmpeg() (*ffmpeg, error) {
	cmdPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return nil, err
	}

	return &ffmpeg{&exec.Cmd{Path: cmdPath}}, nil
}

func (f *ffmpeg) setDir(dir string) {
	f.Dir = dir
}

func (f *ffmpeg) setArgs(args ...string) {
	f.Args = append([]string{f.Path}, args...)
}

func createConcatedAACFile() error {
	name, err := concatFileNames()
	if err != nil {
		return err
	}

	f, err := newFfmpeg()
	if err != nil {
		return err
	}
	f.setDir(aacPath)
	f.setArgs("-i", "concat:"+name, "-c", "copy", path.Join(radigoPath, "result.aac"))
	return f.Run()
}

func convertAACToMP3() error {
	f, err := newFfmpeg()
	if err != nil {
		return err
	}
	f.setDir(radigoPath)
	f.setArgs("-i", path.Join(radigoPath, "result.aac"), "-c:a", "libmp3lame", "-ac", "2", "-q:a", "2", path.Join(radigoPath, "result.mp3"))
	return f.Run()
}
