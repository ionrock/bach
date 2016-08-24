package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/ionrock/bach"
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

func RunScriptBefore(c *cli.Context) error {
	if c.String("name") != "" && c.String("cluster-addr") != "" {
		config := bach.LocalConfig(c.String("name"))
		sm := &bach.ServiceMap{
			Config:      config,
			ClusterAddr: c.String("cluster-addr"),
		}
		sm.Join()
	}

	return bach.RunScript(c.String("start"))
}

func RunScriptAfter(c *cli.Context) error {
	return bach.RunScript(c.String("after"))
}

func GetHereApp() *cli.App {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		cli.Command{
			Name:   "cluster",
			Usage:  "Manually connect to / create a cluster",
			Action: ClusterAction,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name: "port, p",
				},
			},
		},
		cli.Command{
			Name:   "nodes",
			Usage:  "Connect to a cluster and echo a JSON list of members",
			Action: ClusterNodesAction,
		},
	}

	app.Before = RunScriptBefore
	app.After = RunScriptAfter
	app.Action = bach.CommandAction
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "start, s",
			Value:  "",
			EnvVar: "BACH_PRESENT_START",
			Usage:  "A script / cmd to run when starting your process.",
		},
		cli.StringFlag{
			Name:   "after, a",
			Value:  "",
			EnvVar: "BACH_PRESENT_AFTER",
			Usage:  "A script / cmd to run when finishing your process.",
		},
		cli.StringFlag{
			Name:  "name, n",
			Value: "",
			Usage: "The name of an app to a cluster",
		},
		cli.StringFlag{
			Name:  "cluster-addr, c",
			Value: "",
			Usage: "The address of another node in the cluster",
		},
	}
	return app
}

func main() {
	app := GetHereApp()
	app.Run(os.Args)
}
