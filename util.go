package main

import (
	"encoding/json"
	"fmt"

	"github.com/ghodss/yaml"
)

// printJSON marshalls an interface and prints to STDOUT
func printResource(i interface{}, yamlOutput bool) error {
	result, err := json.MarshalIndent(&i, "", "  ")
	if err != nil {
		return err
	}

	if yamlOutput {
		result, err := yaml.JSONToYAML(result)
		if err != nil {
			return err
		}
		fmt.Printf("---\n%s", result)
	} else {
		fmt.Printf("%s", result)
	}

	return nil
}
