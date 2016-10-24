package main

import (
	"net/url"

	"github.com/mhenderson-so/haproxyctl/cmd/haproxyctl"
)

type HAProxyCtlConfig struct {
	DefaultUsername string
	DefaultPassword string
	LoadBalancers   []LoadBalancer
}

type LoadBalancer struct {
	Name       string
	Url        string
	Username   string
	Password   string
	HAProxyCtl haproxyctl.HAProxyConfig
}

func (c *HAProxyCtlConfig) ProcessInit() {
	for i, x := range c.LoadBalancers {
		thisUsername := x.Username
		if thisUsername == "" {
			thisUsername = c.DefaultUsername
		}
		thisPassword := x.Password
		if thisPassword == "" {
			thisPassword = c.DefaultPassword
		}
		thisURL, _ := url.Parse(x.Url)
		c.LoadBalancers[i].HAProxyCtl = haproxyctl.HAProxyConfig{
			Username: thisUsername,
			Password: thisPassword,
			URL:      *thisURL,
		}
	}
}

const (
	ActionGetDetail = "get"
)
