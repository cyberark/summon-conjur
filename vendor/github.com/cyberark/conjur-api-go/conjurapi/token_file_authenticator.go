package conjurapi

import (
	"time"
	"os"
)

type TokenFileAuthenticator struct {
	TokenFile string `env:"CONJUR_AUTHN_TOKEN_FILE"`
	mTime time.Time
	MaxWaitTime time.Duration
}

func (a *TokenFileAuthenticator) RefreshToken() ([]byte, error) {
	maxWaitTime := a.MaxWaitTime
	if maxWaitTime == 0 {
		maxWaitTime = 10 * time.Millisecond
	}
	bytes, err := waitForTextFile(a.TokenFile, time.After(a.MaxWaitTime))
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
