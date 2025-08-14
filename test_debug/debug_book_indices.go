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

	fmt.Println("=== Testing Book Array Access ===")
	
	// Test 1: Direct array access
	fmt.Println("\n1. Direct book access: $.store.book")
	results1, err := jp.Query("$.store.book", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("Results: %d\n", len(results1))
	for i, result := range results1 {
		fmt.Printf("  Result %d: %s (type: %T)\n", i, result.Path, result.Value)
	}
	
	// Test 2: Individual book elements
	fmt.Println("\n2. Individual book elements: $.store.book[*]")
	results2, err := jp.Query("$.store.book[*]", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("Results: %d\n", len(results2))
	for i, result := range results2 {
		fmt.Printf("  Result %d: %s\n", i, result.Path)
	}
	
	// Test 3: Books with index filter
	fmt.Println("\n3. Books with index filter: $.store.book[?(@property !== 0)]")
	results3, err := jp.Query("$.store.book[?(@property !== 0)]", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("Results: %d (expected: 3)\n", len(results3))
	for i, result := range results3 {
		fmt.Printf("  Result %d: %s\n", i, result.Path)
	}
	
	// Test 4: Recursive descent to book
	fmt.Println("\n4. Recursive descent to book: $..book")
	results4, err := jp.Query("$..book", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("Results: %d\n", len(results4))
	for i, result := range results4 {
		fmt.Printf("  Result %d: %s (type: %T)\n", i, result.Path, result.Value)
	}
	
	// Test 5: The failing case
	fmt.Println("\n5. The failing case: $..book[?(@property !== 0)]")
	results5, err := jp.Query("$..book[?(@property !== 0)]", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("Results: %d (expected: 3)\n", len(results5))
	for i, result := range results5 {
		fmt.Printf("  Result %d: %s\n", i, result.Path)
	}
}