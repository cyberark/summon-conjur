package conjurapi

import (
	"net/url"
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
)

func (c *Client) generateVariableUrl(varId string) string {
	escapedVarId := url.QueryEscape(varId)
	return fmt.Sprintf("%s/secrets/%s/variable/%s", c.config.BaseURL(), c.config.Account, escapedVarId)
}

func (c *Client) RetrieveSecret(variableIdentifier string) (string, error) {
	variableUrl := c.generateVariableUrl(variableIdentifier)
	req, err := http.NewRequest(
		"GET",
		variableUrl,
		nil,
	)
	if err != nil {
		return "", err
	}

	err = c.createAuthRequest(req)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	switch resp.StatusCode {
	case 404:
		return "", fmt.Errorf("%v: Variable '%s' not found", resp.StatusCode, variableIdentifier)
	case 403:
		return "", fmt.Errorf("%v: Invalid permissions on '%s'", resp.StatusCode, variableIdentifier)
	case 200:
		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return "", err
		}
		return string(body), nil
	default:
		return "", fmt.Errorf("%v: %s", resp.StatusCode, resp.Status)
	}
}

func (c *Client) AddSecret(variableIdentifier string, secretValue string) (error) {
	variableUrl := c.generateVariableUrl(variableIdentifier)
	req, err := http.NewRequest(
		"POST",
		variableUrl,
		strings.NewReader(secretValue),
	)
	if err != nil {
		return err
	}

	err = c.createAuthRequest(req)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case 404:
		return fmt.Errorf("%v: Variable '%s' not found", resp.StatusCode, variableIdentifier)
	case 403:
		return fmt.Errorf("%v: Invalid permissions on '%s'", resp.StatusCode, variableIdentifier)
	case 201:
		return nil
	default:
		return fmt.Errorf("%v: %s", resp.StatusCode, resp.Status)
	}
}
