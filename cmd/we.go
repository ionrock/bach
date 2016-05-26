package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/ionrock/bach"
)

func WeBefore(c *cli.Context) error {
	bach.InitLogging(c.Bool("debug"))

	if c.String("env") != "" {
		err := bach.WithEnv(c.String("env"))
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Before = WeBefore
	app.Action = bach.CommandAction
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
		cli.BoolFlag{
			Name:  "debug, D",
			Usage: "Turn on debugging output",
		},
	}

	app.Run(os.Args)
}
