package conjurapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"io"
)

func (c *Client) LoadPolicy(policyIdentifier string, policy io.Reader) (string, error) {
	policyUrl := fmt.Sprintf("%s/policies/%s/policy/%s", c.config.BaseURL(), c.config.Account, policyIdentifier)
	req, err := http.NewRequest(
		"PUT",
		policyUrl,
		policy,
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
	case 201:
		defer resp.Body.Close()

		var responseText []byte
		responseText, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		return string(responseText), err
	default:
		return "", fmt.Errorf("%v: %s", resp.StatusCode, resp.Status)
	}
}
