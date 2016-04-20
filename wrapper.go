package main

import (
	"os"
	"os/exec"

	"github.com/codegangsta/cli"
)

func NewCmd(parts ...string) *exec.Cmd {
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func CommandAction(c *cli.Context) {
	cmd := NewCmd(c.Args()...)
	cmd.Run()
}
