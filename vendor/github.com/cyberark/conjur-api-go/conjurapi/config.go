package conjurapi

import (
	"fmt"
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v1"
	"strings"
)

type Config struct {
	Account      string `yaml:"account"`
	ApplianceURL string `yaml:"appliance_url"`
	NetRCPath    string `yaml:"netrc_path"`
	SSLCert      string
	SSLCertPath  string `yaml:"cert_file"`
	Https        bool
	V4           bool
}

func (c *Config) validate() (error) {
	errors := []string{}

	if c.ApplianceURL == "" {
		errors = append(errors, fmt.Sprintf("Must specify an ApplianceURL in %v", c))
	}

	if c.Account == "" {
		errors = append(errors, fmt.Sprintf("Must specify an Account in %v", c))
	}

	c.Https = c.SSLCertPath != "" || c.SSLCert != ""

	if len(errors) == 0 {
		return nil
	}
	return fmt.Errorf("%s", strings.Join(errors, " -- "))
}

func (c *Config) ReadSSLCert() ([]byte, error) {
	if c.SSLCert != "" {
		return []byte(c.SSLCert), nil
	}
	return ioutil.ReadFile(c.SSLCertPath)
}

func (c *Config) BaseURL() string {
	prefix := ""
	if !strings.HasPrefix(c.ApplianceURL, "http") {
		if c.Https {
			prefix = "https://"
		} else {
			prefix = "http://"
		}
	}
	return prefix + c.ApplianceURL
}

func mergeValue(a, b string) string {
	if len(b) != 0 {
		return b
	}
	return a
}

func (c *Config) merge(o *Config) {
	c.ApplianceURL = mergeValue(c.ApplianceURL, o.ApplianceURL)
	c.Account = mergeValue(c.Account, o.Account)
	c.SSLCert = mergeValue(c.SSLCert, o.SSLCert)
	c.SSLCertPath = mergeValue(c.SSLCertPath, o.SSLCertPath)
	c.NetRCPath = mergeValue(c.NetRCPath, o.NetRCPath)
	c.V4 = c.V4 || o.V4
}

func (c *Config) mergeYAML(filename string) {
	var tmp Config

	buf, err := ioutil.ReadFile(filename)

	if err != nil {
		return
	}

	err = yaml.Unmarshal(buf, &tmp)

	if err != nil {
		return
	}

	c.merge(&tmp)
}

func (c *Config) mergeEnv() {
	env := Config{
		ApplianceURL: os.Getenv("CONJUR_APPLIANCE_URL"),
		SSLCert:      os.Getenv("CONJUR_SSL_CERTIFICATE"),
		SSLCertPath:  os.Getenv("CONJUR_CERT_FILE"),
		Account:      os.Getenv("CONJUR_ACCOUNT"),
		NetRCPath:    os.Getenv("CONJUR_NETRC_PATH"),
		V4:           os.Getenv("CONJUR_MAJOR_VERSION") == "4",
	}

	c.merge(&env)
}

func LoadConfig() (Config) {
	config := Config{}

	config.mergeYAML("/etc/conjur.conf")

	conjurrc := os.Getenv("CONJURRC")

	if conjurrc != "" {
		config.mergeYAML(conjurrc)
	} else {
		path := os.ExpandEnv("$HOME/.conjurrc")
		config.mergeYAML(path)

		path = os.ExpandEnv("$PWD/.conjurrc")
		config.mergeYAML(path)
	}

	config.mergeEnv()

	return config
}
