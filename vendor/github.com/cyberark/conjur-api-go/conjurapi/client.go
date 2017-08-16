package conjurapi

import (
	"net/http"
	"time"
	"fmt"
)

type Authenticator interface {
	RefreshToken() ([]byte, error)
	NeedsTokenRefresh() bool
}

type Client struct {
	config     Config
	authToken  *AuthnToken
	httpclient *http.Client
	authenticator Authenticator
}

func AuthnURL(applianceURL, account string) string {
	return fmt.Sprintf("%s/authn/%s/%s/authenticate", applianceURL, account, "%s")
}

func NewClientFromKey(config Config, login string, aPIKey string) (*Client, error) {
	return newClientWithAuthenticator(
		config,
		&APIKeyAuthenticator{
			AuthnURLTemplate: AuthnURL(config.ApplianceURL, config.Account),
			Login:            login,
			APIKey:           aPIKey,
		},
	)
}

func NewClientFromTokenFile(config Config, tokenFile string) (*Client, error) {
	return newClientWithAuthenticator(
		config,
		&TokenFileAuthenticator{
			TokenFile: tokenFile,
		},
	)
}

func newClientWithAuthenticator(config Config, authenticator Authenticator) (*Client, error) {
	var (
		err error
	)

	err = config.validate()

	if err != nil {
		return nil, err
	}

	return &Client{
		config:     config,
		authenticator:  authenticator,
		httpclient: &http.Client{
			Timeout: time.Second * 10,
		},
	}, nil
}

