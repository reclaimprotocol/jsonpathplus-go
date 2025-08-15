package jsonpathplus

import (
	"fmt"
	
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func debug_property_newMain() {
	// Test data
	data := `{"users":{"1":{"name":"Alice","age":25,"active":true},"2":{"name":"Bob","age":30,"active":false},"10":{"name":"Charlie","age":35,"active":true}}}`
	
	queries := []string{
		"$.users",
		"$.users.*",
		"$.users[?(@property === '1')]",
		"$.users[?(@property)]",
		"$.users.1[?(@parentProperty === '1')]",
		"$.users.1[?(@parentProperty === 'users')]",
	}
	
	for _, query := range queries {
		fmt.Printf("\n=== Testing: %s ===\n", query)
		results, err := jp.Query(query, data)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			continue
		}
		
		fmt.Printf("Count: %d\n", len(results))
		for i, result := range results {
			fmt.Printf("[%d] Value: %v\n", i, result.Value)
			fmt.Printf("    Path: %s\n", result.Path) 
			fmt.Printf("    Parent: %v\n", result.Parent)
			fmt.Printf("    ParentProperty: '%s'\n", result.ParentProperty)
		}
	}
}