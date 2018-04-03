# Go Radiko

[![godoc](https://godoc.org/github.com/yyoshiki41/go-radiko?status.svg)](https://godoc.org/github.com/yyoshiki41/go-radiko)
[![build](https://travis-ci.org/yyoshiki41/go-radiko.svg?branch=master)](https://travis-ci.org/yyoshiki41/go-radiko)
[![codecov](https://codecov.io/gh/yyoshiki41/go-radiko/branch/master/graph/badge.svg)](https://codecov.io/gh/yyoshiki41/go-radiko)
[![go report](https://goreportcard.com/badge/github.com/yyoshiki41/go-radiko)](https://goreportcard.com/report/github.com/yyoshiki41/go-radiko)

The __unofficial__ [radiko.jp](https://radiko.jp/) APIs Client Library for Go

## Installation


- Go 1.7 or newer

```bash
$ go get github.com/yyoshiki41/go-radiko
```

## Usage

### ■ Default

```go
// authentication token is not necessary.
client, err := radiko.New("")
if err != nil {
	panic(err)
}
// Get programs data
stations, err := client.GetNowPrograms(context.Background())
if err != nil {
	log.Fatal(err)
}
fmt.Printf("%v", stations)
```

### ■ Get & Set authentication token

```go
// 1. Create a new Client.
client, err := radiko.New("")
if err != nil {
	panic(err)
}

// 2. Enables and sets the auth_token.
// After client.AuthorizeToken() has succeeded,
// the client has the enabled auth_token internally.
authToken, err := client.AuthorizeToken(context.Background(), authKeyPath)
if err != nil {
	log.Fatal(err)
}
```

#### Premium member (Enable to use the [area free](http://radiko.jp/rg/premium/).)

```go
// 1. Create a new Client.
client, err := radiko.New("")
if err != nil {
	panic(err)
}

// 2. Login as the premium member
// After client.Login() has succeeded,
// the client has the valid cookie internally.
ctx := context.Background()
login, err := client.Login(ctx, "example@mail.com", "example_password")
if err != nil {
	log.Fatal(err)
}
if login.StatusCode() != 200 {
	log.Fatalf("Failed to login premium member.\nInvalid status code: %d",
		login.StatusCode())
}

// 3. Enables and sets the auth_token.
// After client.AuthorizeToken() has succeeded,
// the client has the enabled auth_token internally.
authToken, err := client.AuthorizeToken(context.Background(), authKeyPath)
if err != nil {
	log.Fatal(err)
}
```


### ■ Use your authentication token

```go
// If the auth_token is cached, set your token in HTTP Header like below.
client, err = radiko.New("auth_token")
if err != nil {
	panic(err)
}
```

## Examples

It is possible to try [examples](https://github.com/yyoshiki41/go-radiko/tree/master/examples).

```bash
# Get programs data
$ go run ./examples/main.go
# Get & Set auth_token
$ go run ./examples/auth/main.go
```

## Projects using go-radiko

- [yyoshiki41/radigo](https://github.com/yyoshiki41/radigo) - Record a radiko program.
- [juntaki/toyotomimi](https://github.com/juntaki/toyotomimi)

## License 
The MIT License

## Author

Yoshiki Nakagawa
