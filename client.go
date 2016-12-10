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

func getClient(ctx context.Context, authToken, areaID string) (*radiko.Client, error) {
	var client *radiko.Client
	var err error

	switch {
	case areaID != "" && areaID != currentAreaID:
		// When a currentAreaID is not equal to the given areaID,
		// it is neccessary to use the area free as the premium member.
		client, err = newClientPremiumMember(ctx, authToken)
		if err != nil {
			return client, err
		}
		client.SetAreaID(areaID)
	default:
		client, err = radiko.New("")
	}

	if err != nil {
		return nil, err
	}
	return client, nil
}

func newClientPremiumMember(ctx context.Context, authToken string) (*radiko.Client, error) {
	client, err := radiko.New(authToken)
	if err != nil {
		return nil, err
	}

	mail := os.Getenv(envRadikoMail)
	password := os.Getenv(envRadikoPassword)

	login, err := client.Login(ctx, mail, password)
	switch {
	case err != nil:
		return nil, err
	case login.StatusCode() != 200:
		return nil, fmt.Errorf(
			"invalid login status code: %d", login.StatusCode())
	}

	return client, nil
}

func downloadSwfPlayer(flagForce bool) error {
	_, err := os.Stat(swfPlayer)
	if flagForce && os.IsExist(err) {
		os.Remove(swfPlayer)
	}

	if flagForce || os.IsNotExist(err) {
		err := radiko.DownloadPlayer(swfPlayer)
		if err != nil {
			return err
		}
	}
	return nil
}
