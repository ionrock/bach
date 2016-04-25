package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/ionrock/bach/core"
)

func RunScriptBefore(c *cli.Context) error {
	script := c.String("script")
	fmt.Printf("Running Script: %s\n", script)

	if script != "" {
		cmd := core.NewCmd(script)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func GetHereApp() *cli.App {
	app := cli.NewApp()
	app.Before = RunScriptBefore
	app.Action = core.CommandAction
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "script, s",
			Value:  "",
			EnvVar: "BACH_HERE_SCRIPT",
			Usage:  "A script / cmd to run when starting your process.",
		},
	}
	return app
}

func main() {
	app := GetHereApp()
	app.Run(os.Args)
}
