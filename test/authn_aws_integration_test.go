package test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var awsRolePolicy = `
# The host ID needs to match the AWS ARN of the role we wish to authenticate
- !host 601277729239/InstanceReadJenkinsExecutorHostFactoryToken

- !variable db/password
- !permit
  role: !host 601277729239/InstanceReadJenkinsExecutorHostFactoryToken
  privilege: [ read, execute ]
  resource: !variable db/password
`

var awsAuthnPolicy = `
- !policy
  id: conjur/authn-iam/test
  body:
    - !webservice

    - !group clients

    - !permit
      role: !group clients
      privilege: [ read, authenticate ]
      resource: !webservice

    # Give the host permission to authenticate using the IAM Authenticator
    - !grant
      role: !group clients
      member: !host /601277729239/InstanceReadJenkinsExecutorHostFactoryToken
`

func TestAuthnAWSIntegration(t *testing.T) {
	if strings.ToLower(os.Getenv("TEST_AWS")) != "true" {
		t.Skip("Skipping AWS IAM authn test")
	}

	ApplianceURL := os.Getenv("CONJUR_APPLIANCE_URL")
	Account := os.Getenv("CONJUR_ACCOUNT")
	Login := os.Getenv("CONJUR_AUTHN_LOGIN")
	APIKey := os.Getenv("CONJUR_AUTHN_API_KEY")
	Path := os.Getenv("PATH")

	t.Run("Given a Conjur instance with an authn-iam authenticator", func(t *testing.T) {
		// Load AWS IAM authn policy

		conjur, _ := createConjurClient(ApplianceURL, Account, Login, APIKey)

		rootCleanup, err := loadPolicy(conjur, "root", awsRolePolicy)
		require.NoError(t, err)
		defer rootCleanup()

		_, err = loadPolicy(conjur, "root", awsAuthnPolicy)
		require.NoError(t, err)

		conjur.EnableAuthenticator("iam", "test", true)

		variableIdentifier := "db/password"
		secretValue := addSecretWithRandomValue(conjur, variableIdentifier)

		t.Run("Given invalid authn-iam configuration", func(t *testing.T) {
			e := setupTestEnvironment(Path, ApplianceURL, Account, "", "")
			os.Setenv("CONJUR_AUTHN_TYPE", "iam")
			os.Setenv("CONJUR_SERVICE_ID", "nonexistent-service")
			os.Setenv("CONJUR_AUTHN_JWT_HOST_ID", "invalid/host-id")
			defer e.RestoreEnv()

			t.Run("Fails to authenticate and returns an error", func(t *testing.T) {
				assertCommandError(t, variableIdentifier, "401 Unauthorized")
			})
		})

		t.Run("Given valid authn-iam configuration", func(t *testing.T) {
			e := setupTestEnvironment(Path, ApplianceURL, Account, "", "")
			os.Setenv("CONJUR_AUTHN_TYPE", "iam")
			os.Setenv("CONJUR_SERVICE_ID", "test")
			os.Setenv("CONJUR_AUTHN_JWT_HOST_ID", "601277729239/InstanceReadJenkinsExecutorHostFactoryToken")
			defer e.RestoreEnv()

			t.Run("Retrieves a variable", func(t *testing.T) {
				// Then attempt to authenticate and retrieve a secret using authn-iam
				stdout, _, err := RunCommand(PackageName, variableIdentifier)
				assert.NoError(t, err)
				assert.Equal(t, secretValue, stdout.String())
			})
		})
	})
}
