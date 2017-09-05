package wrapper

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"github.com/cyberark/conjur-api-go/conjurapi/authn"
)

func SetRequestAuthorization(req *http.Request, base64EncodedAuthToken string) {
	req.Header.Set(
		"Authorization",
		fmt.Sprintf("Token token=\"%s\"", base64EncodedAuthToken),
	)
}

func AuthenticateRequest(applianceURL, account string, loginPair authn.LoginPair) (*http.Request, error) {
	authenticateUrl := fmt.Sprintf("%s/authn/%s/%s/authenticate", applianceURL, account, url.QueryEscape(loginPair.Login))

	req, err := http.NewRequest("POST", authenticateUrl, strings.NewReader(loginPair.APIKey))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/plain")

	return req, nil
}

func AuthenticateResponse(resp *http.Response) ([]byte, error) {
	switch resp.StatusCode {
	case 200:
		return ByteResponseTransformer(resp)
	default:
		return nil, fmt.Errorf("Authentication: %s", resp.Status)
	}
}
