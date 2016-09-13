package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/ionrock/bach"
)

func WeBefore(c *cli.Context) error {
	bach.InitLogging(c.Bool("debug"))
	log.Info("args: ", os.Args[1:])
	return bach.WithEnv(os.Args[1:])
}

func main() {
	app := cli.NewApp()
	app.Name = "we"
	app.Usage = "Add environment variables via YAML or scripts before running a command."
	app.Before = WeBefore
	app.Action = bach.CommandAction

	// NOTE: These flags are essentially ignored b/c we need ordered flags
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "env, e",
			Usage: "A YAML/JSON file to include in the environment.",
		},

		cli.StringSliceFlag{
			Name:  "script, s",
			Usage: "Execute a script that outputs YAML.",
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

		cli.StringSliceFlag{
			Name:  "envvar, E",
			Usage: "Override a single environment variable.",
		},

		cli.BoolFlag{
			Name:  "debug, D",
			Usage: "Turn on debugging output",
		},
	}

	app.Run(os.Args)
}
