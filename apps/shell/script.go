package main

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type ExperimentScript struct {
	Name string `yaml:"name"`
	Info struct {
		Version     string `yaml:"version"`
		Question    string `yaml:"question"`
		Description string `yaml:"description"`
	} `yaml:"info"`
	Env    []string          `yaml:"env"`
	Params map[string]string `yaml:"params"`
	Setup  struct {
		SQL      []string `yaml:"sql"`
		Diligent []string `yaml:"diligent"`
	} `yaml:"setup"`
	Experiment []string `yaml:"experiment"`
	Conclusion []string `yaml:"conclusion"`
}

type Replacements struct {
	Name   string
	Env    map[string]string
	Params map[string]string
}

func LoadScript(scriptFileName string) (*ExperimentScript, error) {
	var err error

	file, err := os.Open(scriptFileName)
	if err != nil {
		return nil, err
	}
	yamlBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	script := &ExperimentScript{}
	err = yaml.Unmarshal(yamlBytes, &script)
	if err != nil {
		return nil, err
	}

	return script, nil
}
