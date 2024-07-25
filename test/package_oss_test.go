package test

import (
	"fmt"
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

	t.Run("version flag", func(t *testing.T) {
		e := ClearEnv()
		defer e.RestoreEnv()
		os.Setenv("PATH", Path)

		stdout, stderr, err := RunCommand(PackageName, "--version")

		assert.NoError(t, err)
		assert.Empty(t, stderr.String())
		assert.Equal(t, "unset-unset\n", stdout.String())
	})

	t.Run("Given no configuration and no authentication information", func(t *testing.T) {
		e := ClearEnv()
		defer e.RestoreEnv()
		os.Setenv("PATH", Path)

		variableIdentifier := "variable"
		_, stderr, err := RunCommand(PackageName, variableIdentifier)

		//When both config and auth information is missing, then config errors take priority
		assert.Error(t, err)
		assert.Contains(t, stderr.String(), "Failed creating a Conjur client: Must specify an ApplianceURL -- Must specify an Account")
	})

	t.Run("Given valid OSS configuration", func(t *testing.T) {
		e := ClearEnv()
		defer e.RestoreEnv()
		os.Setenv("PATH", Path)

		os.Setenv("CONJUR_APPLIANCE_URL", ApplianceURL)
		os.Setenv("CONJUR_ACCOUNT", Account)

		t.Run("Given valid APIKey credentials", func(t *testing.T) {
			os.Setenv("CONJUR_AUTHN_LOGIN", Login)
			os.Setenv("CONJUR_AUTHN_API_KEY", APIKey)
			os.Setenv("HOME", "/root") //Workaround for Conjur API sending a warning to the stderr

			t.Run("Given interactive mode active", func(t *testing.T) {
				t.Run("Retrieves multiple existing variable's values", func(t *testing.T) {
					variableIdentifierUsername := "db/username"
					variableIdentifierPassword := "db/password"

					secretValueUsername := fmt.Sprintf("secret-value-username")
					secretValuePassword := fmt.Sprintf("secret-value-%v", rand.Intn(123456))
					policy := fmt.Sprintf(`
- !variable %s
- !variable %s
`, variableIdentifierUsername, variableIdentifierPassword)

					config := conjurapi.Config{
						ApplianceURL: ApplianceURL,
						Account:      Account,
					}
					conjur, _ := conjurapi.NewClientFromKey(config, conjur_authn.LoginPair{Login: Login, APIKey: APIKey})

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

					conjur.AddSecret(variableIdentifierUsername, secretValueUsername)
					conjur.AddSecret(variableIdentifierPassword, secretValuePassword)
					values := []string{variableIdentifierUsername, variableIdentifierPassword}
					output, err := RunCommandInteractively(PackageName, values)

					assert.Nil(t, err)
					assert.Equal(t, EncodeStringToBase64(secretValueUsername), output[0])
					assert.Equal(t, EncodeStringToBase64(secretValuePassword), output[1])
				})
				t.Run("Returns error on non-existent variables", func(t *testing.T) {
					variableIdentifier1 := "non-existent-variable1"
					variableIdentifier2 := "non-existent-variable2"

					values := []string{variableIdentifier1, variableIdentifier2}

					_, err := RunCommandInteractively(PackageName, values)

					assert.Contains(t, string(err), "404 Not Found")
				})
			})
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
				conjur, _ := conjurapi.NewClientFromKey(config, conjur_authn.LoginPair{Login: Login, APIKey: APIKey})

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
				assert.Equal(t, secretValue, stdout.String())
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

			tokenFile, _ := os.CreateTemp("", "existent-token-file")
			tokenFileName := tokenFile.Name()
			tokenFileContents := stdout.String()
			os.Remove(tokenFileName)
			go func() {
				os.WriteFile(tokenFileName, []byte(tokenFileContents), 0600)
			}()
			defer os.Remove(tokenFileName)

			os.Setenv("CONJUR_AUTHN_TOKEN_FILE", tokenFileName)

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
				conjur, _ := conjurapi.NewClientFromKey(config, conjur_authn.LoginPair{Login: Login, APIKey: APIKey})

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
				assert.Equal(t, secretValue, stdout.String())
			})

			t.Run("Returns error on non-existent variable", func(t *testing.T) {
				variableIdentifier := "non-existent-variable"

				_, stderr, err := RunCommand(PackageName, variableIdentifier)

				assert.Error(t, err)
				assert.Contains(t, stderr.String(), "CONJ00076E Variable cucumber:variable:non-existent-variable is empty or not found")
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

			t.Run("Returns with error on non-existent variable", func(t *testing.T) {
				variableIdentifier := "existent-or-non-existent-variable"

				_, stderr, err := RunCommand(PackageName, variableIdentifier)

				assert.Error(t, err)
				assert.Contains(t, stderr.String(), "Failed creating a Conjur client")
			})
		})
	})
}
