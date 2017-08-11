package conjurapi

import (
	"net/http"
)

type Config struct {
	Account      string
	APIKey       string
	ApplianceUrl string
	Username     string
}

type Client struct {
	config     Config
	AuthToken  string
	httpClient *http.Client
}

func NewClient(c Config) *Client {
	return &Client{
		config:     c,
		httpClient: &http.Client{},
	}
}
