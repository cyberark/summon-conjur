package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type ConjurClient struct {
	ApplianceUrl string
	Username     string
	APIKey       string
	SSLCertPath  string
	httpClient   *http.Client
}

func NewConjurClient() (*ConjurClient, error) {
	sslCertPath := os.Getenv("GO_CONJUR_SSL_CERTIFICATE_PATH")
	httpClient, err := NewConjurHTTPClient(sslCertPath)
	if err != nil {
		return nil, err
	}

	return &ConjurClient{
		ApplianceUrl: os.Getenv("GO_CONJUR_APPLIANCE_URL"),
		Username:     os.Getenv("GO_CONJUR_AUTHN_LOGIN"),
		APIKey:       os.Getenv("GO_CONJUR_AUTHN_API_KEY"),
		SSLCertPath:  sslCertPath,
		httpClient:   httpClient,
	}, nil
}

func (c *ConjurClient) getAuthToken() (string, error) {
	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/authn/users/%s/authenticate", c.ApplianceUrl, c.Username),
		"text/plain",
		strings.NewReader(c.APIKey),
	)

	var token []byte

	if err != nil {
		return "", err
	} else {
		defer resp.Body.Close()
		token, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
	}

	return string(token), err
}

func (c *ConjurClient) RetrieveVariable(path string) (string, error) {
	variable := url.QueryEscape(path)
	variableUrl := fmt.Sprintf("%v/variables/%v/value", c.ApplianceUrl, variable)

	token, err := c.getAuthToken()
	if err != nil {
		return "", err
	}
	tokenHeader := fmt.Sprintf("Token token=\"%s\"", base64.StdEncoding.EncodeToString([]byte(token)))

	req, err := http.NewRequest("GET", variableUrl, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set(
		"Authorization",
		tokenHeader,
	)
	// This is needed to avoid parsing the %2F in variable name
	req.URL.Opaque = variableUrl

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
