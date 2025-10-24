package test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var azureRolesPolicyTemplate = `
# The host ID needs to match the Azure ARN of the role we wish to authenticate
- !host
  id: azureVM
  annotations:
    authn-azure/subscription-id: %q
    authn-azure/resource-group: %q

- !variable db/password

- !permit
  role: !host azureVM
  privilege: [ read, execute ]
  resource: !variable db/password
`

var azureAuthnPolicy = `
- !policy
  id: conjur/authn-azure/test
  body:
  - !webservice

  - !variable
    id: provider-uri

  - !group apps

  - !permit
    role: !group apps
    privilege: [ read, authenticate ]
    resource: !webservice

  # Give the host permission to authenticate using the Azure Authenticator
  - !grant
    role: !group apps
    member: !host /azureVM
`

func TestAuthnAzureIntegration(t *testing.T) {
	if strings.ToLower(os.Getenv("TEST_AZURE")) != "true" {
		t.Skip("Skipping Azure authn test")
	}

	if os.Getenv("AZURE_SUBSCRIPTION_ID") == "" ||
		os.Getenv("AZURE_RESOURCE_GROUP") == "" {
		t.Fatal("AZURE_SUBSCRIPTION_ID and AZURE_RESOURCE_GROUP must be set to run this test")
	}

	ApplianceURL := os.Getenv("CONJUR_APPLIANCE_URL")
	Account := os.Getenv("CONJUR_ACCOUNT")
	Login := os.Getenv("CONJUR_AUTHN_LOGIN")
	APIKey := os.Getenv("CONJUR_AUTHN_API_KEY")
	Path := os.Getenv("PATH")

	t.Run("Given a Conjur instance with an authn-azure authenticator", func(t *testing.T) {
		rolesPolicy := fmt.Sprintf(azureRolesPolicyTemplate,
			os.Getenv("AZURE_SUBSCRIPTION_ID"),
			os.Getenv("AZURE_RESOURCE_GROUP"))

		conjur, _ := createConjurClient(ApplianceURL, Account, Login, APIKey)

		rootCleanup, err := loadPolicy(conjur, "root", rolesPolicy)
		require.NoError(t, err)
		defer rootCleanup()

		_, err = loadPolicy(conjur, "root", azureAuthnPolicy)
		require.NoError(t, err)

		err = conjur.AddSecret("conjur/authn-azure/test/provider-uri", "https://sts.windows.net/df242c82-fe4a-47e0-b0f4-e3cb7f8104f1/")
		require.NoError(t, err)

		conjur.EnableAuthenticator("azure", "test", true)

		variableIdentifier := "db/password"
		secretValue := addSecretWithRandomValue(conjur, variableIdentifier)

		t.Run("Given invalid authn-azure configuration", func(t *testing.T) {
			e := setupTestEnvironment(Path, ApplianceURL, Account, "", "")
			os.Setenv("CONJUR_AUTHN_TYPE", "azure")
			os.Setenv("CONJUR_SERVICE_ID", "nonexistent-service")
			os.Setenv("CONJUR_AUTHN_JWT_HOST_ID", "invalid/host-id")
			defer e.RestoreEnv()

			t.Run("Fails to authenticate and returns an error", func(t *testing.T) {
				assertCommandError(t, variableIdentifier, "401 Unauthorized")
			})
		})

		t.Run("Given valid authn-azure configuration", func(t *testing.T) {
			e := setupTestEnvironment(Path, ApplianceURL, Account, "", "")
			os.Setenv("CONJUR_AUTHN_TYPE", "azure")
			os.Setenv("CONJUR_SERVICE_ID", "test")
			os.Setenv("CONJUR_AUTHN_JWT_HOST_ID", "azureVM")
			defer e.RestoreEnv()

			t.Run("Retrieves a variable", func(t *testing.T) {
				// Then attempt to authenticate and retrieve a secret using authn-azure
				stdout, _, err := RunCommand(PackageName, variableIdentifier)
				assert.NoError(t, err)
				assert.Equal(t, secretValue, stdout.String())
			})
		})
	})
}
