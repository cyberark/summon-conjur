package conjurapi

import (
	"fmt"
	"net/http"
	"encoding/json"
	"time"
	"net/url"
	"strings"
	"io/ioutil"
)

func (c *Client) RefreshToken() (error) {
	var (
		token AuthnToken
		err error
	)

	if c.NeedsTokenRefresh() {
		var tokenBytes []byte
		tokenBytes, err = c.authenticator.RefreshToken()
		if err == nil {
			if err = json.Unmarshal(tokenBytes, &token); err == nil && token.Key != "" {
				c.authToken = token
			}
		}
	}

	return err
}

func (c *Client) NeedsTokenRefresh() bool {
	return &c.authToken == nil || !c.authToken.ValidAtTime(time.Now()) || c.authenticator.NeedsTokenRefresh()
}

func (c *Client) createAuthRequest(req *http.Request) (error) {
	if err := c.RefreshToken(); err != nil {
		return err
	}

	req.Header.Set(
		"Authorization",
		fmt.Sprintf("Token token=\"%s\"", c.authToken.Base64encoded),
	)

	return nil
}

func (c *Client) Authenticate(loginPair LoginPair) ([]byte, error) {
	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/authn/%s/%s/authenticate", c.config.BaseURL(), c.config.Account, url.QueryEscape(loginPair.Login)),
		"text/plain",
		strings.NewReader(loginPair.APIKey),
	)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case 200:
		defer resp.Body.Close()

		var tokenPayload []byte
		tokenPayload, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return tokenPayload, err
	default:
		return nil, fmt.Errorf("%v: %s", resp.StatusCode, resp.Status)
	}
}
