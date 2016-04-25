package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/codegangsta/cli"
)

func CleanScript(filename string) {
	err := os.Remove(filename)
	if err != nil {
		fmt.Print(err)
	}
}

func CreateTestScript() string {
	tmp, err := ioutil.TempFile("", "test_script.sh")
	if err != nil {
		panic(err)
	}

	defer CleanScript(tmp.Name())

	script := "#!/bin/bash\n\necho 'TEST SCRIPT'\n"
	tmp.WriteString(script)
	tmp.Close()
	return tmp.Name()
}

func TestHere_WithScript(t *testing.T) {
	called := false
	script := CreateTestScript()

	app := GetHereApp()
	OrigBefore := app.Before
	app.Before = func(c *cli.Context) error {
		OrigBefore(c)
		called = true
		return nil
	}
	app.Run([]string{"-s", script, "echo", "DONE"})
}
