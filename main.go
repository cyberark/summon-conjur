package main

import (
	"fmt"
	"os"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/karrick/golf"
	log "github.com/sirupsen/logrus"
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
	var help = golf.BoolP('h', "help", false, "show help")
	var version = golf.BoolP('V', "version", false, "show version")
	var verbose = golf.BoolP('v', "verbose", false, "be verbose")

	golf.Parse()

	args := golf.Args()
	if len(args) == 0 || *help {
		golf.Usage()
		os.Exit(1)
	}

	if *verbose {
		log.SetFormatter(&log.TextFormatter{})
		log.SetLevel(log.DebugLevel)
	}

	if !*version {
		RetrieveSecret(args[0])
	} else {
		fmt.Println(VERSION)
	}
}

func printAndExit(err error) {
	os.Stderr.Write([]byte(err.Error()))
	os.Exit(1)
}
