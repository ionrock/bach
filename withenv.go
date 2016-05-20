package bach

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Env struct {
	Vars []EnvVar
}

type EnvVar struct {
	Key   string
	Value string
}

func NewEnv(path string) (*Env, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var container Env

	err = yaml.Unmarshal(b, &container)
	if err != nil {
		return nil, err
	}

	return &container, nil
}

func WithEnv(envfile string) error {
	env, err := NewEnv(envfile)
	if err != nil {
		panic(err)
	}

	for _, ev := range env.Vars {
		err = os.Setenv(ev.Key, ev.Value)
		if err != nil {
			panic(err)
		}
	}
	return nil
}
