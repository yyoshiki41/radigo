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

func newClientPremiumMember(authToken string) (*radiko.Client, error) {
	client, err := radiko.New(authToken)
	if err != nil {
		return nil, err
	}

	mail := os.Getenv(envRadikoMail)
	password := os.Getenv(envRadikoPassword)
	login, err := client.Login(context.Background(), mail, password)

	switch {
	case err != nil:
		return nil, err
	case login.StatusCode() != 200:
		return nil, fmt.Errorf(
			"invalid login status code: %d", login.StatusCode())
	}

	return client, nil
}
