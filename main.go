package main

import (
	"fmt"
	"os"
	"github.com/conjurinc/api-go/conjurapi"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("A variable name must be given as the first and only argument!")
		os.Exit(-1)
	}
	variableName := os.Args[1]

	config := conjurapi.Config{
		ApplianceUrl: os.Getenv("CONJUR_APPLIANCE_URL"),
		Account: os.Getenv("CONJUR_ACCOUNT"),
		Username: os.Getenv("CONJUR_AUTHN_LOGIN"),
		APIKey: os.Getenv("CONJUR_AUTHN_API_KEY"),
	}

	conjur := conjurapi.NewClient(config)

	value, err := conjur.RetrieveSecret(variableName)
	if err != nil {
		printAndExit(err)
	}

	fmt.Print(value)
}

func printAndExit(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}
