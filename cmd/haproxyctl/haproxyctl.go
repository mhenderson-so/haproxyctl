package haproxyctl

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gocarina/gocsv"
)

// GetStats gets the latest set of statistics from HAProxy
func (c *HAProxyConfig) GetStats() (*Statistics, error) {
	req, err := http.NewRequest("GET", c.GetRequestURI(), nil)
	if err != nil {
		return nil, err
	}
	if c.Username != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code %v", resp.StatusCode)
	}

	csvBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var theseStats Statistics
	err = gocsv.UnmarshalBytes(csvBody, &theseStats)
	if err != nil {
		return nil, err
	}

	return &theseStats, nil
}

// GetRequestURI returns the URL to be used when sending a request
func (c *HAProxyConfig) GetRequestURI() string {
	return fmt.Sprintf("%vhaproxy;csv", c.URL.String())
}

// SetCredentialsFromAuthString is used when you have credentails in an auth string, but don't want to send
// the separate username/passwords
func (c *HAProxyConfig) SetCredentialsFromAuthString(authstring string) error {
	decoded, err := base64.StdEncoding.DecodeString(authstring)
	if err != nil {
		return err
	}
	decodedParts := strings.Split(string(decoded), ":")
	if len(decodedParts) != 2 {
		return fmt.Errorf("auth string is not a username/password combination")
	}

	c.Username = decodedParts[0]
	c.Password = decodedParts[1]

	return nil
}
