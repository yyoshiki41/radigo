package radigo

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
)

func bulkDownload(list []string) error {
	const maxAttempts = 5

	var errFlag bool
	var wg sync.WaitGroup
	for _, v := range list {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()

			var err error
			for i := 0; i < maxAttempts; i++ {
				err = download(link)
				if err == nil {
					break
				}
			}
			if err != nil {
				log.Printf("Failed to download aac file: %s", err)
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

func download(link string) error {
	resp, err := http.Get(link)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, fileName := path.Split(link)
	file, err := os.Create(path.Join(aacPath, fileName))
	if err != nil {
		return err
	}

	_, err = io.Copy(file, resp.Body)
	if closeErr := file.Close(); err == nil {
		err = closeErr
	}
	return err
}
