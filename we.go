package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/ionrock/bach/core"
)

func We(c *cli.Context) error {
	return core.WithEnv(c.String("env"))
}

func main() {
	app := cli.NewApp()
	app.Before = We
	app.Action = core.CommandAction
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "env, e",
			Value: "",
			Usage: "A YAML file to include in the environment.",
		},
		cli.StringFlag{
			Name:  "directory, d",
			Value: "",
			Usage: "A directory containing YAML files to recursively applyt to the environment.",
		},
		cli.StringFlag{
			Name:  "alias, a",
			Value: "",
			Usage: "A YAML file containing a list of file/directory entries to apply to the environment.",
		},
		cli.StringFlag{
			Name:  "envvar, E",
			Value: "",
			Usage: "Override a single environment variable.",
		},
	}

	app.Run(os.Args)
}
