package main

import (
	"encoding/json"
	"fmt"
)

// printJSON marshalls an interface and prints to STDOUT
func printJSON(i interface{}) error {
	result, err := json.MarshalIndent(&i, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(result))
	return nil
}
