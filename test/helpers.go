package test

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func splitEq(s string) (string, string) {
	a := strings.SplitN(s, "=", 2)
	return a[0], a[1]
}

type envSnapshot struct {
	env []string
}

func ClearEnv() *envSnapshot {
	e := os.Environ()

	for _, s := range e {
		k, _ := splitEq(s)
		os.Setenv(k, "")
	}
	return &envSnapshot{env: e}
}

func (e *envSnapshot) RestoreEnv() {
	ClearEnv()
	for _, s := range e.env {
		k, v := splitEq(s)
		os.Setenv(k, v)
	}
}

func RunCommand(name string, arg ...string) (bytes.Buffer, bytes.Buffer, error) {
	cmd := exec.Command(name, arg...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout, stderr, err
}

func WithoutArgs(t *testing.T) {
	t.Run("Given summon-conjur is run with no arguments", func(t *testing.T) {
		_, stderr, err := RunCommand(PackageName)

		t.Run("Returns with error", func(t *testing.T) {
			assert.Error(t, err)
			assert.Equal(t, stderr.String(), `Usage of summon-conjur:
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
