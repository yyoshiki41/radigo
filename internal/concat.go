package internal

import (
	"io/ioutil"
	"os"
	"sort"
)

type sliceFileInfo []os.FileInfo

func (f sliceFileInfo) Len() int           { return len(f) }
func (f sliceFileInfo) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f sliceFileInfo) Less(i, j int) bool { return f[i].Name() < f[j].Name() }

func ConcatFileNames(output string) (string, error) {
	files, err := ioutil.ReadDir(output)
	if err != nil {
		return "", err
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
