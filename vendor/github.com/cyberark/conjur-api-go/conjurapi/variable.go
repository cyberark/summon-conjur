package conjurapi

import (
	"github.com/cyberark/conjur-api-go/conjurapi/wrapper"
	"net/http"
	"github.com/cyberark/conjur-api-go/conjurapi/wrapper_v4"
)

func (c *Client) RetrieveSecret(variableIdentifier string) ([]byte, error) {
	var (
		req *http.Request
		err error
	)

	if c.config.V4 {
		req, err = wrapper_v4.RetrieveSecretRequest(c.config.ApplianceURL, variableIdentifier)
	} else {
		req, err = wrapper.RetrieveSecretRequest(c.config.ApplianceURL, c.config.Account, variableIdentifier)
	}

	if err != nil {
		return nil, err
	}

	resp, err := c.SubmitRequest(req)
	if err != nil {
		return nil, err
	}

	if c.config.V4 {
		return wrapper_v4.RetrieveSecretResponse(variableIdentifier, resp)
	} else {
		return wrapper.RetrieveSecretResponse(variableIdentifier, resp)
	}
}

func (c *Client) AddSecret(variableIdentifier string, secretValue string) error {
	req, err := wrapper.AddSecretRequest(c.config.ApplianceURL, c.config.Account, variableIdentifier, secretValue)
	if err != nil {
		return err
	}

	resp, err := c.SubmitRequest(req)
	if err != nil {
		return err
	}

	return wrapper.AddSecretResponse(variableIdentifier, resp)
}
