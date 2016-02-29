package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("A variable name must be given as the first and only argument!")
		os.Exit(-1)
	}
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
	os.Stderr.Write([]byte(err.Error()))
	os.Exit(1)
}
