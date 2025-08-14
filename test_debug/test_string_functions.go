package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	jsonData := `{
		"users": [
			{"name": "Alice Johnson", "email": "alice@example.com", "age": 30},
			{"name": "Bob Smith", "email": "bob@test.org", "age": 25},
			{"name": "Charlie Brown", "email": "charlie@example.com", "age": 35},
			{"name": "Diana Prince", "email": "diana@hero.gov", "age": 28}
		]
	}`

	fmt.Println("=== Testing String Functions ===")

	// Test each function individually
	tests := []struct {
		name     string
		jsonpath string
		expected int
	}{
		{"Contains function", "$.users[?(@.email.contains('example'))]", 2},
		{"StartsWith function", "$.users[?(@.name.startsWith('A'))]", 1},
		{"EndsWith function", "$.users[?(@.email.endsWith('.com'))]", 2},
		{"Match function", "$.users[?(@.email.match(/.*@example\\.com$/))]", 2},
	}

	for _, test := range tests {
		fmt.Printf("\n=== %s ===\n", test.name)
		fmt.Printf("JSONPath: %s\n", test.jsonpath)

		results, err := jp.Query(test.jsonpath, jsonData)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			continue
		}

		fmt.Printf("Results: %d (expected: %d)\n", len(results), test.expected)
		if len(results) == test.expected {
			fmt.Printf("✅ Test passed\n")
		} else {
			fmt.Printf("❌ Test failed\n")
		}

		for i, result := range results {
			fmt.Printf("  [%d] %v\n", i, result.Value)
		}
	}
}
