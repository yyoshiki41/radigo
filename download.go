package radigo

import (
	"io"
	"net/http"
	"os"
	"path"
	"sync"
)

const maxAttempts = 5

func bulkDownload(list []string) error {
	var wg sync.WaitGroup

	for _, v := range list {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()

			for i = 0; i < maxAttempts; i++ {
				err := download(link)
				if err == nil {
					break
				}
			}
		}(v)
	}
	wg.Wait()

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
