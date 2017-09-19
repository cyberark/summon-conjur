package wrapper_v4

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"github.com/cyberark/conjur-api-go/conjurapi/authn"
	"github.com/cyberark/conjur-api-go/conjurapi/wrapper"
)

func AuthenticateRequest(applianceURL string, loginPair authn.LoginPair) (*http.Request, error) {
	authenticateUrl := fmt.Sprintf("%s/authn/users/%s/authenticate", applianceURL, url.QueryEscape(loginPair.Login))

	req, err := http.NewRequest("POST", authenticateUrl, strings.NewReader(loginPair.APIKey))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/plain")

	return req, nil
}

func AuthenticateResponse(resp *http.Response) ([]byte, error) {
	return wrapper.AuthenticateResponse(resp)
}
