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
		log.Errorf("Failed creating a Conjur client: %s\n", err.Error())
		os.Exit(1)
	}

	value, err := conjur.RetrieveSecret(variableName)
	if err != nil {
		log.Errorln(err.Error())
		os.Exit(1)
	}

	os.Stdout.Write([]byte(value))

}

func main() {
	var help = golf.BoolP('h', "help", false, "show help")
	var version = golf.BoolP('V', "version", false, "show version")
	var verbose = golf.BoolP('v', "verbose", false, "be verbose")

	golf.Parse()
	args := golf.Args()

	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}
	if *help {
		golf.Usage()
		os.Exit(0)
	}
	if len(args) == 0 {
		golf.Usage()
		os.Exit(1)
	}

	log.SetFormatter(&log.TextFormatter{DisableTimestamp: true, DisableLevelTruncation: true})
	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	RetrieveSecret(args[0])
}
