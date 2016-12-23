package radigo

import (
	"io/ioutil"
	"os"
	"path"
	"time"
)

var (
	tokenCache = path.Join(cachePath, "auth_token")
)

func getTokenCache() (string, error) {
	b, err := ioutil.ReadFile(tokenCache)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func removeTokenCache() error {
	f, err := os.Create(tokenCache)
	if err != nil {
		return err
	}
	defer f.Close()

	// overwrite token cache
	if _, err = f.Write(([]byte)("")); err != nil {
		return err
	}
	return nil
}

func isExpiredCache() bool {
	f, err := os.Stat(tokenCache)
	if err != nil {
		return false
	}

	// cache lifetime is 2 hours.
	expireationDate := f.ModTime().Add(2 * time.Hour)
	return time.Now().After(expireationDate)
}

func saveToken(authToken string) error {
	f, err := os.Create(tokenCache)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.Write(([]byte)(authToken)); err != nil {
		return err
	}
	return nil
}
