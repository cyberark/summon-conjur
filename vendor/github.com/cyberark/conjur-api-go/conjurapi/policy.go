package conjurapi

import (
	"io"
	"github.com/cyberark/conjur-api-go/conjurapi/wrapper"
)

func (c *Client) LoadPolicy(policyIdentifier string, policy io.Reader) ([]byte, error) {

	req, err := wrapper.LoadPolicyRequest(c.config.ApplianceURL, c.config.Account, policyIdentifier, policy)
	if err != nil {
		return nil, err
	}

	err = c.createAuthRequest(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return wrapper.LoadPolicyResponse(resp)
}
