package test

import (
	"fmt"
	"maps"
	"math/rand"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

		// When both config and auth information is missing, then config errors take priority
		assertCommandError(t, "variable", "Failed creating a Conjur client: Must specify an ApplianceURL -- Must specify an Account")
	})

	t.Run("Given valid OSS configuration", func(t *testing.T) {
		e := ClearEnv()
		defer e.RestoreEnv()
		os.Setenv("PATH", Path)

		os.Setenv("CONJUR_APPLIANCE_URL", ApplianceURL)
		os.Setenv("CONJUR_ACCOUNT", Account)

		t.Run("Given valid APIKey credentials", func(t *testing.T) {
			e := setupTestEnvironment(Path, ApplianceURL, Account, Login, APIKey)
			defer e.RestoreEnv()

			t.Run("Given interactive mode active", func(t *testing.T) {
				t.Run("Retrieves multiple existing variable's values", func(t *testing.T) {
					variableIdentifierUsername := "db/username"
					// file deepcode ignore HardcodedPassword/test: This is a test file
					variableIdentifierPassword := "db/password"

					secretValueUsername := "secret-value-username"
					// file deepcode ignore InsecurelyGeneratedPassword/test: This is a test file
					secretValuePassword := fmt.Sprintf("secret-value-%v", rand.Intn(123456))

					conjur, _ := createConjurClient(ApplianceURL, Account, Login, APIKey)
					cleanup, err := setupVariablePolicy(conjur, variableIdentifierUsername, variableIdentifierPassword)
					require.NoError(t, err)
					defer cleanup()

					conjur.AddSecret(variableIdentifierUsername, secretValueUsername)
					conjur.AddSecret(variableIdentifierPassword, secretValuePassword)
					values := []string{variableIdentifierUsername, variableIdentifierPassword}
					output, errStr := RunCommandInteractively(PackageName, values)

					assert.Nil(t, errStr)
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
				t.Run("Retrieves large number of variables", func(t *testing.T) {
					numVariables := 250
					variableMap := make(map[string]string)

					for i := range numVariables {
						variableIdentifier := fmt.Sprintf("variable-%d", i)
						secretValue := generateRandomString(1024)
						variableMap[variableIdentifier] = secretValue
					}

					variableNames := slices.Collect(maps.Keys(variableMap))

					conjur, _ := createConjurClient(ApplianceURL, Account, Login, APIKey)
					cleanup, err := setupVariablePolicy(conjur, variableNames...)
					require.NoError(t, err)
					defer cleanup()

					for key, value := range variableMap {
						err := conjur.AddSecret(key, value)
						assert.NoError(t, err, fmt.Sprintf("Failed to add secret for variable %s", key))
					}

					output, errStr := RunCommandInteractively(PackageName, variableNames)

					assert.Nil(t, errStr)
					assert.Len(t, output, len(variableNames))
					for i, value := range variableNames {
						assert.Equal(t, EncodeStringToBase64(variableMap[value]), output[i])
					}
				})
			})
			t.Run("Retrieves existing variable's defined value", func(t *testing.T) {
				variableIdentifier := "db/password"

				conjur, _ := createConjurClient(ApplianceURL, Account, Login, APIKey)
				cleanup, err := setupVariablePolicy(conjur, variableIdentifier)
				require.NoError(t, err)
				defer cleanup()

				secretValue := addSecretWithRandomValue(conjur, variableIdentifier)
				stdout, _, errStr := RunCommand(PackageName, variableIdentifier)

				assert.NoError(t, errStr)
				assert.Equal(t, secretValue, stdout.String())
			})

			t.Run("Returns error on non-existent variable", func(t *testing.T) {
				assertCommandError(t, "non-existent-variable", "not found")
			})

			t.Run("Given a non-existent Login is set", func(t *testing.T) {
				os.Setenv("CONJUR_AUTHN_LOGIN", "non-existent-user")

				t.Run("Returns 401", func(t *testing.T) {
					assertCommandError(t, "existent-or-non-existent-variable", "401 Unauthorized")
				})
			})

			// Cleanup
			os.Unsetenv("CONJUR_AUTHN_LOGIN")
			os.Unsetenv("CONJUR_AUTHN_API_KEY")
		})

		t.Run("Given valid TokenFile credentials", func(t *testing.T) {
			e := setupTestEnvironment(Path, ApplianceURL, Account, "", "")
			defer e.RestoreEnv()

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

				conjur, err := createConjurClient(ApplianceURL, Account, Login, APIKey)
				require.NoError(t, err)
				cleanup, err := setupVariablePolicy(conjur, variableIdentifier)
				require.NoError(t, err)
				defer cleanup()

				secretValue := addSecretWithRandomValue(conjur, variableIdentifier)
				stdout, _, err := RunCommand(PackageName, variableIdentifier)

				require.NoError(t, err)
				assert.Equal(t, secretValue, stdout.String())
			})

			t.Run("Returns error on non-existent variable", func(t *testing.T) {
				assertCommandError(t, "non-existent-variable", "CONJ00076E Variable cucumber:variable:non-existent-variable is empty or not found")
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
				assertCommandError(t, "existent-or-non-existent-variable", "Failed creating a Conjur client")
			})
		})
	})
}
