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
		os.Stderr.Write([]byte("A variable name or version flag must be given as the first and only argument!"))
		os.Exit(-1)
	}

	singleArgument := os.Args[1]
	switch singleArgument {
	case "-v","--version":
		os.Stdout.Write([]byte(VERSION))
	default:
		RetrieveSecret(singleArgument)
	}
}

func printAndExit(err error) {
	os.Stderr.Write([]byte(err.Error()))
	os.Exit(1)
}
