package test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackageEnterprise(t *testing.T) {
	Account := os.Getenv("CONJUR_ACCOUNT")
	Login := os.Getenv("CONJUR_AUTHN_LOGIN")

	ApplianceURL_V4 := os.Getenv("CONJUR_V4_APPLIANCE_URL")
	SSLCert_V4 := os.Getenv("CONJUR_V4_SSL_CERTIFICATE")
	APIKey_V4 := os.Getenv("CONJUR_V4_AUTHN_API_KEY")

	Path := os.Getenv("PATH")

	t.Run("Given valid V4 appliance configuration", func(t *testing.T) {
		e := ClearEnv()
		defer e.RestoreEnv()
		os.Setenv("PATH", Path)

		os.Setenv("CONJUR_MAJOR_VERSION", "4")
		os.Setenv("CONJUR_APPLIANCE_URL", ApplianceURL_V4)
		os.Setenv("CONJUR_ACCOUNT", Account)
		os.Setenv("CONJUR_AUTHN_LOGIN", Login)
		os.Setenv("CONJUR_SSL_CERTIFICATE", SSLCert_V4)

		t.Run("Given valid APIKey credentials", func(t *testing.T) {
			os.Setenv("CONJUR_AUTHN_API_KEY", APIKey_V4)

			WithoutArgs(t)

			t.Run("Retrieves existing variable's defined value", func(t *testing.T) {
				variableIdentifier := "existent-variable-with-defined-value"
				secretValue := "existent-variable-defined-value"

				stdout, _, err := RunCommand(PackageName, variableIdentifier)

				assert.NoError(t, err)
				assert.Equal(t, stdout.String(), secretValue)
			})

			t.Run("Returns error on existent-variable-undefined-value", func(t *testing.T) {
				variableIdentifier := "existent-variable-with-undefined-value"

				_, stderr, err := RunCommand(PackageName, variableIdentifier)

				assert.Error(t, err)
				assert.Contains(t, stderr.String(), "Not Found")
			})

			t.Run("Returns error on non-existent variable", func(t *testing.T) {
				variableIdentifier := "non-existent-variable"

				_, stderr, err := RunCommand(PackageName, variableIdentifier)

				assert.Error(t, err)
				assert.Contains(t, stderr.String(), "not found")
			})
		})
	})
}
