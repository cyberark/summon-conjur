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
	. "github.com/smartystreets/goconvey/convey"
)

func TestPackageOSS(t *testing.T) {
	ApplianceURL := os.Getenv("CONJUR_APPLIANCE_URL")
	Account := os.Getenv("CONJUR_ACCOUNT")
	Login := os.Getenv("CONJUR_AUTHN_LOGIN")
	APIKey := os.Getenv("CONJUR_AUTHN_API_KEY")

	Path := os.Getenv("PATH")

	Convey("Given no configuration and no authentication information", t, func() {
		e := ClearEnv()
		defer e.RestoreEnv()
		os.Setenv("PATH", Path)

		WithoutArgs()
	})

	Convey("Given valid V5 OSS configuration", t, func() {
		e := ClearEnv()
		defer e.RestoreEnv()
		os.Setenv("PATH", Path)

		os.Setenv("CONJUR_APPLIANCE_URL", ApplianceURL)
		os.Setenv("CONJUR_ACCOUNT", Account)

		Convey("Given valid APIKey credentials", func() {
			os.Setenv("CONJUR_AUTHN_LOGIN", Login)
			os.Setenv("CONJUR_AUTHN_API_KEY", APIKey)

			WithoutArgs()

			Convey("Retrieves existing variable's defined value", func() {
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

				So(err, ShouldBeNil)
				So(stdout.String(), ShouldEqual, secretValue)
			})

			Convey("Returns error on non-existent variable", func() {
				variableIdentifier := "non-existent-variable"

				_, stderr, err := RunCommand(PackageName, variableIdentifier)

				So(err, ShouldNotBeNil)
				So(stderr.String(), ShouldContainSubstring, "not found")
			})

			Convey("Given a non-existent Login is set", func() {
				os.Setenv("CONJUR_AUTHN_LOGIN", "non-existent-user")

				Convey("Returns 401", func() {
					variableIdentifier := "existent-or-non-existent-variable"

					_, stderr, err := RunCommand(PackageName, variableIdentifier)

					So(err, ShouldNotBeNil)
					So(stderr.String(), ShouldContainSubstring, "401")
				})
			})
		})

		Convey("Given valid TokenFile credentials", func() {

			getToken := fmt.Sprintf(`
token=$(curl --data "%s" "$CONJUR_APPLIANCE_URL/authn/$CONJUR_ACCOUNT/%s/authenticate")
echo $token
`, APIKey, Login)
			stdout, _, err := RunCommand("bash", "-c", getToken)

			So(err, ShouldBeNil)
			So(stdout.String(), ShouldContainSubstring, "signature")

			tokenFile, _ := ioutil.TempFile("", "existent-token-file")
			tokenFileName := tokenFile.Name()
			tokenFileContents := stdout.String()
			os.Remove(tokenFileName)
			go func() {
				ioutil.WriteFile(tokenFileName, []byte(tokenFileContents), 0600)
			}()
			defer os.Remove(tokenFileName)

			os.Setenv("CONJUR_AUTHN_TOKEN_FILE", tokenFileName)

			WithoutArgs()

			Convey("Retrieves existent variable's defined value", func() {
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

				So(err, ShouldBeNil)
				So(stdout.String(), ShouldEqual, secretValue)
			})

			Convey("Returns error on non-existent variable", func() {
				variableIdentifier := "non-existent-variable"

				_, stderr, err := RunCommand(PackageName, variableIdentifier)

				So(err, ShouldNotBeNil)
				So(stderr.String(), ShouldContainSubstring, "not found in account")
			})

			Convey("Given a non-existent TokenFile is set", func() {
				os.Setenv("CONJUR_AUTHN_TOKEN_FILE", "non-existent-token-file")

				Convey("Waits for longer than a second", func() {
					timeout := time.After(1 * time.Second)
					unexpected_response := make(chan int)

					go func() {
						variableIdentifier := "existent-or-non-existent-variable"
						RunCommand(PackageName, variableIdentifier)
						unexpected_response <- 1
					}()

					select {
					case <-unexpected_response:
						So("receive unexpected response", ShouldEqual, "not receive unexpected response")
					case <-timeout:
						So(true, ShouldEqual, true)
					}
				})
			})
		})

		Convey("Given no authentication credentials", func() {

			WithoutArgs()

			Convey("Returns with error on non-existent variable", func() {
				variableIdentifier := "existent-or-non-existent-variable"

				_, stderr, err := RunCommand(PackageName, variableIdentifier)

				So(err, ShouldNotBeNil)
				So(stderr.String(), ShouldContainSubstring, "at least one authentication strategy")
			})
		})
	})
}
