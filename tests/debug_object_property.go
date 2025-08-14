package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	data := `{"info":{"0":"string_zero","1":"string_one"}}`
	
	fmt.Println("=== Testing Simple Object Property Filtering ===")
	
	queries := []string{
		"$.info[*]",                    // All info properties to see context
		"$.info[?(@property === '0')]", // Should return info['0']
		"$.info[?(@property === '1')]", // Should return info['1']
		"$.info[?(@property !== '0')]", // Should return info['1']
	}
	
	for _, query := range queries {
		fmt.Printf("\nQuery: %s\n", query)
		results, err := jp.Query(query, data)
		if err != nil {
			fmt.Printf("Go Error: %s\n", err.Error())
		} else {
			fmt.Printf("Go Count: %d\n", len(results))
			for i, r := range results {
				fmt.Printf("  [%d] Path: %s, ParentProperty: '%s'\n", i, r.Path, r.ParentProperty)
			}
		}
	}
}