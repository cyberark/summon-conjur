package main

import (
	"os"
	"github.com/cyberark/conjur-api-go/conjurapi"
)

func RetrieveSecret(variableName string) {
	config := conjurapi.LoadConfig()

	conjur, err := conjurapi.NewClientFromEnvironment(config)

	if err != nil {
		printAndExit(err)
	}

	value, err := conjur.RetrieveSecret(variableName)
	if err != nil {
		printAndExit(err)
	}

	os.Stdout.Write([]byte(value))

}

func main() {
	if len(os.Args) != 2 {
		os.Stderr.Write([]byte("A variable name must be given as the first and only argument!"))
		os.Exit(-1)
	}
	variableName := os.Args[1]

	RetrieveSecret(variableName)
}

func printAndExit(err error) {
	os.Stderr.Write([]byte(err.Error()))
	os.Exit(1)
}
