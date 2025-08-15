package jsonpathplus

import (
	"fmt"
	
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func debug_mainMain() {
	// Test data
	data := `{"users":{"1":{"name":"Alice","age":25,"active":true},"2":{"name":"Bob","age":30,"active":false}}}`
	
	// Test queries step by step
	queries := []string{
		"$.users",
		"$.users.1",
		"$.users.1.*",
		"$.users.1[?(@parentProperty === '1')]",
		"$.users['1']",
		"$.users['1'][?(@parentProperty === '1')]",
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
			fmt.Printf("[%d] Value: %v, Path: %s, Parent: %v, ParentProperty: %s\n", 
				i, result.Value, result.Path, result.Parent, result.ParentProperty)
		}
	}
}