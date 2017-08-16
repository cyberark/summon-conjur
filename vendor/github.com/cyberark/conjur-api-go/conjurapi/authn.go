package conjurapi

import (
	"fmt"
	"net/http"
	"encoding/json"
	"time"
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
				c.authToken = &token
			}
		}
	}

	return err
}

func (c *Client) NeedsTokenRefresh() bool {
	return c.authToken == nil || !c.authToken.ValidAtTime(time.Now()) || c.authenticator.NeedsTokenRefresh()
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
