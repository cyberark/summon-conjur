package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
)

func NewConjurHTTPClient(certPath string) (*http.Client, error) {
	pem, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM(pem)
	if !ok {
		return nil, fmt.Errorf("Can't append cert %s", certPath)
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: pool},
	}
	return &http.Client{Transport: tr}, nil
}
