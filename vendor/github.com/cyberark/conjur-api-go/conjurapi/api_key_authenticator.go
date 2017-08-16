package conjurapi

import (
	"strings"
	"io/ioutil"
	"fmt"
	"net/http"
	"time"
	"net/url"
)

type APIKeyAuthenticator struct {
	AuthnURLTemplate string
	Login string `env:"CONJUR_AUTHN_LOGIN"`
	APIKey string `env:"CONJUR_AUTHN_API_KEY"`
}

func (a *APIKeyAuthenticator) RefreshToken() ([]byte, error) {
	httpclient := http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := httpclient.Post(
		fmt.Sprintf(a.AuthnURLTemplate, url.QueryEscape(a.Login)),
		"text/plain",
		strings.NewReader(a.APIKey),
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
		return nil, fmt.Errorf("%v: %s\n", resp.StatusCode, resp.Status)
	}
}

func (a *APIKeyAuthenticator) NeedsTokenRefresh() bool {
	return false
}
