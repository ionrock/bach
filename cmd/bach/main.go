package main

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
	"github.com/ionrock/bach"
	"github.com/urfave/cli"
)

var builddate = ""
var gitref = ""

func applyWithEnv(v interface{}) {
	log.Infof("%#v", v)
	items := v.([]interface{})

	entries := make([]map[string]string, len(items))

	for i, item := range items {
		for k, v := range item.(map[string]interface{}) {
			action := make(map[string]string)
			action[k] = v.(string)
			entries[i] = action
		}
	}
	envAlias := bach.EnvAlias{}
	envAlias.ApplyFromMap(entries)
}

func applyToConfig(v interface{}) {
	tc := v.(map[string]interface{})

	bach.ApplyConfig(
		tc["template"].(string),
		tc["config"].(string),
	)

}

func RunBachBefore(c *cli.Context) error {
	b, err := ioutil.ReadFile(c.String("config"))
	if err != nil {
		log.Fatal(err)
	}

	log.Info(string(b))
	conf := make([]map[string]interface{}, 0)

	yaml.Unmarshal(b, &conf)

	for _, cmd := range conf {
		for k, v := range cmd {
			switch k {
			case "withenv":
				applyWithEnv(v)

			case "toconfig":
				applyToConfig(v)
			}
		}
	}

	return nil
}

func RunBachAfter(c *cli.Context) error {
	return nil
}

func main() {
	app := cli.NewApp()
	app.Version = fmt.Sprintf("%s-%s", gitref, builddate)

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "A Bach yaml config",
		},
	}
	app.Before = RunBachBefore
	app.Action = bach.CommandAction
	app.After = RunBachAfter
	app.Run(os.Args)
}
