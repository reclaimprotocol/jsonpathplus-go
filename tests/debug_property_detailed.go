package main

import (
	"encoding/json"
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	data := `{"store":{"book":[{"title":"Book0"},{"title":"Book1"}],"info":{"0":"string_zero","1":"string_one"}}}`
	
	fmt.Println("=== Testing Property Value Types ===")
	
	// Test different property filter queries to see what types we get
	queries := []string{
		"$..book[*]",        // Array elements - should get numbers
		"$..info[*]",        // Object keys - should get strings
	}
	
	for _, query := range queries {
		fmt.Printf("\nQuery: %s\n", query)
		results, err := jp.Query(query, data)
		if err != nil {
			fmt.Printf("Go Error: %s\n", err.Error())
		} else {
			fmt.Printf("Go Count: %d\n", len(results))
			for i, r := range results {
				valueStr, _ := json.Marshal(r.Value)
				fmt.Printf("  [%d] Path: %s, Value: %s\n", i, r.Path, string(valueStr))
				fmt.Printf("      ParentProperty: '%s' (string)\n", r.ParentProperty)
				// We can't directly access GetPropertyValue() from here since it's on Context
				// But we can infer what it should be based on the parent type
			}
		}
	}
	
	// Test specific filter cases
	fmt.Printf("\n=== Testing Specific Filter Cases ===\n")
	testCases := []struct {
		query string
		desc  string
	}{
		{"$..book[?(@property === 0)]", "Array: @property === 0"},
		{"$..book[?(@property === '0')]", "Array: @property === '0'"},
		{"$..book[?(@property !== 0)]", "Array: @property !== 0"},
		{"$..info[?(@property === 0)]", "Object: @property === 0"},
		{"$..info[?(@property === '0')]", "Object: @property === '0'"},
		{"$..info[?(@property !== 0)]", "Object: @property !== 0"},
	}
	
	for _, test := range testCases {
		fmt.Printf("\n%s\n", test.desc)
		fmt.Printf("Query: %s\n", test.query)
		results, err := jp.Query(test.query, data)
		if err != nil {
			fmt.Printf("Go Error: %s\n", err.Error())
		} else {
			fmt.Printf("Go Count: %d\n", len(results))
			if len(results) > 0 {
				fmt.Printf("Go Paths: [")
				for i, r := range results {
					if i > 0 { fmt.Print(", ") }
					fmt.Printf("\"%s\"", r.Path)
				}
				fmt.Println("]")
			}
		}
	}
}