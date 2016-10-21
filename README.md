# haproxyctl

HAProxyCTL is a small Golang library for retriving the control settings
from a haproxy instance over HTTP.

It is used for querying remote HAProxy instances, and can send server actions such as putting 
a server into maintenance, or disabling health checks.

For details about usage, including commands, see the [GoDoc documentation below](#godoc-documentation-haproxyctl).

<!-- TOC -->

- [haproxyctl](#haproxyctl)
    - [Example usage](#example-usage)
        - [Discovering HAProxy statistics](#discovering-haproxy-statistics)
        - [Performing a HAProxy action command](#performing-a-haproxy-action-command)
- [GoDoc documentation (haproxyctl)](#godoc-documentation-haproxyctl)
    - [Usage](#usage)
            - [type Action](#type-action)
            - [type Duration](#type-duration)
            - [func (*Duration) String](#func-duration-string)
            - [func (*Duration) UnmarshalCSV](#func-duration-unmarshalcsv)
            - [type EntryType](#type-entrytype)
            - [type HAProxyConfig](#type-haproxyconfig)
            - [func (*HAProxyConfig) GetRequestURI](#func-haproxyconfig-getrequesturi)
            - [func (*HAProxyConfig) GetStats](#func-haproxyconfig-getstats)
            - [func (*HAProxyConfig) SendAction](#func-haproxyconfig-sendaction)
            - [func (*HAProxyConfig) SetCredentialsFromAuthString](#func-haproxyconfig-setcredentialsfromauthstring)
            - [type Statistic](#type-statistic)
            - [type Statistics](#type-statistics)

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

# GoDoc documentation (haproxyctl)
--
    import "github.com/mhenderson-so/haproxyctl/cmd/haproxyctl"


## Usage

#### type Action

```go
type Action string
```

Action is a set of actions that we can send to a HAProxy server

```go
const (
	ActionSetStateToReady     Action = "ready"
	ActionSetStateToDrain     Action = "drain"
	ActionSetStateToMaint     Action = "maint"
	ActionHealthDisableChecks Action = "dhlth"
	ActionHealthEnableChecks  Action = "ehlth"
	ActionHealthForceUp       Action = "hrunn"
	ActionHealthForceNoLB     Action = "hnolb"
	ActionHealthForceDown     Action = "hdown"
	ActionAgentDisablechecks  Action = "dagent"
	ActionAgentEnablechecks   Action = "eagent"
	ActionAgentForceUp        Action = "arunn"
	ActionAgentForceDown      Action = "adown"
	ActionKillSessions        Action = "shutdown"
)
```

#### type Duration

```go
type Duration struct {
	time.Duration
}
```

Duration is a type that we can attach CSV marshalling to for getting
time.Duration

#### func (*Duration) String

```go
func (date *Duration) String() string
```
You could also use the standard Stringer interface

#### func (*Duration) UnmarshalCSV

```go
func (date *Duration) UnmarshalCSV(csv string) (err error)
```
UnmarshalCSV converts the seconds timestamp into a golang time.Duration

#### type EntryType

```go
type EntryType int
```

EntryType can be a Frontend, Backend, Server or Socket

```go
const (
	// Frontend indicates this is a front-end
	Frontend EntryType = iota
	// Backend indicates this is a back-end
	Backend
	// Server indicates this is a server
	Server
	// Socket indicates this is a socket
	Socket
)
```

#### type HAProxyConfig

```go
type HAProxyConfig struct {
	URL      url.URL
	Username string
	Password string
}
```

HAProxyConfig holds the basic configuration options for haproxyctl

#### func (*HAProxyConfig) GetRequestURI

```go
func (c *HAProxyConfig) GetRequestURI(csv bool) string
```
GetRequestURI returns the URL to be used when sending a request

#### func (*HAProxyConfig) GetStats

```go
func (c *HAProxyConfig) GetStats() (*Statistics, error)
```
GetStats gets the latest set of statistics from HAProxy

#### func (*HAProxyConfig) SendAction

```go
func (c *HAProxyConfig) SendAction(servers []string, backend string, action Action) (done bool, allok bool, err error)
```
SendAction sends an action to HAProxy to perform on a list of servers. For
example, putting servers into maintenance mode, or disabling health checks. The
return parameters indicate whether the request was serviced, whether everything
was OK, and any resulting errors. For example, a request may have been applied
to some of the nodes requested, but not others. In which case "done" will be
true, but "allok" will be false, and the error will contain a brief text.

#### func (*HAProxyConfig) SetCredentialsFromAuthString

```go
func (c *HAProxyConfig) SetCredentialsFromAuthString(authstring string) error
```
SetCredentialsFromAuthString is used when you have credentails in an auth
string, but don't want to send the separate username/passwords

#### type Statistic

```go
type Statistic struct {
	BackendName             string    `csv:"# pxname"`
	FrontendName            string    `csv:"svname"`
	QueueCurrent            uint64    `csv:"qcur"`
	QueueMax                uint64    `csv:"qmax"`
	SessionsCurrent         uint64    `csv:"scur"`
	SessionsMax             uint64    `csv:"smax"`
	SessionLimit            uint64    `csv:"slim"`
	SessionsTotal           uint64    `csv:"stot"`
	BytesIn                 uint64    `csv:"bin"`
	BytesOut                uint64    `csv:"bout"`
	DeniedRequests          uint64    `csv:"dreq"`
	DeniedResponses         uint64    `csv:"dresp"`
	ErrorsRequests          uint64    `csv:"ereq"`
	ErrorsConnections       uint64    `csv:"econ"`
	ErrorsResponses         uint64    `csv:"eresp"`
	WarningsRetries         uint64    `csv:"wretr"`
	WarningsDispatches      uint64    `csv:"wredis"`
	Status                  string    `csv:"status"`
	Weight                  uint64    `csv:"weight"`
	IsActive                uint64    `csv:"act"`
	IsBackup                uint64    `csv:"bck"`
	CheckFailed             uint64    `csv:"chkfail"`
	CheckDowned             uint64    `csv:"chkdown"`
	StatusLastChanged       Duration  `csv:"lastchg"`
	Downtime                uint64    `csv:"downtime"`
	QueueLimit              uint64    `csv:"qlimit"`
	ProcessID               uint64    `csv:"pid"`
	ProxyID                 uint64    `csv:"iid"`
	ServiceID               uint64    `csv:"sid"`
	Throttle                uint64    `csv:"throttle"`
	LBTotal                 uint64    `csv:"lbtot"`
	Tracked                 uint64    `csv:"tracked"`
	Type                    EntryType `csv:"type"`
	Rate                    uint64    `csv:"rate"`
	RateLimit               uint64    `csv:"rate_lim"`
	RateMax                 uint64    `csv:"rate_max"`
	CheckStatus             string    `csv:"check_status"`
	CheckCode               string    `csv:"check_code"`
	CheckDuration           uint64    `csv:"check_duration"`
	HTTPResponse1xx         uint64    `csv:"hrsp_1xx"`
	HTTPResponse2xx         uint64    `csv:"hrsp_2xx"`
	HTTPResponse3xx         uint64    `csv:"hrsp_3xx"`
	HTTPResponse4xx         uint64    `csv:"hrsp_4xx"`
	HTTPResponse5xx         uint64    `csv:"hrsp_5xx"`
	HTTPResponseOther       uint64    `csv:"hrsp_other"`
	CheckFailedDets         uint64    `csv:"hanafail"`
	RequestRate             uint64    `csv:"req_rate"`
	RequestRateMax          uint64    `csv:"req_rate_max"`
	RequestTotal            uint64    `csv:"req_tot"`
	AbortedByClient         uint64    `csv:"cli_abrt"`
	AbortedByServer         uint64    `csv:"srv_abrt"`
	CompressedBytesIn       uint64    `csv:"comp_in"`
	CompressedBytesOut      uint64    `csv:"comp_out"`
	CompressedBytesBypassed uint64    `csv:"comp_byp"`
	CompressedResponses     uint64    `csv:"comp_rsp"`
	LastSession             Duration  `csv:"lastsess"`
	LastCheck               string    `csv:"last_chk"`
	LastAgentCheck          string    `csv:"last_agt"`
	AvgQueueTime            uint64    `csv:"qtime"`
	AvgConnectTime          uint64    `csv:"ctime"`
	AvgResponseTime         uint64    `csv:"rtime"`
	AvgTotalTime            uint64    `csv:"ttime"`
}
```

Statistic contains a set of HAProxy Statistics

#### type Statistics

```go
type Statistics []Statistic
```

Statistics is a slice of HAProxy Statistics
