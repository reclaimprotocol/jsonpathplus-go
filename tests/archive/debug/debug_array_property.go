package jsonpathplus

import (
	"fmt"
	
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func debug_array_propertyMain() {
	// Test data - array with mixed types
	data := `{"data":[42,"hello",true,null,{"key":"value"},[1,2,3]]}`
	
	queries := []string{
		"$.data",
		"$.data[*]",
		"$.data[?(@property === 0)]",
		"$.data[?(@property === 1)]",
		"$.data[?(@property !== 0)]",
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
			fmt.Printf("[%d] Value: %v, Path: %s, ParentProperty: '%s'\n", 
				i, result.Value, result.Path, result.ParentProperty)
		}
	}
}