package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	jsonData := `{
		"store": {
			"book": [
				{
					"category": "reference",
					"author": "Nigel Rees",
					"title": "Sayings of the Century",
					"price": 8.95
				},
				{
					"category": "fiction", 
					"author": "Evelyn Waugh",
					"title": "Sword of Honour",
					"price": 12.99
				},
				{
					"category": "fiction",
					"author": "Herman Melville", 
					"title": "Moby Dick",
					"isbn": "0-553-21311-3",
					"price": 8.99
				},
				{
					"category": "fiction",
					"author": "J. R. R. Tolkien",
					"title": "The Lord of the Rings",
					"isbn": "0-395-19395-8",
					"price": 22.99
				}
			]
		}
	}`

	fmt.Println("=== Testing @parentProperty for book properties ===")
	
	// Test 1: Get all book properties to see the structure
	fmt.Println("\n1. All book properties: $..book.*")
	results1, err := jp.Query("$..book.*", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("Results: %d\n", len(results1))
	for i, result := range results1[:8] { // Show first 8
		fmt.Printf("  Result %d: %s -> %v\n", i, result.Path, result.Value)
	}
	if len(results1) > 8 {
		fmt.Printf("  ... and %d more\n", len(results1)-8)
	}
	
	// Test 2: Filter by @parentProperty
	fmt.Println("\n2. Book properties where parent property != 0: $..book.*[?(@parentProperty !== 0)]")
	results2, err := jp.Query("$..book.*[?(@parentProperty !== 0)]", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("Results: %d (expected: 12)\n", len(results2))
	for i, result := range results2 {
		fmt.Printf("  Result %d: %s -> %v\n", i, result.Path, result.Value)
	}
	
	// Test 3: Test book index 1 properties specifically
	fmt.Println("\n3. Book index 1 properties: $.store.book[1].*")
	results3, err := jp.Query("$.store.book[1].*", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("Results: %d\n", len(results3))
	for i, result := range results3 {
		fmt.Printf("  Result %d: %s -> %v\n", i, result.Path, result.Value)
	}
	
	// Test 4: Test specific property with filter
	fmt.Println("\n4. Book 1 author with parent property filter: $.store.book[1].author[?(@parentProperty !== 0)]")
	results4, err := jp.Query("$.store.book[1].author[?(@parentProperty !== 0)]", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("Results: %d\n", len(results4))
	for i, result := range results4 {
		fmt.Printf("  Result %d: %s -> %v\n", i, result.Path, result.Value)
	}
}