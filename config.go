package main

import (
	"fmt"
	"github.com/bgentry/go-netrc/netrc"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"os"
)

// Contains information to connect to a Conjur appliance
type Config struct {
	// The url, like "https://conjur.companyname.com/api"
	ApplianceUrl string `yaml:"appliance_url"`
	// Your conjur username (aka login)
	Username string

	// Alternate URL for authentication service
	AltAuthnUrl string `yaml:"authn_url"`

	// Alternate URL for core services
	AltCoreUrl string `yaml:"core_url"`

	// Your conjur API key (aka password)
	APIKey string
	// Path to the certificate for your Conjur appliance
	SSLCertPath string `yaml:"cert_file"`
}

func (c *Config) AuthnUrl() string {
	if c.AltAuthnUrl != "" {
		return c.AltAuthnUrl
	}

	return c.ApplianceUrl + "/authn"
}

func (c *Config) CoreUrl() string {
	if c.AltCoreUrl != "" {
		return c.AltCoreUrl
	}
	return c.ApplianceUrl
}

func mergeValue(a, b string) string {
	if len(a) != 0 {
		return a
	}
	return b
}

func (c *Config) merge(o *Config) {
	c.ApplianceUrl = mergeValue(c.ApplianceUrl, o.ApplianceUrl)
	c.Username = mergeValue(c.Username, o.Username)
	c.SSLCertPath = mergeValue(c.SSLCertPath, o.SSLCertPath)
	c.APIKey = mergeValue(c.APIKey, o.APIKey)
	c.AltAuthnUrl = mergeValue(c.AltAuthnUrl, o.AltAuthnUrl)
	c.AltCoreUrl = mergeValue(c.AltCoreUrl, o.AltCoreUrl)
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
		ApplianceUrl: os.Getenv("CONJUR_APPLIANCE_URL"),
		Username:     os.Getenv("CONJUR_AUTHN_LOGIN"),
		SSLCertPath:  os.Getenv("CONJUR_CERT_FILE"),
		APIKey:       os.Getenv("CONJUR_AUTHN_API_KEY"),
		AltCoreUrl:   os.Getenv("CONJUR_CORE_URL"),
		AltAuthnUrl:  os.Getenv("CONJUR_AUTHN_URL"),
	}

	c.merge(&env)
}

func (c *Config) mergeNetrc() {
	rc, err := netrc.ParseFile(os.ExpandEnv("$HOME/.netrc"))
	if err != nil {
		return
	}

	m := rc.FindMachine(c.ApplianceUrl + "/authn")

	if m != nil {
		c.APIKey = m.Password
		c.Username = m.Login
	}
}

func (c *Config) validate() error {
	// check urls, a bit more complex than the other stuff
	if (c.AltCoreUrl == "" || c.AltAuthnUrl == "") && c.ApplianceUrl == "" {
		return fmt.Errorf("Must specify either authn and core urls or an appliance url in %v", c)
	}

	if c.Username == "" || c.APIKey == "" || c.SSLCertPath == "" {
		return fmt.Errorf("Missing config info in %v", c)
	}
	return nil
}

// Gathers configuration information as follows:
//  * Always load /etc/conjur.conf
//  * If $CONJURRC is set, load configuration from that file
//  * Otherwise, read ~/.conjurrc and ./.conjurrc
//  * Load credentials from the environment if they are present
//  * Load them from ~/.netrc if it exists and the values are found
//  * Fail if no credentials are found.
func LoadConfig() (*Config, error) {
	c := Config{}

	// read /etc/conjur.conf
	c.mergeYAML("/etc/conjur.conf")

	// check for $CONJURRC
	conjurrc := os.Getenv("CONJURRC")

	if conjurrc != "" {
		c.mergeYAML(conjurrc)
	} else {
		// merge ~/.conjurrc and ./.conjurrc
		path := os.ExpandEnv("$HOME/.conjurrc")
		c.mergeYAML(path)

		path = os.ExpandEnv("$HOME/.conjurrc")
		c.mergeYAML(path)
	}

	c.mergeEnv()

	// merge credentials from netrc
	c.mergeNetrc()

	err := c.validate()

	if err != nil {
		return nil, err
	}

	return &c, nil
}
