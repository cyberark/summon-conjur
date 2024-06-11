package test

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strings"
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
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			line := scanner.Bytes()
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
