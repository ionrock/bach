package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/ionrock/bach"
)

func RunScriptBefore(c *cli.Context) error {
	script := c.String("script")
	fmt.Printf("Running Script: %s\n", script)

	if script != "" {
		fmt.Println(script)
		cmd := bach.NewCommand(script)
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
	app.Action = bach.CommandAction
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
