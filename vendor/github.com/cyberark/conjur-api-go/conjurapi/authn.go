package conjurapi

import (
	"net/http"
	"encoding/json"
	"time"
	"github.com/cyberark/conjur-api-go/conjurapi/authn"
	"github.com/cyberark/conjur-api-go/conjurapi/wrapper"
	"github.com/cyberark/conjur-api-go/conjurapi/wrapper_v4"
)

func (c *Client) RefreshToken() (error) {
	var (
		token authn.AuthnToken
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

	wrapper.SetRequestAuthorization(req, c.authToken.Base64encoded)

	return nil
}

func (c *Client) Authenticate(loginPair authn.LoginPair) ([]byte, error) {
	req, err := wrapper.AuthenticateRequest(c.config.ApplianceURL, c.config.Account, loginPair)
	if c.config.V4 {
		req, err = wrapper_v4.AuthenticateRequest(c.config.ApplianceURL, loginPair)
	} else {
		req, err = wrapper.AuthenticateRequest(c.config.ApplianceURL, c.config.Account, loginPair)
	}

	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}


	if c.config.V4 {
		return wrapper_v4.AuthenticateResponse(resp)
	} else {
		return wrapper.AuthenticateResponse(resp)
	}
}
