package bach

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type Flattener interface {
	Flatten(string) map[string]string
}

type JsonFlattener struct {
	path string
}

func (f JsonFlattener) load(path string) (map[string]interface{}, error) {
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

	return m, nil
}

func flatKey(prefix []string, key string) string {
	return strings.Join(append(prefix, key), "_")
}

func (f JsonFlattener) flattenEnv(env []map[string]interface{}) map[string]string {
	fenv := make(map[string]string)
	for _, ev := range env {
		// must be a map
		f.flattenMap(fenv, ev, []string{})
	}
	return fenv
}

func (f JsonFlattener) flattenMap(env map[string]string, ev map[string]interface{}, prefix []string) map[string]string {
	for k, v := range ev {
		switch t := v.(type) {
		case string:
			key := flatKey(prefix, k)
			env[key] = v.(string)
		case map[string]interface{}:
			f.flattenMap(env, v.(map[string]interface{}), append(prefix, k))
		default:
			log.Debug("Prefix: ", prefix)
			log.Debugf("Default: %#v", t)
		}
	}

	return env
}
