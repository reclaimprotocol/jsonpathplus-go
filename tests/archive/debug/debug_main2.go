package jsonpathplus

import (
	"fmt"

	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func debug_main2Main() {
	// Test the exact failing case
	data := `{"users":{"1":{"name":"Alice","age":25,"active":true}}}`

	fmt.Println("=== Testing: $.users.1 ===")
	results1, err := jp.Query("$.users.1", data)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	for i, result := range results1 {
		fmt.Printf("[%d] Value: %v\n", i, result.Value)
		fmt.Printf("    Path: %s\n", result.Path)
		fmt.Printf("    Parent: %v\n", result.Parent)
		fmt.Printf("    ParentProperty: '%s'\n", result.ParentProperty)
		fmt.Printf("    Index: %d\n", result.Index)
		fmt.Printf("    OriginalIndex: %d\n", result.OriginalIndex)

		// Check if this is an object we can iterate
		if obj, ok := result.Value.(map[string]interface{}); ok {
			fmt.Printf("    Object with %d properties\n", len(obj))
			for key, value := range obj {
				fmt.Printf("      '%s' = %v\n", key, value)
			}
		}
	}

	fmt.Println("\n=== Testing: $.users.1[?(@parentProperty === '1')] ===")
	results2, err := jp.Query("$.users.1[?(@parentProperty === '1')]", data)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	fmt.Printf("Count: %d\n", len(results2))
	for i, result := range results2 {
		fmt.Printf("[%d] Value: %v, Path: %s, ParentProperty: '%s'\n",
			i, result.Value, result.Path, result.ParentProperty)
	}
}
