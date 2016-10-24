package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/mhenderson-so/haproxyctl/cmd/haproxyctl"
	"github.com/olekukonko/tablewriter"
)

var (
	tomlLoc = flag.String("config", "config.toml", "config.toml location")
)

func main() {
	flag.Parse()

	var Config HAProxyCtlConfig

	if *tomlLoc != "" {
		if _, err := toml.DecodeFile(*tomlLoc, &Config); err != nil {
			log.Fatal(err)
		}
	}

	args := flag.Args()

	if len(args) != 1 && len(args) != 3 {
		printHelp()
		log.Fatal("Invalid number of arguments, must specify one or three arguments")
		return
	}

	var argCommand haproxyctl.Action
	var argServerName string
	var argBackendName string

	argCommand = haproxyctl.Action(strings.ToLower(args[0]))

	if len(args) == 3 {
		argServerName = strings.ToLower(args[1])
		argBackendName = strings.ToLower(args[2])
	}

	if argCommand == "" {
		printHelp()
		log.Fatal("Cannot specify a blank command")
		return
	}

	if !(argCommand == ActionGetDetail ||
		argCommand == haproxyctl.ActionSetStateToReady ||
		argCommand == haproxyctl.ActionSetStateToDrain ||
		argCommand == haproxyctl.ActionSetStateToMaint ||
		argCommand == haproxyctl.ActionHealthDisableChecks ||
		argCommand == haproxyctl.ActionHealthEnableChecks ||
		argCommand == haproxyctl.ActionHealthForceUp ||
		argCommand == haproxyctl.ActionHealthForceNoLB ||
		argCommand == haproxyctl.ActionHealthForceDown ||
		argCommand == haproxyctl.ActionAgentDisablechecks ||
		argCommand == haproxyctl.ActionAgentEnablechecks ||
		argCommand == haproxyctl.ActionAgentForceUp ||
		argCommand == haproxyctl.ActionAgentForceDown ||
		argCommand == haproxyctl.ActionKillSessions) {
		printHelp()
		log.Fatal(fmt.Sprintf("Invalid command specified (%v)", argCommand))
		return
	}

	if argServerName == "" && argCommand != ActionGetDetail {
		printHelp()
		log.Fatal(fmt.Sprintf("You must specify at least one server name when using the '%v' command", argCommand))
		return
	}

	if argBackendName == "" && argCommand != ActionGetDetail {
		printHelp()
		log.Fatal(fmt.Sprintf("You must specify a backend when using the '%v' command", argCommand))
		return
	}

	argServers := strings.Split(argServerName, ",")

	Config.ProcessInit()
	//fmt.Println("Load balancers found:", len(*Config.LoadBalancers))

	var outputTable *tablewriter.Table
	if argCommand == ActionGetDetail {
		outputTable = Config.getDetails(argServerName)
	} else {
		outputTable = Config.sendAction(argCommand, argBackendName, argServers)
	}

	outputTable.Render()

}

func (c *HAProxyCtlConfig) sendAction(action haproxyctl.Action, backend string, servers []string) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"LoadBalancer", "Done", "All OK", "Error"})
	for _, h := range c.LoadBalancers {
		done, ok, err := h.HAProxyCtl.SendAction(servers, backend, action)
		table.Append([]string{
			h.Name,
			fmt.Sprintf("%v", done),
			fmt.Sprintf("%v", ok),
			fmt.Sprintf("%v", err),
		})
	}
	return table
}

func (c *HAProxyCtlConfig) getDetails(server string) *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"LoadBalancer", "Backend", "Server", "Status", "LastCheck", "Downtime", "Error"})

	for _, h := range c.LoadBalancers {
		stats, err := h.HAProxyCtl.GetStats()
		if err != nil {
			table.Append([]string{
				h.Name,
				"",
				"",
				"ERROR",
				"",
				"",
				fmt.Sprintf("%v", err),
			})
			continue
		}

		for _, s := range *stats {
			if s.FrontendName == "FRONTEND" || s.FrontendName == "BACKEND" {
				continue
			}
			table.Append([]string{
				h.Name,
				s.BackendName,
				s.FrontendName,
				s.Status,
				s.LastCheck,
				s.Downtime.String(),
				"",
			})
		}
	}

	return table
}

func printHelp() {
	fmt.Println()
	fmt.Println("HAPROXYCTL HELP")
	fmt.Println("haproxyctl is a command-line utility to the haproxyctl library.")
	fmt.Println()
	fmt.Println("It is used for interacting with haproxy servers via their web admin interface.")
	fmt.Println()
	fmt.Println("Usage: haproxyctl [-config config.toml] action server1,server2 backend")
	fmt.Println("    -config config.toml - Optional parameter to the configuration file for your haproxy nodes")
	fmt.Println("    action - the action to perform (see below for valid actions)")
	fmt.Println("    server1,server2 - A comma-seperated list of back-end servers to perform the action on")
	fmt.Println("    backend - The name of the backend to apply the action to")
	fmt.Println()
	fmt.Println("Example: haproxyctl get")
	fmt.Println("Example: haproxyctl ready ny-web01,ny-web02 prod-web")
	fmt.Println()
	fmt.Println("Valid actions are:")
	fmt.Println("    get      - Gets the status of the backends. No additional arguments are required")
	fmt.Println("    ready    - Sets the server state to 'ready'")
	fmt.Println("    drain    - Sets the server state to 'drain")
	fmt.Println("    maint    - Sets the server state to 'maintenance'")
	fmt.Println("    dhlth    - Disables health checks")
	fmt.Println("    ehlth    - Enables health checks")
	fmt.Println("    hrunn    - Forces the server to be UP")
	fmt.Println("    hnolb    - Forces the server to disable load balancing")
	fmt.Println("    hdown    - Forces the server to be DOWN")
	fmt.Println("    dagent   - Disables agent checks")
	fmt.Println("    eagent   - Enables agent checks")
	fmt.Println("    arunn    - Forces agent to be UP")
	fmt.Println("    adown    - Forces agent to be DOWN")
	fmt.Println("    shutdown - Kills all sessions")
	fmt.Println()
}
