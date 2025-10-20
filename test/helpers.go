package test

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/cyberark/conjur-api-go/conjurapi/authn"
	"github.com/stretchr/testify/assert"
)

func createConjurClient(applianceURL, account, login, apiKey string) (*conjurapi.Client, error) {
	config := conjurapi.Config{
		ApplianceURL: applianceURL,
		Account:      account,
	}
	return conjurapi.NewClientFromKey(
		config,
		authn.LoginPair{Login: login, APIKey: apiKey},
	)
}

func setupVariablePolicy(conjur *conjurapi.Client, variableIdentifiers ...string) (func(), error) {
	policyBuilder := strings.Builder{}
	for _, identifier := range variableIdentifiers {
		policyBuilder.WriteString(fmt.Sprintf("- !variable %s\n", identifier))
	}

	// Load the variables in the root policy and return the cleanup function
	return loadPolicy(conjur, "root", policyBuilder.String())
}

func loadPolicy(conjur *conjurapi.Client, branch, policy string) (func(), error) {
	_, err := conjur.LoadPolicy(
		conjurapi.PolicyModePost,
		branch,
		strings.NewReader(policy),
	)

	if err != nil {
		return nil, err
	}

	// Return cleanup function that removes all content from specified policy branch
	return func() {
		conjur.LoadPolicy(
			conjurapi.PolicyModePut,
			branch,
			strings.NewReader(""),
		)
	}, nil
}

func addSecretWithRandomValue(conjur *conjurapi.Client, variableIdentifier string) string {
	secretValue := fmt.Sprintf("secret-value-%v", rand.Intn(123456))
	conjur.AddSecret(variableIdentifier, secretValue)
	// Return the generated secret value
	return secretValue
}

// assertCommandError runs the command and checks that it fails with expected error message
func assertCommandError(t *testing.T, variableIdentifier, expectedError string) {
	_, stderr, err := RunCommand(PackageName, variableIdentifier)
	assert.Error(t, err)
	assert.Contains(t, stderr.String(), expectedError)
}

// setupTestEnvironment sets up common environment variables for tests
func setupTestEnvironment(path, applianceURL, account, login, apiKey string) *envSnapshot {
	e := ClearEnv()
	os.Setenv("PATH", path)
	os.Setenv("CONJUR_APPLIANCE_URL", applianceURL)
	os.Setenv("CONJUR_ACCOUNT", account)
	if login != "" {
		os.Setenv("CONJUR_AUTHN_LOGIN", login)
	}
	if apiKey != "" {
		os.Setenv("CONJUR_AUTHN_API_KEY", apiKey)
	}
	os.Setenv("HOME", "/root") // Workaround for Conjur API sending a warning to stderr
	return e
}

func generateRandomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"
	ret := make([]byte, n)
	for i := range n {
		num := rand.Intn(len(letters))
		ret[i] = letters[num]
	}

	return string(ret)
}

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

// RunCommandInteractively takes multiple paths to secrets and returns their values in Base64 and a last error that occurred
func RunCommandInteractively(command string, values []string) ([][]byte, []byte) {
	errChan := make(chan []byte, 1)
	defer close(errChan)
	doneChan := make(chan bool, 1)
	cmd := exec.Command(command)

	stdinPipe, _ := cmd.StdinPipe()
	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()

	cmd.Start()

	go func() {
		defer stdinPipe.Close()
		for _, value := range values {
			fmt.Fprintln(stdinPipe, value)
		}
	}()

	var output [][]byte
	go func() {
		defer close(doneChan)
		defer stdoutPipe.Close()
		reader := bufio.NewReader(stdoutPipe)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				errChan <- []byte(fmt.Sprintf("Reader error: %v", err))
				break
			}

			line = bytes.TrimRight(line, "\r\n")
			output = append(output, line)
		}
	}()

	go func() {
		defer stderrPipe.Close()
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Bytes()
			errChan <- line
		}
	}()

	select {
	case err := <-errChan:
		_ = cmd.Process.Signal(os.Kill)
		return output, err
	case <-doneChan:
		return output, nil
	}
}

// EncodeStringToBase64 encodes a string into a Base64 byte array
func EncodeStringToBase64(inputString string) []byte {
	data := []byte(inputString)
	encodedLen := base64.StdEncoding.EncodedLen(len(data))
	encodedData := make([]byte, encodedLen)
	base64.StdEncoding.Encode(encodedData, data)
	return encodedData
}

const PackageName = "summon-conjur"
