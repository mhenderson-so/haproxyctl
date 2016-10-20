# haproxyctl

HAProxyCTL is a small Golang library for retriving the control settings
from a haproxy instance over HTTP.

It is used for querying remote HAProxy instances.

## ToDo

- Add controlling the backends themselves (up/down/drain etc)

## Example usage

```Go
package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/mhenderson-so/haproxyctl/cmd/haproxyctl"
)

func main() {
	//Create a simple URL to our HAProxy endpoint
	endpoint, _ := url.Parse("http://10.1.8.20:7003/")
	//Create our config item. We could specify our username/password here, but we will specify it encoded later on
	thisConfig := haproxyctl.HAProxyConfig{
		URL: *endpoint,
	}
	//Pass through our encoded username/password
	err := thisConfig.SetCredentialsFromAuthString("dXNlcm5hbWU6cGFzc3dvcmQ=")
	if err != nil {
		fmt.Println(err)
		return
	}

	//Get the latest statistics from this HAProxy server
	stats, err := thisConfig.GetStats()
	if err != nil {
		fmt.Println(err)
		return
	}

	//Print the output to the console in a readable format
	for _, x := range *stats {
		data := strings.Replace(fmt.Sprintf("%+v\n", x), " ", "\n\t", -1)
		fmt.Println(data[1 : len(data)-2])
		fmt.Println()
	}
}
```