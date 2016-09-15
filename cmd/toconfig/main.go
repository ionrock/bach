package main

import (
	"os"

	"github.com/ionrock/bach"
	"github.com/urfave/cli"
)

func ApplyConfigAction(c *cli.Context) error {
	return bach.ApplyConfig(c.String("template"), c.String("config"))
}

func main() {
	app := cli.NewApp()
	app.Name = "toconfig"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "template, t",
			Usage: "The template(s) to fill in from env vars",
		},
		cli.StringFlag{
			Name:  "config, c",
			Value: "",
			Usage: "Where to write the template output",
		},
	}
	app.Before = ApplyConfigAction
	app.Action = bach.CommandAction
	app.Run(os.Args)
}
