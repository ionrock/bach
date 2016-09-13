package bach

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type Action interface {
	Apply() map[string]string
}

type EnvVar struct {
	field string
}

func (e EnvVar) Apply() (env map[string]string) {
	parts := strings.Split(e.field, "=")
	if len(parts) != 2 {
		log.Fatal("Invalid env var format. Use %s=%s")
	}
	key := parts[0]
	value := parts[1]
	env[key] = value

	log.Debugf("export %s = %s", key, value)
	err := os.Setenv(key, value)
	if err != nil {
		log.Fatal(err)		
	}
	return env
}

type EnvFile struct {
	path string
}

func (e EnvFile) Parse() (map[string]string, error) {
	f := Flattener{e.path}

	env, err := f.Flatten()
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range env {
		log.Debugf("%s = %s", k, v)
	}

	return env, nil
}

func (e EnvFile) Apply() map[string]string {
	env, err := e.Parse()
	if err != nil {
		panic(err)
	}

	for k, v := range env {
		log.Debugf("export %s=%s", k, v)
		err = os.Setenv(k, v)
		if err != nil {
			panic(err)
		}
	}

	return env
}

type EnvScript struct {
	cmd string
}

func (e EnvScript) Apply() map[string]string {
	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmp.Name())

	parts := SplitCommand(e.cmd)

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = tmp
	err = cmd.Run()

	if err != nil {
		panic(err)
	}

	tmp.Close()

	ef := EnvFile{path: tmp.Name()}
	return ef.Apply()
}

func WithEnv(args []string) error {
	in_flag := ""
	for _, f := range args {
		log.Debugf("default: %s", f)

		switch {
		case f == "--env" || f == "-e":
			in_flag = "env"
		case in_flag == "env":
			log.Debug("Applying env: ", f)
			action := EnvFile{path: f}
			action.Apply()
			in_flag = ""

		case f == "--script" || f == "-s":
			in_flag = "script"
		case in_flag == "script":
			log.Debug("Applying script: ", f)
			action := EnvScript{cmd: f}
			action.Apply()
			in_flag = ""

		case f == "--envvar" || f == "-E":
			in_flag = "envvar"
		case in_flag == "envvar":
			log.Debug("Applying single var: ", f)
			action := EnvVar{field: f}
			action.Apply()
			in_flag = ""			

		default:
			break
		}
	}

	return nil
}
