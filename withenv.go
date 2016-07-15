package bach

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/flynn/go-shlex"
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

func flatKey(prefix []string, key string) string {
	return strings.Join(append(prefix, key), "_")
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

func yamlToEnv(path string) (map[string]string, error) {
	env, err := loadYaml(path)
	if err != nil {
		return nil, err
	}

	e := flattenEnv(env)
	log.Debug("Flattened Env")
	for k, v := range e {
		log.Debugf("%s = %s", k, v)
	}

	return e, nil
}

type Action interface {
	Apply() map[string]string
}

type EnvFile struct {
	path string
}

func (e EnvFile) Apply() map[string]string {
	env, err := yamlToEnv(e.path)
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

	parts, err := shlex.Split(e.cmd)
	if err != nil {
		panic(err)
	}

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
			log.Info("Applying  env: ", f)
			action := EnvFile{path: f}
			action.Apply()
			in_flag = ""

		case f == "--script" || f == "-s":
			in_flag = "script"
		case in_flag == "script":
			log.Info("Applying script: ", f)
			action := EnvScript{cmd: f}
			action.Apply()
			in_flag = ""

		default:
			break
		}
	}

	return nil
}
