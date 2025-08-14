package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	// Use exact same data and expressions as in the failing tests
	jsonData := `{
		"users": [
			{"name": "Alice Johnson", "email": "alice@example.com", "age": 30},
			{"name": "Bob Smith", "email": "bob@test.org", "age": 25},
			{"name": "Charlie Brown", "email": "charlie@example.com", "age": 35},
			{"name": "Diana Prince", "email": "diana@hero.gov", "age": 28}
		],
		"products": [
			{"title": "JavaScript Guide", "category": "programming"},
			{"title": "Python Cookbook", "category": "programming"},
			{"title": "Design Patterns", "category": "architecture"},
			{"title": "Clean Code", "category": "programming"}
		],
		"tags": ["frontend", "backend", "mobile", "web", "api"],
		"descriptions": [
			"This is a great product",
			"Excellent quality item",
			"Good value for money",
			"Outstanding performance"
		]
	}`

	fmt.Println("=== Testing Exact Test Cases from Compatibility Tests ===")

	// These are the exact expressions from the skipped tests
	tests := []struct {
		name        string
		jsonpath    string
		expected    int
		description string
	}{
		{"String contains", "$.users[?(@.email.contains('example'))]", 2, "Find users with 'example' in email"},
		{"String startsWith", "$.users[?(@.name.startsWith('A'))]", 1, "Find users whose name starts with 'A'"},
		{"String endsWith", "$.users[?(@.email.endsWith('.com'))]", 2, "Find users with .com email addresses"},
		{"Regex match", "$.users[?(@.email.match(/.*@example\\.com$/))]", 2, "Find users with example.com email using regex"},
		{"Array length", "$.tags[?(@.length === 5)]", 0, "Check if tags array has 5 elements (applied to individual elements)"},
		{"String includes", "$.products[?(@.category.includes('prog'))]", 0, "Products in categories containing 'prog'"},
		{"typeof string", "$.users[*].name[?(@.typeof() === 'string')]", 4, "Find all string-type names"},
		{"typeof number", "$.users[*].age[?(@.typeof() === 'number')]", 4, "Find all number-type ages"},
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
