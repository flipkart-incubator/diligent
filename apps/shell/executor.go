package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"text/template"
	"time"
)

const (
	sqlTimeout = 5 * time.Second
)

type Executor struct {
	script       *ExperimentScript
	values       *ExperimentValues
	replacements *Replacements
	db           *sql.DB
	dryRun       bool
}

func NewExecutor(script *ExperimentScript, values *ExperimentValues, dryRun bool) (*Executor, error) {
	env, err := getEnvValues(script, values)
	if err != nil {
		return nil, err
	}
	params, err := getParamValues(script, values)
	if err != nil {
		return nil, err
	}
	replacements := &Replacements{
		Name:   script.Name,
		Env:    env,
		Params: params,
	}

	dbDriver := replacements.Env["dbDriver"]
	switch dbDriver {
	case "mysql", "pgx":
	case "":
		return nil, fmt.Errorf("required env not found: %s", "dbDriver")
	default:
		return nil, fmt.Errorf("invalid env value dbDriver='%s'. Allowed values are 'mysql', 'pgx'", dbDriver)
	}

	// Validate URL
	dbUrl := replacements.Env["dbUrl"]
	if dbUrl == "" {
		return nil, fmt.Errorf("required env not found: %s", "dbUrl")
	}

	// Open new connection
	db, err := sql.Open(dbDriver, dbUrl)
	if err != nil {
		return nil, err
	}

	return &Executor{
		script:       script,
		values:       values,
		replacements: replacements,
		db:           db,
		dryRun:       dryRun,
	}, nil
}

func (e *Executor) Execute() error {
	fmt.Printf("Running benchmark: %s\n", e.script.Name)
	fmt.Printf("Dry run: %v\n", e.dryRun)
	fmt.Printf("General information:\n")
	fmt.Printf("\tVersion: %s\n", e.script.Info.Version)
	fmt.Printf("\tQuestion: %s\n", e.script.Info.Question)
	fmt.Printf("\tDescription: %s\n", e.script.Info.Description)
	fmt.Printf("\nEnvironment:\n")
	for k, v := range e.replacements.Env {
		fmt.Printf("\t%s: %s\n", k, v)
	}
	fmt.Printf("\nParams:\n")
	for k, v := range e.replacements.Params {
		fmt.Printf("\t%s: %s\n", k, v)
	}

	fmt.Printf("\nStarting Setup Phase:\n")
	fmt.Printf("Executing SQL commands:\n")
	for _, cmd := range e.script.Setup.SQL {
		err := e.executeSQLCmd(cmd)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Executing Diligent commands:\n")
	for _, cmd := range e.script.Setup.Diligent {
		err := e.executeDiligentCmd(cmd)
		if err != nil {
			return err
		}
	}
	fmt.Printf("Setup Phase Completed.\n")

	fmt.Printf("\nStarting Execution Phase:\n")
	fmt.Printf("Executing Diligent commands:\n")
	for _, cmd := range e.script.Experiment {
		err := e.executeDiligentCmd(cmd)
		if err != nil {
			return err
		}
	}
	fmt.Printf("Execution Phase Completed.\n")

	fmt.Printf("\nStarting Conclusion Phase:\n")
	fmt.Printf("Executing Diligent commands:\n")
	for _, cmd := range e.script.Conclusion {
		err := e.executeDiligentCmd(cmd)
		if err != nil {
			return err
		}
	}
	fmt.Printf("Conclusion Phase Completed.\n")
	return nil
}

func getEnvValues(script *ExperimentScript, values *ExperimentValues) (map[string]string, error) {
	env := make(map[string]string)
	notFound := make([]string, 0)

	for _, key := range script.Env {
		value, ok := values.Env[key]
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

func getParamValues(script *ExperimentScript, values *ExperimentValues) (map[string]string, error) {
	params := make(map[string]string)
	unknowns := make([]string, 0)

	// Make a copy of input parameters
	for key, val := range script.Params {
		params[key] = val
	}

	// Go over the provided overrides
	for key, val := range values.Overrides {
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

func (e *Executor) executeSQLCmd(cmdTmplStr string) error {
	tmpl, err := template.New("t1").Parse(cmdTmplStr)
	if err != nil {
		return err
	}
	var sb strings.Builder
	err = tmpl.Execute(&sb, e.replacements)
	if err != nil {
		return err
	}
	cmd := sb.String()
	fmt.Println(">>", "sql", cmd)

	if !e.dryRun {
		sqlCtx, sqlCancel := context.WithTimeout(context.Background(), sqlTimeout)
		_, err = e.db.ExecContext(sqlCtx, cmd)
		sqlCancel()
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Executor) executeDiligentCmd(cmdTmplStr string) error {
	tmpl, err := template.New("t1").Parse(cmdTmplStr)
	if err != nil {
		return err
	}
	var sb strings.Builder
	err = tmpl.Execute(&sb, e.replacements)
	if err != nil {
		return err
	}
	cmdStr := sb.String()
	fmt.Println(">>", cmdStr)

	if !e.dryRun {
		tokens := strings.Split(cmdStr, " ")
		return grumbleApp.RunCommand(tokens)
	}

	return nil
}
