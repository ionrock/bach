package bach

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func LoadYaml(path string) ([]map[interface{}]interface{}, error) {
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

func FlattenMap(env map[string]string, ev map[interface{}]interface{}, prefix []string) map[string]string {
	for k, v := range ev {
		switch t := v.(type) {
		case string:
			key := flatKey(prefix, k.(string))
			env[key] = v.(string)
		case map[interface{}]interface{}:
			FlattenMap(env, v.(map[interface{}]interface{}), append(prefix, k.(string)))
		default:
			log.Debug("Prefix: ", prefix)
			log.Debugf("Default: %#v", t)
		}
	}

	return env
}

func FlattenEnv(env []map[interface{}]interface{}) map[string]string {
	fenv := make(map[string]string)
	for _, ev := range env {
		// must be a map
		FlattenMap(fenv, ev, []string{})
	}
	return fenv
}
func YamlToEnv(path string) (map[string]string, error) {
	env, err := LoadYaml(path)
	if err != nil {
		return nil, err
	}

	e := FlattenEnv(env)
	log.Debug("Flattened Env")
	for k, v := range e {
		log.Debugf("%s = %s", k, v)
	}

	return e, nil
}

func WithEnv(envfile string) error {
	env, err := YamlToEnv(envfile)
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
	return nil
}
