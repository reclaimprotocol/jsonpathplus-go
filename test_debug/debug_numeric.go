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

	fmt.Println("=== Debugging Numeric Comparison ===")
	fmt.Printf("JSONPath: $.users[?(@.age > 30)]\n")
	fmt.Println("Ages in data:")
	fmt.Println("- Alice: 30 (not > 30)")
	fmt.Println("- Bob: 25 (not > 30)") 
	fmt.Println("- Charlie: 35 (> 30)")
	fmt.Println("- Diana: 28 (not > 30)")
	fmt.Println("Expected: 1 result (only Charlie)")
	
	results, err := jp.Query("$.users[?(@.age > 30)]", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("Actual results: %d\n", len(results))
	for i, result := range results {
		fmt.Printf("  [%d] %v\n", i, result.Value)
	}
}