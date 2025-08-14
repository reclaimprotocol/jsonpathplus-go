package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	jsonStr := `[{"id":1,"name":"first"},{"id":2,"name":"second"}]`

	fmt.Printf("Testing query: $[*].id on JSON: %s\n", jsonStr)

	results, err := jp.Query("$[*].id", jsonStr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Results: %d\n", len(results))
	for i, result := range results {
		fmt.Printf("  [%d] Value: %v, Path: %s\n", i, result.Value, result.Path)
	}

	// Also test the individual parts
	fmt.Printf("\nTesting $[*]:\n")
	results2, err := jp.Query("$[*]", jsonStr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Results: %d\n", len(results2))
		for i, result := range results2 {
			fmt.Printf("  [%d] Value: %v, Path: %s\n", i, result.Value, result.Path)
		}
	}
}
