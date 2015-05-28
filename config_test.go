package main
import (
    "testing"
    "os"
    "github.com/stretchr/testify/assert"
    "io/ioutil"
)

func cloberEnv(home, pwd string, f func()) {
    oldHome := os.Getenv("HOME")
    oldPwd  := os.Getenv("PWD")

    // is this the right way to defer multiple statements?
    cleanup := func(){
        os.Setenv("HOME", oldHome)
        os.Setenv("PWD", oldPwd)
    }
    defer cleanup()

    os.Setenv("HOME", home)
    os.Setenv("PWD", pwd)
    f()
}
func TestLoadEnvConfig(t *testing.T) {
    cloberEnv("", "", func(){
        // don't load from conjurrc
        os.Setenv("CONJURRC", "")

        expected := &Config{
            APIKey: "env-api-key",
            ApplianceUrl: "env-app-url",
            SSLCertPath:  "env-cert-file",
            Username: "env-username",
        }

        // set env vars
        os.Setenv("CONJUR_API_KEY", expected.APIKey)
        os.Setenv("CONJUR_APPLIANCE_URL", expected.ApplianceUrl)
        os.Setenv("CONJUR_CERT_FILE", expected.SSLCertPath)
        os.Setenv("CONJUR_AUTHN_LOGIN", expected.Username)

        c,err := LoadConfig()
        assert.Nil(t, err)
        assert.Equal(t, expected, c, "loaded the wrong env config")
  })
}

func TestLoadConjurRC(t *testing.T) {
    rcPath,err := ioutil.TempFile("", "conjurrc-test")

    if err != nil { panic(err) }

    rcPath.WriteString(`
appliance_url: rc-app-url
cert_file: rc-cert-file
    `)
    expected := &Config{
        ApplianceUrl: "rc-app-url",
        SSLCertPath:  "rc-cert-file",
        APIKey: "dummy",
        Username: "dummy",
    }
    cloberEnv("","", func(){
        os.Setenv("CONJUR_AUTHN_LOGIN", "dummy")
        os.Setenv("CONJUR_API_KEY", "dummy")
        os.Setenv("CONJURRC", rcPath.Name())
        actual, err := LoadConfig()
        assert.Nil(t,err)
        assert.Equal(t, actual, expected, "loaded the wrong config")
    })
}


func TestReadNetrc(t *testing.T) {
    // BUT HOW?? Maybe testify/mock?  TODO for now
}

