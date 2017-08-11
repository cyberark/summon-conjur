package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
	"os/exec"
	"strings"
	"fmt"
	"github.com/conjurinc/api-go/conjurapi"
	"math/rand"
)

func WithoutArgs()  {
	Convey("When I run the package with no arguments an error should be presented", func() {
		out, err := exec.Command("summon-conjur").Output()

		So(err, ShouldNotBeNil)
		So(string(out), ShouldEqual, "A variable name must be given as the first and only argument!\n")
	})
}

func TestPackage(t *testing.T) {

	if os.Getenv("TEST_PACKAGE") == "" {
		return
	}

	Convey("Given a compiled summon-conjur package", t, func() {
		appliance_url := os.Getenv("CONJUR_APPLIANCE_URL")
		account := os.Getenv("CONJUR_ACCOUNT")
		api_key := os.Getenv("CONJUR_API_KEY")
		path := os.Getenv("PATH")

		Convey("Given no configuration", func() {
			e := ClearEnv()
			defer e.RestoreEnv()
			os.Setenv("PATH", path)

			WithoutArgs()
		})

		Convey("Given all configuration", func() {
			e := ClearEnv()
			defer e.RestoreEnv()
			os.Setenv("PATH", path)

			os.Setenv("CONJUR_APPLIANCE_URL", appliance_url)
			os.Setenv("CONJUR_ACCOUNT", account)
			os.Setenv("CONJUR_AUTHN_LOGIN", "admin")
			os.Setenv("CONJUR_AUTHN_API_KEY", api_key)

			WithoutArgs()

			Convey("Existent and assigned variable is retrieved", func() {
				variable_identifier := "db/password"
				secret_value := fmt.Sprintf("secret-value-%v", rand.Intn(123456))
				policy := fmt.Sprintf(`
- !variable %s
`, variable_identifier)

				config := conjurapi.Config{
					ApplianceUrl: os.Getenv("CONJUR_APPLIANCE_URL"),
					Account: os.Getenv("CONJUR_ACCOUNT"),
					Username: os.Getenv("CONJUR_AUTHN_LOGIN"),
					APIKey: os.Getenv("CONJUR_AUTHN_API_KEY"),
				}
				conjur := conjurapi.NewClient(config)

				conjur.LoadPolicy(
					"root",
					strings.NewReader(policy),
				)
				defer conjur.LoadPolicy(
					"root",
					strings.NewReader(""),
				)

				conjur.AddSecret(variable_identifier, secret_value)

				out, err := exec.Command("summon-conjur", variable_identifier).Output()

				So(err, ShouldBeNil)
				So(string(out), ShouldEqual, secret_value)
			})

			Convey("Non-existent variable returns 404", func() {
				variable_identifier := "non-existent-variable"

				out, err := exec.Command("summon-conjur", variable_identifier).Output()

				So(err, ShouldNotBeNil)
				So(string(out), ShouldContainSubstring, "404")
			})

			Convey("When a non-existent username is set", func() {
				os.Setenv("CONJUR_AUTHN_LOGIN", "non-existent-user")

				Convey("Variable fetching returns 401", func() {
					variable_identifier := "existent-or-non-existent-variable"

					out, err := exec.Command("summon-conjur", variable_identifier).Output()

					So(err, ShouldNotBeNil)
					So(string(out), ShouldContainSubstring, "401")
				})
			})
		})
	})
}
