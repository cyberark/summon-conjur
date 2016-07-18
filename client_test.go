package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadAuthToken(t *testing.T) {
	Convey("Given an environment that contains a Conjur auth token", t, func() {
		e := ClearEnv()
		defer e.RestoreEnv()

		pem, _ := ioutil.ReadFile("test/files/real.pem")
		authToken := "eyJkYXRhIjo..."

		os.Setenv("CONJUR_AUTHN_TOKEN", authToken)
		os.Setenv("CONJUR_APPLIANCE_URL", "https://url.to.conjur")
		os.Setenv("CONJUR_SSL_CERTIFICATE", string(pem))

		Convey("When I create a Conjur client", func() {
			c, err := NewConjurClient()
			Convey("There shouldn't be an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the client uses the auth token from environment", func() {
				So(c.AuthToken, ShouldEqual, authToken)
			})
		})
	})
}
