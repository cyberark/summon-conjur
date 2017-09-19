package wrapper_v4

import (
	"net/url"
	"fmt"
	"net/http"
	"github.com/cyberark/conjur-api-go/conjurapi/wrapper"
)

func RetrieveSecretRequest(applianceURL, variableIdentifier string) (*http.Request, error) {
	return http.NewRequest(
		"GET",
		VariableURL(applianceURL, variableIdentifier),
		nil,
	)
}

func RetrieveSecretResponse(variableIdentifier string, resp *http.Response) ([]byte, error) {
	return wrapper.RetrieveSecretResponse(variableIdentifier, resp)
}

func VariableURL(applianceURL, variableIdentifier string) string {
	return fmt.Sprintf("%s/variables/%s/value", applianceURL, url.QueryEscape(variableIdentifier))
}
