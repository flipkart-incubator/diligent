package main

import (
	"fmt"
	"strings"
	"text/template"
)

type Executor struct {
	script *BenchmarkScript
	values *BenchmarkValues
}

func NewExecutor(script *BenchmarkScript, values *BenchmarkValues) *Executor {
	return &Executor{
		script: script,
		values: values,
	}
}

func (e *Executor) Execute() error {
	env, err := e.getEnvValues()
	if err != nil {
		return err
	}
	params, err := e.getParamValues()
	if err != nil {
		return err
	}
	replacements := &Replacements{
		Name:   e.script.Name,
		Env:    env,
		Params: params,
	}

	fmt.Printf("Running benchmark: %s\n", e.script.Name)
	fmt.Printf("General information:\n")
	fmt.Printf("\tVersion: %s\n", e.script.Info.Version)
	fmt.Printf("\tQuestion: %s\n", e.script.Info.Question)
	fmt.Printf("\tDescription: %s\n", e.script.Info.Description)
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
	for _, cmd := range e.script.Setup.SQL {
		err := e.executeSQLCmd(cmd, replacements)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Executing Diligent commands:\n")
	for _, cmd := range e.script.Setup.Diligent {
		err := e.executeDiligentCmd(cmd, replacements)
		if err != nil {
			return err
		}
	}
	fmt.Printf("Setup Phase Completed.\n")

	fmt.Printf("\nStarting Execution Phase:\n")
	fmt.Printf("Executing Diligent commands:\n")
	for _, cmd := range e.script.Experiment {
		err := e.executeDiligentCmd(cmd, replacements)
		if err != nil {
			return err
		}
	}
	fmt.Printf("Execution Phase Completed.\n")

	fmt.Printf("\nStarting Conclusion Phase:\n")
	fmt.Printf("Executing Diligent commands:\n")
	for _, cmd := range e.script.Conclusion {
		err := e.executeDiligentCmd(cmd, replacements)
		if err != nil {
			return err
		}
	}
	fmt.Printf("Conclusion Phase Completed.\n")
	return nil
}

func (e *Executor) getEnvValues() (map[string]string, error) {
	env := make(map[string]string)
	notFound := make([]string, 0)

	for _, key := range e.script.Env {
		value, ok := e.values.Env[key]
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

func (e *Executor) getParamValues() (map[string]string, error) {
	params := make(map[string]string)
	unknowns := make([]string, 0)

	// Make a copy of input parameters
	for key, val := range e.script.Params {
		params[key] = val
	}

	// Go over the provided overrides
	for key, val := range e.values.Overrides {
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

func (e *Executor) executeSQLCmd(cmd string, replacements *Replacements) error {
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

func (e *Executor) executeDiligentCmd(cmd string, replacements *Replacements) error {
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
