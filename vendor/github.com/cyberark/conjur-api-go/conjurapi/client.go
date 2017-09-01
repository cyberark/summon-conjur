package conjurapi

import (
	"net/http"
	"time"
	"fmt"
	"crypto/x509"
	"crypto/tls"
	"os"
	"github.com/bgentry/go-netrc/netrc"
)

type Authenticator interface {
	RefreshToken() ([]byte, error)
	NeedsTokenRefresh() bool
}

type Client struct {
	config        Config
	authToken     AuthnToken
	httpClient    *http.Client
	authenticator Authenticator
}

func NewClientFromKey(config Config, loginPair LoginPair) (*Client, error) {
	authenticator := &APIKeyAuthenticator{
		LoginPair: loginPair,
	}
	client, err := newClientWithAuthenticator(
		config,
		authenticator,
	)
	authenticator.Authenticate = client.Authenticate
	return client, err
}

func NewClientFromTokenFile(config Config, tokenFile string) (*Client, error) {
	return newClientWithAuthenticator(
		config,
		&TokenFileAuthenticator{
			TokenFile: tokenFile,
			MaxWaitTime: time.Second*10,
		},
	)
}

func LoginPairFromEnv() (*LoginPair, error) {
	return &LoginPair{
		Login: os.Getenv("CONJUR_AUTHN_LOGIN"),
		APIKey: os.Getenv("CONJUR_AUTHN_API_KEY"),
	}, nil
}

func LoginPairFromNetRC(config Config) (*LoginPair, error) {
	if config.NetRCPath == "" {
		config.NetRCPath = os.ExpandEnv("$HOME/.netrc")
	}

	rc, err := netrc.ParseFile(config.NetRCPath)
	if err != nil {
		return nil, err
	}

	m := rc.FindMachine(config.ApplianceURL + "/authn")

	if m == nil {
		return nil, fmt.Errorf("No credentials found in NetRCPath")
	}

	return &LoginPair{Login: m.Login, APIKey: m.Password}, nil
}

func NewClientFromEnvironment(config Config) (*Client, error) {
	err := config.validate()

	if err != nil {
		return nil, err
	}

	authnTokenFile := os.Getenv("CONJUR_AUTHN_TOKEN_FILE")
	if authnTokenFile != ""  {
		return NewClientFromTokenFile(config, authnTokenFile)
	}

	loginPair, err := LoginPairFromEnv()
	if err == nil && loginPair.Login != "" && loginPair.APIKey != ""  {
		return NewClientFromKey(config, *loginPair)
	}

	loginPair, err = LoginPairFromNetRC(config)
	if err == nil && loginPair.Login != "" && loginPair.APIKey != ""  {
		return NewClientFromKey(config, *loginPair)
	}

	return nil, fmt.Errorf("Environment variables and machine identity files satisfying at least one authentication strategy must be present!")
}

func newClientWithAuthenticator(config Config, authenticator Authenticator) (*Client, error) {
	var (
		err error
	)

	err = config.validate()

	if err != nil {
		return nil, err
	}

	var httpClient *http.Client

	if config.Https {
		cert, err := config.ReadSSLCert()
		if err != nil {
			return nil, err
		}
		httpClient, err = newHTTPSClient(cert)
		if err != nil {
			return nil, err
		}
	} else {
		httpClient = &http.Client{Timeout: time.Second * 10}
	}

	return &Client{
		config:        config,
		authenticator: authenticator,
		httpClient: httpClient,
	}, nil
}

func newHTTPSClient(cert []byte) (*http.Client, error) {
	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM(cert)
	if !ok {
		return nil, fmt.Errorf("Can't append Conjur SSL cert")
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: pool},
	}
	return &http.Client{Transport: tr, Timeout: time.Second * 10}, nil
}
