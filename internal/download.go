package internal

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

const (
	maxAttempts    = 4
	maxConcurrents = 64
)

var sem = make(chan struct{}, maxConcurrents)

func BulkDownload(list []string, output string) error {
	var errFlag bool
	var wg sync.WaitGroup

	for _, v := range list {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()

			var err error
			for i := 0; i < maxAttempts; i++ {
				sem <- struct{}{}
				err = download(link, output)
				<-sem
				if err == nil {
					break
				}
			}
			if err != nil {
				log.Printf("Failed to download: %s", err)
				errFlag = true
			}
		}(v)
	}
	wg.Wait()

	if errFlag {
		return errors.New("Lack of aac files")
	}
	return nil
}

func download(link, output string) error {
	resp, err := http.Get(link)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, fileName := filepath.Split(link)
	file, err := os.Create(filepath.Join(output, fileName))
	if err != nil {
		return err
	}

	_, err = io.Copy(file, resp.Body)
	if closeErr := file.Close(); err == nil {
		err = closeErr
	}
	return err
}
