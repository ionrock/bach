package bach

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
)

type EnvVar struct {
	Name  string
	Value string
}

// func CompileEnv(newEnv map[string]interface{}) {
// 	for k, v := range newEnv {
// 		value, ok := v.(string)
// 		if !ok {
// 			os.Setenv(k, UpdateEnv(v))
// 		}
// 	}
// }

func LoadYaml(path string) (map[string]interface{}, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var container map[string]interface{}

	err = yaml.Unmarshal(b, &container)
	if err != nil {
		return nil, err
	}

	return container, nil
}

func WithEnv(envfile string) error {
	fmt.Println(envfile)
	b, err := ioutil.ReadFile(envfile)
	if err != nil {
		return err
	}

	var container map[string]interface{}

	err = yaml.Unmarshal(b, &container)
	if err != nil {
		return err
	}

	return nil
}
