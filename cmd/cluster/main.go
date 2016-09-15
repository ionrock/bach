package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/ionrock/bach"
	"github.com/urfave/cli"
)

func ClusterAction(c *cli.Context) error {
	name := c.Args()[0]
	addr := ""
	if len(c.Args()) > 1 {
		addr = c.Args()[1]
	}

	config := bach.LocalConfig(name)
	if c.Int("port") != 0 {
		config.BindPort = c.Int("port")
	}

	sm := &bach.ServiceMap{
		Config:      config,
		ClusterAddr: addr,
	}

	sm.Load()
	sm.Join()

	log.Printf("Listing on: %s:%d", config.BindAddr, config.BindPort)
	for {
		time.Sleep(time.Second)
	}
	return nil
}

func ClusterNodesAction(c *cli.Context) error {
	args := c.Args()

	if len(args) == 0 {
		log.Fatal("No cluster address provided")
	}
	addr := args[0]

	config := bach.LocalConfig("tmp")

	config.LogOutput = ioutil.Discard

	sm := &bach.ServiceMap{
		Config:      config,
		ClusterAddr: addr,
	}

	sm.Join()
	sm.CopyJsonTo(os.Stdout)

	return nil
}

var sm *bach.ServiceMap

func JoinClusterBefore(c *cli.Context) error {
	if c.String("name") != "" && c.String("cluster-addr") != "" {
		config := bach.LocalConfig(c.String("name"))
		sm = &bach.ServiceMap{
			Config:      config,
			ClusterAddr: c.String("cluster-addr"),
		}
		sm.Join()
	}

	return nil
}

func LeaveClusterAfter(c *cli.Context) error {
	if sm != nil {
		err := sm.Leave()

		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func GetClusterApp() *cli.App {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		cli.Command{
			Name:   "start",
			Usage:  "Manually connect to / create a cluster",
			Action: ClusterAction,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name: "port, p",
				},
			},
		},
		cli.Command{
			Name:      "nodes",
			Usage:     "Connect to a cluster and echo a JSON list of members",
			ArgsUsage: "[CLUSTER_ADDRESS]",
			Action:    ClusterNodesAction,
		},
	}

	app.Before = JoinClusterBefore
	app.After = LeaveClusterAfter
	app.Action = bach.CommandAction
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "name, n",
			Value:  "",
			EnvVar: "BACH_CLUSTER_NAME",
			Usage:  "The name of an app to a cluster",
		},
		cli.StringFlag{
			Name:   "cluster-addr, c",
			Value:  "",
			EnvVar: "BACH_CLUSTER_ADDR",
			Usage:  "The address of another node in the cluster",
		},
	}
	return app
}

func main() {
	app := GetClusterApp()
	app.Run(os.Args)
}
