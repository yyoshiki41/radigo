package radigo

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

func bulkDownload(list []string) error {
	var wg sync.WaitGroup

	for _, v := range list {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()

			err := download(link)
			if err != nil {
				fmt.Println(err)
			}
		}(v)
	}
	wg.Wait()

	return nil
}

func download(link string) error {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: time.Duration(300 * time.Second)}
	resp, err := client.Do(req)
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
