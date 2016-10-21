package haproxyctl

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gocarina/gocsv"
)

// GetStats gets the latest set of statistics from HAProxy
func (c *HAProxyConfig) GetStats() (*Statistics, error) {
	req, err := http.NewRequest("GET", c.GetRequestURI(true), nil)
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
func (c *HAProxyConfig) GetRequestURI(csv bool) string {
	c.setupClient()
	if csv {
		return fmt.Sprintf("%vhaproxy;csv", c.URL.String())
	}
	return fmt.Sprintf("%vhaproxy", c.URL.String())
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

// SendAction sends an action to HAProxy to perform on a list of servers. For example, putting servers into
// maintenance mode, or disabling health checks. The return parameters indicate whether the request was serviced,
// whether everything was OK, and any resulting errors. For example, a request may have been applied to some of the
// nodes requested, but not others. In which case "done" will be true, but "allok" will be false, and the error will
// contain a brief text.
func (c *HAProxyConfig) SendAction(servers []string, backend string, action Action) (done bool, allok bool, err error) {

	//Build our form that we're going to POST to HAProxy
	var POSTData []string
	for _, s := range servers {
		POSTData = append(POSTData, fmt.Sprintf("s=%v", url.QueryEscape(s)))
	}
	POSTData = append(POSTData, fmt.Sprintf("action=%v", url.QueryEscape(string(action))))
	POSTData = append(POSTData, fmt.Sprintf("b=%v", url.QueryEscape(string(backend))))

	POSTBuffer := bytes.NewBufferString(strings.Join(POSTData, "&"))

	//Create our request
	req, err := http.NewRequest("POST", c.GetRequestURI(false), POSTBuffer)
	if err != nil {
		return false, false, err
	}
	if c.Username != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	//Send our request to HAProxy
	response, err := c.client.Do(req)
	if err != nil {
		return false, false, err
	}

	//We are expecting a 303 SEE OTHER response
	if response.StatusCode != 303 {
		return false, false, fmt.Errorf("status code %v", response.StatusCode)
	}

	//To see if we were successful, look at the redirection header (remember this is basically screen scraping)
	responseParts := strings.Split(response.Header.Get("Location"), "=")
	if len(responseParts) != 2 {
		return false, false, fmt.Errorf("unrecognised response: %v", response.Header.Get("Location"))
	}

	//These are our "OK" responses, where at least something was applied
	switch responseParts[1] {
	case "DONE":
		return true, true, nil
	case "PART":
		return true, false, fmt.Errorf("partially applied")
	case "NONE":
		return true, false, fmt.Errorf("no changes were applied")
	}

	return false, false, fmt.Errorf("haproxy response: %v", responseParts[1])
}
