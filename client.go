package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type ConjurClient struct {
	AuthnUrl    string
	CoreUrl     string
	Username    string
	APIKey      string
	SSLCertPath string
	httpClient  *http.Client
}

func NewConjurClient() (*ConjurClient, error) {

	config, err := LoadConfig()

	if err != nil {
		return nil, err
	}

	httpClient, err := NewConjurHTTPClient(config.SSLCertPath)

	if err != nil {
		return nil, err
	}

	return &ConjurClient{
		AuthnUrl:    config.AuthnUrl(),
		CoreUrl:     config.CoreUrl(),
		Username:    config.Username,
		APIKey:      config.APIKey,
		SSLCertPath: config.SSLCertPath,
		httpClient:  httpClient,
	}, nil
}

func (c *ConjurClient) getAuthToken() (string, error) {
	escapedUsername := url.QueryEscape(c.Username)
	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/users/%s/authenticate", c.AuthnUrl, escapedUsername),
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
	variableUrl := fmt.Sprintf("%v/variables/%v/value", c.CoreUrl, variable)

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
	switch resp.StatusCode {
	case 404:
		return "", fmt.Errorf("%v Variable '%s' not found", resp.StatusCode, path)
	case 403:
		return "", fmt.Errorf("%s Invalid permissions on '%s'", resp.StatusCode, path)
	case 200:
		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return "", err
		}
		return string(body), nil
	default:
		return "", fmt.Errorf("%s %s", resp.StatusCode, resp.Status)
	}
}
