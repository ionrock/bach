package core

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func WithEnv(envfile string) error {
	fmt.Println(envfile)
	b, err := ioutil.ReadFile(envfile)
	if err != nil {
		return err
	}

	var container map[string]interface{}

	err = yaml.Unmarshal(b, container)
	if err != nil {
		return err
	}
	fmt.Println(container)
	return nil
}
