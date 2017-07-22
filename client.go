package radigo

import (
	"context"
	"fmt"
	"os"

	"github.com/yyoshiki41/go-radiko"
)

const (
	envRadikoMail     = "RADIKO_MAIL"
	envRadikoPassword = "RADIKO_PASSWORD"
)

func getClient(ctx context.Context, areaID string) (*radiko.Client, error) {
	var client *radiko.Client
	var err error

	client, err = radiko.New("")
	if err != nil {
		return nil, err
	}

	// When a currentAreaID is not equal to the given areaID,
	// it is neccessary to use the area free as the premium member.
	if areaID != "" && areaID != currentAreaID {
		mail := os.Getenv(envRadikoMail)
		password := os.Getenv(envRadikoPassword)

		login, err := client.Login(ctx, mail, password)
		switch {
		case err != nil:
			return nil, err
		case login.StatusCode() != 200:
			return nil, fmt.Errorf(
				"invalid login status code: %d", login.StatusCode())
		default:
			client.SetAreaID(areaID)
		}
	}

	return client, nil
}

func downloadSwfPlayer() error {
	_, err := os.Stat(swfPlayer)
	if !(os.IsExist(err) || os.IsNotExist(err)) {
		return err
	}

	if os.IsExist(err) {
		if err := os.Remove(swfPlayer); err != nil {
			return err
		}
	}
	return radiko.DownloadPlayer(swfPlayer)
}
