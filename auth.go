package radigo

import (
	"io/ioutil"
	"os"
	"path"
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
