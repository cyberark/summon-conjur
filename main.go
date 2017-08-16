package main

import (
	"fmt"
	"os"
	"github.com/cyberark/conjur-api-go/conjurapi"
)

func RetrieveSecret(variableName string) {
	config := &conjurapi.Config{}
	conjurapi.LoadFromEnv(config)

	var (
		conjur *conjurapi.Client
		err error
	)

	if TokenFile := os.Getenv("CONJUR_AUTHN_TOKEN_FILE"); TokenFile != "" {
		conjur, err = conjurapi.NewClientFromTokenFile(*config, TokenFile)
	} else if Login, APIKey := os.Getenv("CONJUR_AUTHN_LOGIN"), os.Getenv("CONJUR_AUTHN_API_KEY"); Login != "" && APIKey != "" {
		conjur, err = conjurapi.NewClientFromKey(*config, Login, APIKey)
	} else {
		fmt.Println("Environment variables satisfying at least one authentication strategy must be provided!")
		os.Exit(-1)
	}

	if err != nil {
		printAndExit(err)
	}

	value, err := conjur.RetrieveSecret(variableName)
	if err != nil {
		printAndExit(err)
	}

	fmt.Print(value)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("A variable name must be given as the first and only argument!")
		os.Exit(-1)
	}
	variableName := os.Args[1]

	RetrieveSecret(variableName)
}

func printAndExit(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}
