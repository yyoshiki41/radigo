package main

import (
	"context"
	"fmt"
	"log"

	radiko "github.com/yyoshiki41/go-radiko"
)

func main() {
	// 1. Create a new Client.
	client, err := radiko.New("")
	if err != nil {
		log.Fatalf("Failed to construct a radiko Client. %s", err)
	}

	// 2. Enables and sets the auth_token.
	// After client.AuthorizeToken() has succeeded,
	// the client has the enabled auth_token internally.
	authToken, err := client.AuthorizeToken(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(authToken)
}
