// +build enterprise

package main

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	. "github.com/playscale/goconvey/convey"
)

func RunCommand(name string, arg ...string) (bytes.Buffer, bytes.Buffer, error) {
	cmd := exec.Command(name, arg...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout, stderr, err
}

func WithoutArgs() {
	Convey("Given summon-conjur is run with no arguments", func() {
		_, stderr, err := RunCommand(PackageName)

		Convey("Returns with error", func() {
			So(err, ShouldNotBeNil)
			So(stderr.String(), ShouldEqual, `Usage of summon-conjur:
  -h, --help
	show help (default: false)
  -V, --version
	show version (default: false)
  -v, --verbose
	be verbose (default: false)
`)
		})
	})
}

const PackageName = "summon-conjur"

func TestPackageEnterprise(t *testing.T) {
	Account := os.Getenv("CONJUR_ACCOUNT")
	Login := os.Getenv("CONJUR_AUTHN_LOGIN")

	ApplianceURL_V4 := os.Getenv("CONJUR_V4_APPLIANCE_URL")
	SSLCert_V4 := os.Getenv("CONJUR_V4_SSL_CERTIFICATE")
	APIKey_V4 := os.Getenv("CONJUR_V4_AUTHN_API_KEY")

	Path := os.Getenv("PATH")

	Convey("Given a compiled summon-conjur package", t, func() {
		Convey("Given a v4 appliance", func() {
			Convey("Given valid configuration", func() {
				e := ClearEnv()
				defer e.RestoreEnv()
				os.Setenv("PATH", Path)

				os.Setenv("CONJUR_MAJOR_VERSION", "4")
				os.Setenv("CONJUR_APPLIANCE_URL", ApplianceURL_V4)
				os.Setenv("CONJUR_ACCOUNT", Account)
				os.Setenv("CONJUR_AUTHN_LOGIN", Login)
				os.Setenv("CONJUR_SSL_CERTIFICATE", SSLCert_V4)

				Convey("Given valid APIKey credentials", func() {
					os.Setenv("CONJUR_AUTHN_API_KEY", APIKey_V4)

					WithoutArgs()

					Convey("Retrieves existing variable's defined value", func() {
						variableIdentifier := "existent-variable-with-defined-value"
						secretValue := "existent-variable-defined-value"

						stdout, _, err := RunCommand(PackageName, variableIdentifier)

						So(err, ShouldBeNil)
						So(stdout.String(), ShouldEqual, secretValue)
					})

					Convey("Returns error on existent-variable-undefined-value", func() {
						variableIdentifier := "existent-variable-with-undefined-value"

						_, stderr, err := RunCommand(PackageName, variableIdentifier)

						So(err, ShouldNotBeNil)
						So(stderr.String(), ShouldContainSubstring, "Not Found")
					})

					Convey("Returns error on non-existent variable", func() {
						variableIdentifier := "non-existent-variable"

						_, stderr, err := RunCommand(PackageName, variableIdentifier)

						So(err, ShouldNotBeNil)
						So(stderr.String(), ShouldContainSubstring, "not found")
					})

				})
			})
		})
	})
}
