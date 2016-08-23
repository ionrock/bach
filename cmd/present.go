package main

import (
	"log"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/ionrock/bach"
)

func ClusterAction(c *cli.Context) error {
	name := c.Args()[0]
	ip := ""
	if len(c.Args()) > 1 {
		ip = c.Args()[1]
	}

	config := bach.LocalConfig(name)

	config.BindPort = c.Int("port")

	bach.InitializeMembership(config, ip)

	log.Printf("Listing on: %s:%d", config.BindAddr, config.BindPort)

	for {
		time.Sleep(time.Second)
	}

	return nil
}

func RunScriptBefore(c *cli.Context) error {
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
			Usage:  "Connect to / Create a cluster",
			Action: ClusterAction,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name: "port, p",
				},
			},
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
	}
	return app
}

func main() {
	app := GetHereApp()
	app.Run(os.Args)
}
