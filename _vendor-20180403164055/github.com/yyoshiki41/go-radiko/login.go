package radiko

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
)

// Login returns the Statuser that has StatusCode method.
func (c *Client) Login(ctx context.Context, mail, password string) (Statuser, error) {
	err := c.login(ctx, mail, password)
	if err != nil {
		return nil, err
	}

	return c.loginCheck(ctx)
}

func (c *Client) login(ctx context.Context, mail, password string) error {
	apiEndpoint := "ap/member/login/login"
	v := url.Values{}
	v.Set("mail", mail)
	v.Set("pass", password)

	req, err := c.newRequest(ctx, "POST", apiEndpoint,
		&Params{body: strings.NewReader(v.Encode())})
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// read the response body
	_, _ = ioutil.ReadAll(resp.Body)
	return nil
}

func (c *Client) loginCheck(ctx context.Context) (Statuser, error) {
	apiEndpoint := "ap/member/webapi/member/login/check"
	req, err := c.newRequest(ctx, "GET", apiEndpoint, &Params{})
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		status := LoginOK{}
		if err := json.Unmarshal(b, &status); err != nil {
			return nil, err
		}
		return status, nil
	}

	status := LoginNG{}
	if err := json.Unmarshal(b, &status); err != nil {
		return nil, err
	}
	return status, nil
}

// Statuser is the interface that wraps StatusCode method.
type Statuser interface {
	StatusCode() int
}

// LoginStatus is a base struct that has a Status field.
type LoginStatus struct {
	Status string `json:"status"`
}

// StatusCode returns StatusCode that is type int.
func (l *LoginStatus) StatusCode() int {
	i, _ := strconv.Atoi(l.Status)
	return i
}

// LoginOK represents login is successful.
type LoginOK struct {
	*LoginStatus
	UserKey    string `json:"user_key"`
	PaidMember string `json:"paid_member"`
	Areafree   string `json:"areafree"`
}

// LoginNG represents login failed.
type LoginNG struct {
	*LoginStatus
	Message string `json:"message"`
	Cause   string `json:"cause"`
}
