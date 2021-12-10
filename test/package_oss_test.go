package test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cyberark/conjur-api-go/conjurapi"
	conjur_authn "github.com/cyberark/conjur-api-go/conjurapi/authn"
	"github.com/stretchr/testify/assert"
)

func TestPackageOSS(t *testing.T) {
	ApplianceURL := os.Getenv("CONJUR_APPLIANCE_URL")
	Account := os.Getenv("CONJUR_ACCOUNT")
	Login := os.Getenv("CONJUR_AUTHN_LOGIN")
	APIKey := os.Getenv("CONJUR_AUTHN_API_KEY")

	Path := os.Getenv("PATH")

	t.Run("Given no configuration and no authentication information", func(t *testing.T) {
		e := ClearEnv()
		defer e.RestoreEnv()
		os.Setenv("PATH", Path)

		WithoutArgs(t)
	})

	t.Run("Given valid V5 OSS configuration", func(t *testing.T) {
		e := ClearEnv()
		defer e.RestoreEnv()
		os.Setenv("PATH", Path)

		os.Setenv("CONJUR_APPLIANCE_URL", ApplianceURL)
		os.Setenv("CONJUR_ACCOUNT", Account)

		t.Run("Given valid APIKey credentials", func(t *testing.T) {
			os.Setenv("CONJUR_AUTHN_LOGIN", Login)
			os.Setenv("CONJUR_AUTHN_API_KEY", APIKey)

			WithoutArgs(t)

			t.Run("Retrieves existing variable's defined value", func(t *testing.T) {
				variableIdentifier := "db/password"
				secretValue := fmt.Sprintf("secret-value-%v", rand.Intn(123456))
				policy := fmt.Sprintf(`
- !variable %s
`, variableIdentifier)

				config := conjurapi.Config{
					ApplianceURL: ApplianceURL,
					Account:      Account,
				}
				conjur, _ := conjurapi.NewClientFromKey(config, conjur_authn.LoginPair{Login, APIKey})

				conjur.LoadPolicy(
					conjurapi.PolicyModePost,
					"root",
					strings.NewReader(policy),
				)
				defer conjur.LoadPolicy(
					conjurapi.PolicyModePut,
					"root",
					strings.NewReader(""),
				)

				conjur.AddSecret(variableIdentifier, secretValue)

				stdout, _, err := RunCommand(PackageName, variableIdentifier)

				assert.NoError(t, err)
				assert.Equal(t, stdout.String(), secretValue)
			})

			t.Run("Returns error on non-existent variable", func(t *testing.T) {
				variableIdentifier := "non-existent-variable"

				_, stderr, err := RunCommand(PackageName, variableIdentifier)

				assert.Error(t, err)
				assert.Contains(t, stderr.String(), "not found")
			})

			t.Run("Given a non-existent Login is set", func(t *testing.T) {
				os.Setenv("CONJUR_AUTHN_LOGIN", "non-existent-user")

				t.Run("Returns 401", func(t *testing.T) {
					variableIdentifier := "existent-or-non-existent-variable"

					_, stderr, err := RunCommand(PackageName, variableIdentifier)

					assert.Error(t, err)
					assert.Contains(t, stderr.String(), "401")
				})
			})

			// Cleanup
			os.Unsetenv("CONJUR_AUTHN_LOGIN")
			os.Unsetenv("CONJUR_AUTHN_API_KEY")
		})

		t.Run("Given valid TokenFile credentials", func(t *testing.T) {

			getToken := fmt.Sprintf(`
token=$(curl --data "%s" "$CONJUR_APPLIANCE_URL/authn/$CONJUR_ACCOUNT/%s/authenticate")
echo $token
`, APIKey, Login)
			stdout, _, err := RunCommand("bash", "-c", getToken)

			assert.NoError(t, err)
			assert.Contains(t, stdout.String(), "signature")

			tokenFile, _ := ioutil.TempFile("", "existent-token-file")
			tokenFileName := tokenFile.Name()
			tokenFileContents := stdout.String()
			os.Remove(tokenFileName)
			go func() {
				ioutil.WriteFile(tokenFileName, []byte(tokenFileContents), 0600)
			}()
			defer os.Remove(tokenFileName)

			os.Setenv("CONJUR_AUTHN_TOKEN_FILE", tokenFileName)

			WithoutArgs(t)

			t.Run("Retrieves existent variable's defined value", func(t *testing.T) {
				variableIdentifier := "db/password"
				secretValue := fmt.Sprintf("secret-value-%v", rand.Intn(123456))
				policy := fmt.Sprintf(`
- !variable %s
`, variableIdentifier)

				config := conjurapi.Config{
					ApplianceURL: ApplianceURL,
					Account:      Account,
				}
				conjur, _ := conjurapi.NewClientFromKey(config, conjur_authn.LoginPair{Login, APIKey})

				conjur.LoadPolicy(
					conjurapi.PolicyModePost,
					"root",
					strings.NewReader(policy),
				)
				defer conjur.LoadPolicy(
					conjurapi.PolicyModePut,
					"root",
					strings.NewReader(""),
				)

				conjur.AddSecret(variableIdentifier, secretValue)

				stdout, _, err := RunCommand(PackageName, variableIdentifier)

				assert.NoError(t, err)
				assert.Equal(t, stdout.String(), secretValue)
			})

			t.Run("Returns error on non-existent variable", func(t *testing.T) {
				variableIdentifier := "non-existent-variable"

				_, stderr, err := RunCommand(PackageName, variableIdentifier)

				assert.Error(t, err)
				assert.Contains(t, stderr.String(), "not found in account")
			})

			t.Run("Given a non-existent TokenFile is set", func(t *testing.T) {
				os.Setenv("CONJUR_AUTHN_TOKEN_FILE", "non-existent-token-file")

				t.Run("Waits for longer than a second", func(t *testing.T) {
					timeout := time.After(1 * time.Second)
					unexpectedResponse := make(chan struct{})

					go func() {
						variableIdentifier := "existent-or-non-existent-variable"
						RunCommand(PackageName, variableIdentifier)
						unexpectedResponse <- struct{}{}
					}()

					select {
					case <-unexpectedResponse:
						assert.Fail(t, "unexpected response")
					case <-timeout:
						assert.True(t, true)
					}
				})

				// Cleanup
				os.Unsetenv("CONJUR_AUTHN_TOKEN_FILE")
			})
		})

		t.Run("Given no authentication credentials", func(t *testing.T) {
			WithoutArgs(t)

			t.Run("Returns with error on non-existent variable", func(t *testing.T) {
				variableIdentifier := "existent-or-non-existent-variable"

				_, stderr, err := RunCommand(PackageName, variableIdentifier)

				assert.Error(t, err)
				assert.Contains(t, stderr.String(), "at least one authentication strategy")
			})
		})
	})
}
