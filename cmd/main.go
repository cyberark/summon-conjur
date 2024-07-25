package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/cyberark/conjur-api-go/conjurapi/logging"
	"github.com/cyberark/summon-conjur/pkg/summon_conjur"
	"github.com/karrick/golf"
	log "github.com/sirupsen/logrus"
)

func makeSecretRetriever() (func(variableName string) ([]byte, error), error) {
	config, err := conjurapi.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("Failed loading Conjur API config: %s\n", err.Error())
	}

	conjur, err := conjurapi.NewClientFromEnvironment(config)
	if err != nil {
		return nil, fmt.Errorf("Failed creating a Conjur client: %s\n", err.Error())
	}

	return func(variableName string) ([]byte, error) {
		value, err := conjur.RetrieveSecret(variableName)
		if err != nil {
			return nil, err
		}

		return value, nil
	}, nil
}

func main() {
	var help = golf.BoolP('h', "help", false, "show help")
	var version = golf.BoolP('V', "version", false, "show version")
	var verbose = golf.BoolP('v', "verbose", false, "be verbose")

	golf.Parse()
	args := golf.Args()

	if *version {
		fmt.Println(summon_conjur.FullVersionName)
		os.Exit(0)
	}
	if *help {
		golf.Usage()
		os.Exit(0)
	}

	log.SetFormatter(&log.TextFormatter{DisableTimestamp: true, DisableLevelTruncation: true})
	if *verbose {
		log.SetLevel(log.DebugLevel)
		logging.ApiLog.SetLevel(log.DebugLevel)
	}

	retrieveSecrets, err := makeSecretRetriever()
	if err != nil {
		log.Errorf("%s", err.Error())
		os.Exit(1)
	}

	if len(args) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		// Breaking out of this loop is controlled by a parent process by sending EOF to the stdin
		for scanner.Scan() {
			variableName := scanner.Text()
			if variableName == "" {
				log.Errorln("Failed to retrieve variable from stdin")
				continue
			}
			value, err := retrieveSecrets(variableName)
			if err != nil {
				log.Errorln(err.Error())
				continue
			}
			base64Value := make([]byte, base64.StdEncoding.EncodedLen(len(value)))
			base64.StdEncoding.Encode(base64Value, value)
			fmt.Fprintln(os.Stdout, string(base64Value))
		}
		if err := scanner.Err(); err != nil {
			log.Errorln(err.Error())
			os.Exit(1)
		}
	} else {
		value, err := retrieveSecrets(args[0])
		if err != nil {
			log.Errorln(err.Error())
			os.Exit(1)
		}
		os.Stdout.Write(value)
	}
}
