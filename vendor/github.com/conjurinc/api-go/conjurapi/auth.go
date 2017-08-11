package conjurapi

import (
	"fmt"
	"net/url"
	"strings"
	"io/ioutil"
	"net/http"
	"encoding/base64"
)

func (c *Client) getAuthToken() (string, error) {
	authUrl := fmt.Sprintf("%s/authn/%s/%s/authenticate", c.config.ApplianceUrl, c.config.Account, url.QueryEscape(c.config.Username))
	resp, err := c.httpClient.Post(
		authUrl,
		"text/plain",
		strings.NewReader(c.config.APIKey),
	)
	if err != nil {
		return "", err
	}

	switch resp.StatusCode {
	case 200:
		defer resp.Body.Close()

		var tokenPayload []byte
		tokenPayload, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		return base64.StdEncoding.EncodeToString(tokenPayload), err
	default:
		return "", fmt.Errorf("%v: %s\n", resp.StatusCode, resp.Status)
	}
}

func (c *Client) createAuthRequest(req *http.Request) (error) {
	token, err := c.getAuthToken()
	if err != nil {
		return err
	}

	req.Header.Set(
		"Authorization",
		fmt.Sprintf("Token token=\"%s\"", token),
	)

	return nil
}
