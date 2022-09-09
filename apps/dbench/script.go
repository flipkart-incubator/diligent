package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

type BenchmarkScript struct {
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

func LoadScript(scriptFileName string) (*BenchmarkScript, error) {
	var err error

	file, err := os.Open(scriptFileName)
	if err != nil {
		return nil, err
	}
	yamlBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	bms := &BenchmarkScript{}
	err = yaml.Unmarshal(yamlBytes, &bms)
	if err != nil {
		return nil, err
	}

	return bms, nil
}

func (s *BenchmarkScript) Execute(bmv *BenchmarkValues) error {
	env, err := s.getEnvValues(bmv)
	if err != nil {
		return err
	}
	params, err := s.getParamValues(bmv)
	if err != nil {
		return err
	}
	replacements := &Replacements{
		Name:   s.Name,
		Env:    env,
		Params: params,
	}

	fmt.Printf("Running benchmark: %s\n", s.Name)
	fmt.Printf("General information:\n")
	fmt.Printf("\tVersion: %s\n", s.Info.Version)
	fmt.Printf("\tQuestion: %s\n", s.Info.Question)
	fmt.Printf("\tDescription: %s\n", s.Info.Description)
	fmt.Printf("\nEnvironment:\n")
	for k, v := range env {
		fmt.Printf("\t%s: %s\n", k, v)
	}
	fmt.Printf("\nParams:\n")
	for k, v := range params {
		fmt.Printf("\t%s: %s\n", k, v)
	}

	fmt.Printf("\nStarting Setup Phase:\n")
	fmt.Printf("Executing SQL commands:\n")
	for _, cmd := range s.Setup.SQL {
		err := s.executeSQLCmd(cmd, replacements)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Executing Diligent commands:\n")
	for _, cmd := range s.Setup.Diligent {
		err := s.executeDiligentCmd(cmd, replacements)
		if err != nil {
			return err
		}
	}
	fmt.Printf("Setup Phase Completed.\n")

	fmt.Printf("\nStarting Execution Phase:\n")
	fmt.Printf("Executing Diligent commands:\n")
	for _, cmd := range s.Experiment {
		err := s.executeDiligentCmd(cmd, replacements)
		if err != nil {
			return err
		}
	}
	fmt.Printf("Execution Phase Completed.\n")

	fmt.Printf("\nStarting Conclusion Phase:\n")
	fmt.Printf("Executing Diligent commands:\n")
	for _, cmd := range s.Conclusion {
		err := s.executeDiligentCmd(cmd, replacements)
		if err != nil {
			return err
		}
	}
	fmt.Printf("Conclusion Phase Completed.\n")
	return nil
}

func (s *BenchmarkScript) getEnvValues(bmv *BenchmarkValues) (map[string]string, error) {
	env := make(map[string]string)
	notFound := make([]string, 0)

	for _, key := range s.Env {
		value, ok := bmv.Env[key]
		if !ok {
			notFound = append(notFound, key)
		} else {
			env[key] = value
		}
	}

	if len(notFound) > 0 {
		return nil, fmt.Errorf("required env parameters not found in values: %s", notFound)
	}

	return env, nil
}

func (s *BenchmarkScript) getParamValues(bmv *BenchmarkValues) (map[string]string, error) {
	params := make(map[string]string)
	unknowns := make([]string, 0)

	// Make a copy of input parameters
	for key, val := range s.Params {
		params[key] = val
	}

	// Go over the provided overrides
	for key, val := range bmv.Overrides {
		// Must be a valid override
		_, ok := params[key]
		if !ok {
			unknowns = append(unknowns, key)
		} else {
			params[key] = val
		}
	}

	if len(unknowns) > 0 {
		return nil, fmt.Errorf("unknown override parameters found in input: %s", unknowns)
	}

	return params, nil
}

func (s *BenchmarkScript) executeSQLCmd(cmd string, replacements *Replacements) error {
	tmpl, err := template.New("t1").Parse(cmd)
	if err != nil {
		return err
	}
	var sb strings.Builder
	err = tmpl.Execute(&sb, replacements)
	if err != nil {
		return err
	}
	fmt.Println("sql>>", sb.String())
	return nil
}

func (s *BenchmarkScript) executeDiligentCmd(cmd string, replacements *Replacements) error {
	tmpl, err := template.New("t1").Parse(cmd)
	if err != nil {
		return err
	}
	var sb strings.Builder
	err = tmpl.Execute(&sb, replacements)
	if err != nil {
		return err
	}
	fmt.Println("dcli>>", sb.String())
	return nil
}
