package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func splitEq(s string) (string, string) {
	a := strings.SplitN(s, "=", 2)
	return a[0], a[1]
}

type envSnapshot struct {
	env []string
}

func clearEnv() *envSnapshot {
	e := os.Environ()

	for _, s := range e {
		k, _ := splitEq(s)
		os.Setenv(k, "")
	}
	return &envSnapshot{env: e}
}

func (e *envSnapshot) restoreEnv() {
	clearEnv()
	for _, s := range e.env {
		k, v := splitEq(s)
		os.Setenv(k, v)
	}
}

func TestLoadConfig(t *testing.T) {
	Convey("Given an environment with conjur config vars only", t, func() {

		e := clearEnv()
		defer e.restoreEnv()

		expected := &Config{
			APIKey:       "env-api-key",
			ApplianceUrl: "env-app-url",
			SSLCertPath:  "env-cert-file",
			Username:     "env-username",
		}

		os.Setenv("CONJUR_API_KEY", expected.APIKey)
		os.Setenv("CONJUR_APPLIANCE_URL", expected.ApplianceUrl)
		os.Setenv("CONJUR_CERT_FILE", expected.SSLCertPath)
		os.Setenv("CONJUR_AUTHN_LOGIN", expected.Username)
		Convey("When I load the config", func() {
			c, err := LoadConfig()
			Convey("There shouldn't be an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("And the result have the right properties", func() {
				So(c, ShouldResemble, expected)
			})
		})
	})

	Convey("When I write a conjurrc to a tempfile and set CONJURRC to its path", t, func() {
		e := clearEnv()
		defer e.restoreEnv()
		rcFile, err := ioutil.TempFile("", "conjurrc-test")
		if err != nil {
			panic(err)
		}
		rcFile.WriteString(`
appliance_url: rc-app-url
cert_file: rc-cert-file
ignore_me: please`)
		os.Setenv("CONJURRC", rcFile.Name())
		// also set these so that the config will validate
		os.Setenv("CONJUR_AUTHN_LOGIN", "dummy")
		os.Setenv("CONJUR_API_KEY", "dummy")
		Convey("And I load the config", func() {
			c, err := LoadConfig()
			Convey("It should succeed", func() {
				So(err, ShouldBeNil)
			})

			Convey("It should have the appliance url and cert path from the rc file", func() {
				So(c.ApplianceUrl, ShouldEqual, "rc-app-url")
				So(c.SSLCertPath, ShouldEqual, "rc-cert-file")
			})
		})
	})

	Convey("When I use a fake home dir with a netrc in it", t, func() {
		env := clearEnv()
		defer env.restoreEnv()

		fakeHome, err := ioutil.TempDir("", "test-netrc")
		if err != nil {
			panic(err)
		}

		os.Setenv("HOME", fakeHome)

		Convey("And I write a netrc to it", func() {
			netrc := `machine https://foo.bar.com/api/authn
    login foo
    password s3cr3t
`
			err := ioutil.WriteFile(fakeHome+"/.netrc", []byte(netrc), 0600)

			if err != nil {
				panic(err)
			}

			// put the remaining vars in the env so the config validates and loads the right machine
			os.Setenv("CONJUR_APPLIANCE_URL", "https://foo.bar.com/api")
			os.Setenv("CONJUR_CERT_FILE", "cert-file")

			Convey("When I load the config", func() {
				c, err := LoadConfig()

				Convey("It should succeed", func() {
					So(err, ShouldBeNil)
				})

				Convey("And the config should have values from the netrc file", func() {
					So(c.Username, ShouldEqual, "foo")
					So(c.APIKey, ShouldEqual, "s3cr3t")
				})
			})

		})
	})
}
