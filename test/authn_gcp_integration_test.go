package test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var gcpRolesPolicy = `
- !host
  id: test-app
  annotations:
    authn-gcp/project-id: %q
- !variable db/password

- !permit
  role: !host test-app
  privilege: [ read, execute ]
  resource: !variable db/password
`

var gcpAuthnPolicy = `
- !policy
  id: conjur/authn-gcp
  body:
  - !webservice

  - !group apps

  - !permit
    role: !group apps
    privilege: [ read, authenticate ]
    resource: !webservice

  # Give the host permission to authenticate using the GCP Authenticator
  - !grant
    role: !group apps
    member: !host /test-app
`

func TestAuthnGCPIntegration(t *testing.T) {
	if strings.ToLower(os.Getenv("TEST_GCP")) != "true" {
		t.Skip("Skipping GCP authn test")
	}

	// Replace placeholder in policy with actual project ID
	projectID := os.Getenv("GCP_PROJECT_ID")
	// Fetch the GCP token from environment variable
	prefetchedToken := os.Getenv("GCP_ID_TOKEN")
	if projectID == "" || prefetchedToken == "" {
		t.Fatal("GCP_PROJECT_ID and GCP_ID_TOKEN must be set to run this test")
	}

	ApplianceURL := os.Getenv("CONJUR_APPLIANCE_URL")
	Account := os.Getenv("CONJUR_ACCOUNT")
	Login := os.Getenv("CONJUR_AUTHN_LOGIN")
	APIKey := os.Getenv("CONJUR_AUTHN_API_KEY")
	Path := os.Getenv("PATH")

	t.Run("Given a Conjur instance with an authn-gcp authenticator", func(t *testing.T) {
		rolesPolicy := fmt.Sprintf(gcpRolesPolicy, projectID)

		conjur, _ := createConjurClient(ApplianceURL, Account, Login, APIKey)

		rootCleanup, err := loadPolicy(conjur, "root", rolesPolicy)
		require.NoError(t, err)
		defer rootCleanup()

		_, err = loadPolicy(conjur, "root", gcpAuthnPolicy)
		require.NoError(t, err)

		conjur.EnableAuthenticator("gcp", "", true)

		variableIdentifier := "db/password"
		secretValue := addSecretWithRandomValue(conjur, variableIdentifier)

		t.Run("Given invalid authn-gcp configuration", func(t *testing.T) {
			e := setupTestEnvironment(Path, ApplianceURL, Account, "", "")
			os.Setenv("CONJUR_AUTHN_TYPE", "gcp")
			os.Setenv("CONJUR_AUTHN_JWT_HOST_ID", "invalid/host-id")
			os.Unsetenv("CONJUR_AUTHN_JWT_TOKEN")
			defer e.RestoreEnv()

			t.Run("Fails to authenticate and returns an error", func(t *testing.T) {
				assertCommandError(t, variableIdentifier, "Request failed for GCP Metadata token")
			})
		})

		t.Run("Given valid authn-gcp configuration", func(t *testing.T) {
			e := setupTestEnvironment(Path, ApplianceURL, Account, "", "")
			os.Setenv("CONJUR_AUTHN_TYPE", "gcp")
			os.Setenv("CONJUR_AUTHN_JWT_HOST_ID", "test-app")
			os.Setenv("CONJUR_AUTHN_JWT_TOKEN", prefetchedToken)
			defer e.RestoreEnv()

			t.Run("Retrieves a variable", func(t *testing.T) {
				// Then attempt to authenticate and retrieve a secret using authn-gcp
				stdout, _, err := RunCommand(PackageName, variableIdentifier)
				assert.NoError(t, err)
				assert.Equal(t, secretValue, stdout.String())
			})
		})
	})
}
