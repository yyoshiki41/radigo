package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yyoshiki41/go-radiko"
)

func main() {
	// if an auth_token is not necessary.
	client, err := radiko.New("")
	if err != nil {
		log.Fatalf("Failed to construct a radiko Client. %s", err)
	}
	// Get stations data
	stations, err := client.GetNowPrograms(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", stations)

	// if an auth_token is cached, set a token header like below.
	client, err = radiko.New("auth_token")
	if err != nil {
		log.Fatalf("Failed to construct a radiko Client. %s", err)
	}
	_ = client
}
