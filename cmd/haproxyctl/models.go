package haproxyctl

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// HAProxyConfig holds the basic configuration options for haproxyctl
type HAProxyConfig struct {
	URL       url.URL
	Username  string
	Password  string
	client    *http.Client
	setupdone bool
}

func (c *HAProxyConfig) setupClient() {
	if c.setupdone {
		return
	}

	c.client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	c.setupdone = true
}

// Statistics is a slice of HAProxy Statistics
type Statistics []Statistic

// Statistic contains a set of HAProxy Statistics
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
	Downtime                Duration  `csv:"downtime"`
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

// Duration is a type that we can attach CSV marshalling to for getting time.Duration
type Duration struct {
	time.Duration
}

// UnmarshalCSV converts the seconds timestamp into a golang time.Duration
func (date *Duration) UnmarshalCSV(csv string) (err error) {
	if csv == "" {
		return nil
	}
	timeString := fmt.Sprintf("%vs", csv)
	date.Duration, err = time.ParseDuration(timeString)
	if err != nil {
		return err
	}
	return nil
}

// You could also use the standard Stringer interface
func (date *Duration) String() string {
	return time.Duration(date.Nanoseconds()).String()
}

// EntryType can be a Frontend, Backend, Server or Socket
type EntryType int

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

// Action is a set of actions that we can send to a HAProxy server
type Action string

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
