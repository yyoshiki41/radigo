package radigo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
)

type sliceFileInfo []os.FileInfo

func (f sliceFileInfo) Len() int           { return len(f) }
func (f sliceFileInfo) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f sliceFileInfo) Less(i, j int) bool { return f[i].Name() < f[j].Name() }

func concatFileNames(aacDir string) (string, error) {
	files, err := ioutil.ReadDir(aacDir)
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

func createConcatedAACFile(aacDir string) error {
	name, err := concatFileNames(aacDir)
	if err != nil {
		return err
	}

	f, err := newFfmpeg(fmt.Sprintf("concat:%s", name))
	if err != nil {
		return err
	}

	f.setDir(aacDir)
	f.setArgs("-c", "copy")
	// TODO: Run 結果の標準出力を拾う
	return f.run(path.Join(radigoPath, "result.aac"))
}

func convertAACToMP3() error {
	f, err := newFfmpeg(path.Join(radigoPath, "result.aac"))
	if err != nil {
		return err
	}

	f.setDir(radigoPath)
	f.setArgs(
		"-c:a", "libmp3lame",
		"-ac", "2",
		"-q:a", "2",
	)
	// TODO: Run 結果の標準出力を拾う
	return f.run(path.Join(radigoPath, "result.mp3"))
}
