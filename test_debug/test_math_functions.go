package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	jsonData := `{
		"users": [
			{"name": "Alice Johnson", "email": "alice@example.com", "age": 30.7},
			{"name": "Bob Smith", "email": "bob@test.org", "age": 25.3},
			{"name": "Charlie Brown", "email": "charlie@example.com", "age": 35.9},
			{"name": "Diana Prince", "email": "diana@hero.gov", "age": 28.1}
		]
	}`

	fmt.Println("=== Testing Math Functions ===")
	
	// Test the math functions that are skipped
	tests := []struct {
		name     string
		jsonpath string
		expected int
		description string
	}{
		{"Math floor", "$.users[?(@.age.floor() > 25)]", 3, "Users with floor(age) > 25"},
		{"Math round", "$.users[?(@.age.round() === 30)]", 2, "Users with rounded age of 30"},  // 30.7 rounds to 31, 28.1 rounds to 28
		{"Math ceil", "$.users[?(@.age.ceil() > 30)]", 2, "Users with ceil(age) > 30"},
	}

	for _, test := range tests {
		fmt.Printf("\n=== %s ===\n", test.name)
		fmt.Printf("JSONPath: %s\n", test.jsonpath)
		fmt.Printf("Description: %s\n", test.description)
		
		results, err := jp.Query(test.jsonpath, jsonData)
		if err != nil {
			fmt.Printf("❌ Query Error: %v\n", err)
			continue
		}

		fmt.Printf("Results: %d (expected: %d)\n", len(results), test.expected)
		if len(results) == test.expected {
			fmt.Printf("✅ Test would pass\n")
		} else {
			fmt.Printf("❌ Test would fail\n")
		}
		
		for i, result := range results {
			fmt.Printf("  [%d] %v\n", i, result.Value)
		}
	}
}