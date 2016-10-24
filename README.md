# haproxyctl

HAProxyCTL is a small Golang library for retriving the control settings
from a haproxy instance over HTTP (not over a Unix socket, as other projects do), using
HAProxy's built-in stats web interface.

It is used for querying remote HAProxy instances, and can send server actions such as putting 
a server into maintenance, or disabling health checks.

For details about usage, including commands, see the [GoDoc documentation](https://godoc.org/github.com/mhenderson-so/haproxyctl/cmd/haproxyctl).

<!-- TOC -->

- [haproxyctl](#haproxyctl)
    - [Example usage](#example-usage)
        - [Discovering HAProxy statistics](#discovering-haproxy-statistics)
        - [Performing a HAProxy action command](#performing-a-haproxy-action-command)
    - [Example program](#example-program)

<!-- /TOC -->

## Example usage

### Discovering HAProxy statistics

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

### Performing a HAProxy action command

```Go
package main

import (
	"fmt"
	"net/url"

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

	//Put these two servers into the "ready" state
	servers := []string{
		"ny-web01",
		"ny-web02",
	}
	done, _, err := thisConfig.SendAction(servers, "prod_web_tier", haproxyctl.ActionSetStateToReady)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !done {
		fmt.Println("not ok")
		return
	}
	fmt.Println("ok")
}
```

## Example program

There is a small example program contained in the root directory that can perform the
haproxy commands exposed by haproxyctl. You can build with `go build`. You will need to edit
the supplied example `config.toml` with your haproxy environment.

```
Usage: haproxyctl [-config config.toml] action server1,server2 backend
    -config config.toml - Optional parameter to the configuration file for your haproxy nodes
    action - the action to perform (see below for valid actions)
    server1,server2 - A comma-seperated list of back-end servers to perform the action on
    backend - The name of the backend to apply the action to

Example: haproxyctl get
Example: haproxyctl ready ny-web01,ny-web02 prod-web

Valid actions are:
    get      - Gets the status of the backends. No additional arguments are required
    ready    - Sets the server state to 'ready'
    drain    - Sets the server state to 'drain
    maint    - Sets the server state to 'maintenance'
    dhlth    - Disables health checks
    ehlth    - Enables health checks
    hrunn    - Forces the server to be UP
    hnolb    - Forces the server to disable load balancing
    hdown    - Forces the server to be DOWN
    dagent   - Disables agent checks
    eagent   - Enables agent checks
    arunn    - Forces agent to be UP
    adown    - Forces agent to be DOWN
    shutdown - Kills all sessions
```