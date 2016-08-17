package bach

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
)

func NewFlattener(path string) Flattener {
	json := []string{"json"}
	yaml := []string{"yaml", "yml"}

	// Use the file extension first
	for _, sfx := range json {
		if strings.HasSuffix(path, sfx) {
			log.Debug("Using JSON parser")
			return JsonFlattener{path}
		}
	}

	for _, sfx := range yaml {
		if strings.HasSuffix(path, sfx) {
			log.Debug("Using YAML parser")
			return YamlFlattener{path}
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
		log.Debug("Using JSON parser")
		return JsonFlattener{path}
	}

	// if the [ exists, we'll try json
	if bytes.HasPrefix(bytes.TrimLeft(start, " "), []byte("[")) {
		log.Debug("Using JSON parser")
		return JsonFlattener{path}
	}

	log.Debug("Using YAML parser")
	return YamlFlattener{path}
}

type Flattener interface {
	Flatten() (map[string]string, error)
}

type JsonFlattener struct {
	path string
}

func (f JsonFlattener) loadMap(b []byte) ([]map[string]interface{}, error) {
	// try loading into a list of maps first
	m := make(map[string]interface{})

	err := json.Unmarshal(b, &m)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return []map[string]interface{}{m}, nil
}

func (f JsonFlattener) loadList(b []byte) ([]map[string]interface{}, error) {
	// try loading into a list of maps first
	m := make([]map[string]interface{}, 5)

	err := json.Unmarshal(b, &m)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return m, nil
}

func (f JsonFlattener) load(path string) ([]map[string]interface{}, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	m, err := f.loadMap(b)
	if err == nil {
		return m, err
	}

	log.Errorf("Error parsing %s as an object", path)

	m, err = f.loadList(b)
	if err != nil {
		log.Fatal(err)
	}
	return m, nil
}

func flatKey(prefix []string, key string) string {
	return strings.Join(append(prefix, key), "_")
}

func (f JsonFlattener) flattenEnv(env []map[string]interface{}) map[string]string {
	fenv := make(map[string]string)
	for _, ev := range env {
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

func (f JsonFlattener) Flatten() (map[string]string, error) {
	env, err := f.load(f.path)
	if err != nil {
		log.Error("Error loading JSON")
		return nil, err
	}

	return f.flattenEnv(env), nil
}

type YamlFlattener struct {
	path string
}

func (f YamlFlattener) load(path string) ([]map[interface{}]interface{}, error) {
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

func (f YamlFlattener) flattenMap(env map[string]string, ev map[interface{}]interface{}, prefix []string) map[string]string {
	for k, v := range ev {
		switch t := v.(type) {
		case string:
			key := flatKey(prefix, k.(string))
			env[key] = v.(string)
		case map[interface{}]interface{}:
			f.flattenMap(env, v.(map[interface{}]interface{}), append(prefix, k.(string)))
		default:
			log.Debug("Prefix: ", prefix)
			log.Debugf("Default: %#v", t)
		}
	}

	return env
}

func (f YamlFlattener) flattenEnv(env []map[interface{}]interface{}) map[string]string {
	fenv := make(map[string]string)
	for _, ev := range env {
		// must be a map
		f.flattenMap(fenv, ev, []string{})
	}
	return fenv
}

func (f YamlFlattener) Flatten() (map[string]string, error) {
	env, err := f.load(f.path)
	if err != nil {
		log.Error("Error loading YAML")
		return nil, err
	}

	return f.flattenEnv(env), nil
}
