package wrapper

import (
	"fmt"
	"net/http"
	"io"
)

func LoadPolicyRequest(applianceURL string, account string, policyIdentifier string, policy io.Reader) (*http.Request, error) {
	policyUrl := fmt.Sprintf("%s/policies/%s/policy/%s", applianceURL, account, policyIdentifier)

	return http.NewRequest(
		"PUT",
		policyUrl,
		policy,
	)
}

func LoadPolicyResponse(resp *http.Response) ([]byte, error) {
	switch resp.StatusCode {
	case 201:
		return ByteResponseTransformer(resp)
	default:
		return nil, fmt.Errorf("%v: %s", resp.StatusCode, resp.Status)
	}
}

