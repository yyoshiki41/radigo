package radigo

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	defaultOutputDir = "output"
	envRadigoHome    = "RADIGO_HOME"
)

// OutputConfig contains the configuration for output files.
type OutputConfig struct {
	DirFullPath  string
	FileBaseName string // base name of the file
	FileFormat   string // aac, mp3
}

func NewOutputConfig(fileBaseName, fileFormat string) (*OutputConfig, error) {
	// If the environment variable RADIGO_HOME is set,
	// override working directory path.
	fullPath := os.Getenv(envRadigoHome)
	switch {
	case fullPath != "" && !filepath.IsAbs(fullPath):
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		fullPath = filepath.Join(wd, fullPath)
	case fullPath == "":
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		fullPath = filepath.Join(wd, defaultOutputDir)
	default:
	}

	return &OutputConfig{
		DirFullPath:  filepath.Clean(fullPath),
		FileBaseName: fileBaseName,
		FileFormat:   fileFormat,
	}, nil
}

// SetupDir configures the output directory or returns an error if failed to create it.
func (c *OutputConfig) SetupDir() error {
	_, err := os.Stat(c.DirFullPath)
	switch {
	case err == nil:
		// Output directory already exists.
	case os.IsNotExist(err):
		// Output directory does not exist.
		if err := os.MkdirAll(c.DirFullPath, 0755); err != nil {
			return err
		}
	default:
		return err
	}
	return nil
}

func (c *OutputConfig) TempAACDir() (string, error) {
	aacDir, err := ioutil.TempDir(c.DirFullPath, "aac")
	if err != nil {
		return "", err
	}

	return aacDir, nil
}

func (c *OutputConfig) AudioFormat() string {
	return c.FileFormat
}

func (c *OutputConfig) AbsPath() string {
	name := fmt.Sprintf("%s.%s", c.FileBaseName, c.FileFormat)
	return filepath.Join(c.DirFullPath, name)
}

func concatAACFiles(ctx context.Context, resourcesDir string) (string, error) {
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
	if err := ConcatAACFilesFromList(ctx, listFile.Name(), concatedFile); err != nil {
		return "", err
	}

	return concatedFile, nil
}
