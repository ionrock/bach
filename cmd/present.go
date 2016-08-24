package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/ionrock/bach"
)

func RunScriptBefore(c *cli.Context) error {
	return bach.RunScript(c.String("start"))
}

func RunScriptAfter(c *cli.Context) error {
	return bach.RunScript(c.String("after"))
}

func GetHereApp() *cli.App {
	app := cli.NewApp()

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
