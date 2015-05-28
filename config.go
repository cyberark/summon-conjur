package main
import (
    "gopkg.in/yaml.v1"
    "github.com/bgentry/go-netrc/netrc"
    "io/ioutil"
    "os"
    "fmt"
)

// Contains information to connect to a Conjur appliance
type Config struct {
    // The url, like "https://conjur.companyname.com/api"
    ApplianceUrl string `yaml:"appliance_url"`
    // Your conjur username (aka login)
    Username     string
    // Your conjur API key (aka password)
    APIKey       string
    // Path to the certificate for your Conjur appliance
    SSLCertPath  string `yaml:"cert_file"`
}

func mergeValue(a,b string) string {
    if len(a) != 0 {
        return a
    }
    return b
}


func (c *Config) merge(o *Config) {
    c.ApplianceUrl = mergeValue(c.ApplianceUrl, o.ApplianceUrl)
    c.Username     = mergeValue(c.Username, o.Username)
    c.SSLCertPath  = mergeValue(c.SSLCertPath, o.SSLCertPath)
    c.APIKey       = mergeValue(c.APIKey, o.APIKey)
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
    env := Config {
        ApplianceUrl: os.Getenv("CONJUR_APPLIANCE_URL"),
        Username: os.Getenv("CONJUR_AUTHN_LOGIN"),
        SSLCertPath: os.Getenv("CONJUR_CERT_FILE"),
        APIKey: os.Getenv("CONJUR_API_KEY"),
    }

    c.merge(&env)
}

func (c *Config) mergeNetrc() {

    rc, err := netrc.ParseFile(os.ExpandEnv("$HOME/.netrc"))
    if err != nil {
        return
    }

    m := rc.FindMachine(c.ApplianceUrl)

    if m != nil {
        c.APIKey = m.Password
        c.Username = m.Login
    }
}


func (c *Config) validate() error {
    if c.ApplianceUrl == "" ||
        c.Username    == "" ||
        c.APIKey      == "" ||
        c.SSLCertPath == "" {
            return fmt.Errorf("Missing config info in %V", c)
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
    }else{
        // merge ~/.conjurrc and ./.conjurrc
        path := os.ExpandEnv("$HOME/.conjurrc")
        c.mergeYAML(path)

        path = os.ExpandEnv("$HOME/.conjurrc")
        c.mergeYAML(path)
    }

    // merge credentials from netrc
    c.mergeNetrc()

    // merge in the environment
    c.mergeEnv()

    err := c.validate()

    if err != nil {
        return nil, err
    }

    return &c, nil
}

