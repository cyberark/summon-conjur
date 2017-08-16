package conjurapi

import (
	"time"
	"os"
)

type TokenFileAuthenticator struct {
	TokenFile string `env:"CONJUR_AUTHN_TOKEN_FILE"`
	mTime time.Time
}

func (a *TokenFileAuthenticator) RefreshToken() ([]byte, error) {
	bytes, err := waitForTextFile(a.TokenFile, time.After(time.Second*10))
	if err == nil {
		fi, _ := os.Stat(a.TokenFile)
		a.mTime = fi.ModTime()
	}
	return bytes, err
}

func (a *TokenFileAuthenticator) NeedsTokenRefresh() bool {
	fi, _ := os.Stat(a.TokenFile)
	return a.mTime != fi.ModTime()
}
