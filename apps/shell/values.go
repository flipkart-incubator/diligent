package main

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type ExperimentValues struct {
	Env       map[string]string `yaml:"env"`
	Overrides map[string]string `yaml:"overrides"`
}

func LoadValues(valuesFileName string) (*ExperimentValues, error) {
	var err error

	file, err := os.Open(valuesFileName)
	if err != nil {
		return nil, err
	}
	yamlBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	values := &ExperimentValues{}
	err = yaml.Unmarshal(yamlBytes, &values)
	if err != nil {
		return nil, err
	}

	return values, nil
}
