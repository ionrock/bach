package main

import (
	"fmt"
	"os"
	"strings"

	// log "github.com/Sirupsen/logrus"
	// "github.com/codegangsta/cli"
	"github.com/ionrock/bach"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("bach")
	viper.AddConfigPath("/etc/bach/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func Do(cmd string) error {
	parts := strings.Split(cmd, " ")
	return bach.RunLogged(parts...)
}

func RunBachBefore(c *cli.Context) error {
	err := Do(viper.GetString("join.cmd"))
	if err != nil {
		return err
	}
	return nil
}

func RunBachAfter(c *cli.Context) error {
	err := Do(viper.GetString("leave.cmd"))
	if err != nil {
		return err
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Before = RunBachBefore
	app.Action = bach.CommandAction
	app.After = RunBachAfter
	app.Run(os.Args)
}
