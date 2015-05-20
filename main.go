package main

import (
	"fmt"
	"os"
)

func main() {
	variableName := os.Args[1]

	conjur, err := NewConjurClient()
	if err != nil {
		printAndExit(err)
	}

	value, err := conjur.RetrieveVariable(variableName)
	if err != nil {
		printAndExit(err)
	}

	fmt.Print(value)
}

func printAndExit(err error) {
	fmt.Print(err.Error())
	os.Exit(1)
}
