package bach

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func loadYaml(path string) ([]map[interface{}]interface{}, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// try loading into a list of maps first
	m := make([]map[interface{}]interface{}, 5)

	err = yaml.Unmarshal(b, &m)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return m, nil
}

func loadJson(path string) ([]map[interface{}]interface{}, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// try loading into a list of maps first
	m := make(map[string]interface{}, 5)

	err = json.Unmarshal(b, &m)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Info(m)

	out := make([]map[interface{}]interface{}, 1)
	return out, nil
}

func flattenMap(env map[string]string, ev map[interface{}]interface{}, prefix []string) map[string]string {
	for k, v := range ev {
		switch t := v.(type) {
		case string:
			key := flatKey(prefix, k.(string))
			env[key] = v.(string)
		case map[interface{}]interface{}:
			flattenMap(env, v.(map[interface{}]interface{}), append(prefix, k.(string)))
		default:
			log.Debug("Prefix: ", prefix)
			log.Debugf("Default: %#v", t)
		}
	}

	return env
}

func flattenEnv(env []map[interface{}]interface{}) map[string]string {
	fenv := make(map[string]string)
	for _, ev := range env {
		// must be a map
		flattenMap(fenv, ev, []string{})
	}
	return fenv
}

type envParser func(string) ([]map[interface{}]interface{}, error)

func findParser(path string) envParser {
	json := []string{"json"}
	yaml := []string{"yaml", "yml"}

	// Use the file extension first
	for _, sfx := range json {
		if strings.HasSuffix(path, sfx) {
			return loadJson
		}
	}

	for _, sfx := range yaml {
		if strings.HasSuffix(path, sfx) {
			return loadYaml
		}
	}

	// peek at the first bytes for a {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(f)
	start, err := reader.Peek(3)
	if err != nil {
		log.Fatal(err)
	}

	f.Close()

	// if the { exists, we'll try json
	if bytes.HasPrefix(bytes.TrimLeft(start, " "), []byte("{")) {
		return loadJson
	}

	return loadYaml
}

type Action interface {
	Apply() map[string]string
}

type EnvFile struct {
	path string
}

func (e EnvFile) Parse() (map[string]string, error) {
	parser := findParser(e.path)

	data, err := parser(e.path)
	if err != nil {
		log.Fatal(err)
	}

	env := flattenEnv(data)
	log.Debug("Flattened Env")
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
	// TODO: Figure out a good way to delete this file
	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}

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
		switch {
		case f == "--env" || f == "-e":
			in_flag = "env"
		case in_flag == "env":
			log.Debug("Applying  env: ", f)
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

		default:
			break
		}
	}

	return nil
}
