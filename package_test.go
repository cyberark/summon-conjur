package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
	"os/exec"
	"strings"
	"fmt"
	"github.com/cyberark/conjur-api-go/conjurapi"
	"math/rand"
	"io/ioutil"
)

func WithoutArgs()  {
	Convey("Given summon-conjur is run with no arguments", func() {
		out, err := exec.Command(PackageName).Output()

		Convey("Returns with error", func() {
			So(err, ShouldNotBeNil)
			So(string(out), ShouldEqual, "A variable name must be given as the first and only argument!\n")
		})
	})
}

const PackageName = "summon-conjur"
func TestPackage(t *testing.T) {

	if os.Getenv("TEST_PACKAGE") == "" {
		return
	}

	Convey("Given a compiled summon-conjur package", t, func() {
		ApplianceURL := os.Getenv("CONJUR_APPLIANCE_URL")
		Account := os.Getenv("CONJUR_ACCOUNT")
		Login := os.Getenv("CONJUR_AUTHN_LOGIN")
		APIKey := os.Getenv("CONJUR_AUTHN_API_KEY")
		Path := os.Getenv("PATH")

		Convey("Given no configuration and no authentication information", func() {
			e := ClearEnv()
			defer e.RestoreEnv()
			os.Setenv("PATH", Path)

			WithoutArgs()
		})

		Convey("Given valid configuration", func() {
			e := ClearEnv()
			defer e.RestoreEnv()
			os.Setenv("PATH", Path)

			os.Setenv("CONJUR_APPLIANCE_URL", ApplianceURL)
			os.Setenv("CONJUR_ACCOUNT", Account)

			Convey("Given valid APIKey credentials", func() {
				os.Setenv("CONJUR_AUTHN_LOGIN", Login)
				os.Setenv("CONJUR_AUTHN_API_KEY", APIKey)

				WithoutArgs()

				Convey("Retrieves existent variable's defined value", func() {
					variableIdentifier := "db/password"
					secretValue := fmt.Sprintf("secret-value-%v", rand.Intn(123456))
					policy := fmt.Sprintf(`
- !variable %s
`, variableIdentifier)

					config := conjurapi.Config{
						ApplianceURL: ApplianceURL,
						Account: Account,
					}
					conjur, _ := conjurapi.NewClientFromKey(config, Login, APIKey)

					conjur.LoadPolicy(
						"root",
						strings.NewReader(policy),
					)
					defer conjur.LoadPolicy(
						"root",
						strings.NewReader(""),
					)

					conjur.AddSecret(variableIdentifier, secretValue)

					out, err := exec.Command(PackageName, variableIdentifier).Output()

					So(err, ShouldBeNil)
					So(string(out), ShouldEqual, secretValue)
				})

				Convey("Returns 404 on non-existent variable", func() {
					variableIdentifier := "non-existent-variable"

					out, err := exec.Command(PackageName, variableIdentifier).Output()

					So(err, ShouldNotBeNil)
					So(string(out), ShouldContainSubstring, "404")
				})

				Convey("Given a non-existent Login is set", func() {
					os.Setenv("CONJUR_AUTHN_LOGIN", "non-existent-user")

					Convey("Returns 401", func() {
						variableIdentifier := "existent-or-non-existent-variable"

						out, err := exec.Command(PackageName, variableIdentifier).Output()

						So(err, ShouldNotBeNil)
						So(string(out), ShouldContainSubstring, "401")
					})
				})
			})

			Convey("Given valid TokenFile credentials", func() {

				getToken := fmt.Sprintf(`
token=$(curl --data "%s" "$CONJUR_APPLIANCE_URL/authn/$CONJUR_ACCOUNT/%s/authenticate")
echo $token
`, APIKey, Login)
				out, err := exec.Command("bash", "-c", getToken).Output()

				So(err, ShouldBeNil)
				So(string(out), ShouldContainSubstring, "data")

				tokenFile :=	"/tmp/existent-token-file"
				tokenFileContents := string(out)
				os.Remove(tokenFile)
				go func() {
					ioutil.WriteFile(tokenFile, []byte(tokenFileContents), 0644)
				}()
				defer os.Remove(tokenFile)

				os.Setenv("CONJUR_AUTHN_TOKEN_FILE", tokenFile)

				WithoutArgs()

				Convey("Retrieves existent variable's defined value", func() {
					variableIdentifier := "db/password"
					secretValue := fmt.Sprintf("secret-value-%v", rand.Intn(123456))
					policy := fmt.Sprintf(`
- !variable %s
`, variableIdentifier)

					config := conjurapi.Config{
						ApplianceURL: ApplianceURL,
						Account: Account,
					}
					conjur, _ := conjurapi.NewClientFromKey(config, Login, APIKey)

					conjur.LoadPolicy(
						"root",
						strings.NewReader(policy),
					)
					defer conjur.LoadPolicy(
						"root",
						strings.NewReader(""),
					)

					conjur.AddSecret(variableIdentifier, secretValue)

					out, err := exec.Command(PackageName, variableIdentifier).Output()

					So(err, ShouldBeNil)
					So(string(out), ShouldEqual, secretValue)
				})

				Convey("Returns 404 on non-existent variable", func() {
					variableIdentifier := "non-existent-variable"

					out, err := exec.Command(PackageName, variableIdentifier).Output()

					So(err, ShouldNotBeNil)
					So(string(out), ShouldContainSubstring, "404")
				})

				Convey("Given a non-existent TokenFile is set", func() {
					os.Setenv("CONJUR_AUTHN_TOKEN_FILE", "non-existent-user")

					Convey("Returns with timed out error", func() {
						variableIdentifier := "existent-or-non-existent-variable"

						out, err := exec.Command(PackageName, variableIdentifier).Output()

						So(err, ShouldNotBeNil)
						So(string(out), ShouldContainSubstring, "timed out")
					})
				})
			})

			Convey("Given no authentication credentials", func() {

				WithoutArgs()

				Convey("Returns with on non-existent variable", func() {
					variableIdentifier := "existent-or-non-existent-variable"

					out, err := exec.Command(PackageName, variableIdentifier).Output()

					So(err, ShouldNotBeNil)
					So(string(out), ShouldContainSubstring, "at least one authentication strategy")
				})
			})
		})


	})
}
