package main

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type BenchmarkValues struct {
	Env       map[string]string `yaml:"env"`
	Overrides map[string]string `yaml:"overrides"`
}

func LoadValues(valuesFileName string) (*BenchmarkValues, error) {
	var err error

	file, err := os.Open(valuesFileName)
	if err != nil {
		return nil, err
	}
	yamlBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	bmv := &BenchmarkValues{}
	err = yaml.Unmarshal(yamlBytes, &bmv)
	if err != nil {
		return nil, err
	}

	return bmv, nil
}
