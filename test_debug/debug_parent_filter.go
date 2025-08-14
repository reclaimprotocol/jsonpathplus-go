package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	jsonData := `{
		"store": {
			"book": [
				{"title": "Book1", "price": 8.95},
				{"title": "Book2", "price": 12.99}
			],
			"bicycle": {
				"color": "red",
				"price": 19.95
			}
		}
	}`

	fmt.Println("=== Testing @parent filters ===")
	
	// Test 1: Simple parent existence
	fmt.Println("\n1. Books where parent has bicycle: $.store.book[?(@parent.bicycle)]")
	results1, err := jp.Query("$.store.book[?(@parent.bicycle)]", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("Results: %d (expected: 2)\n", len(results1))
	for i, result := range results1 {
		fmt.Printf("  Result %d: %s -> %v\n", i, result.Path, result.Value)
	}
	
	// Test 2: Parent property comparison
	fmt.Println("\n2. Books where parent bicycle is red: $.store.book[?(@parent.bicycle.color === 'red')]")
	results2, err := jp.Query("$.store.book[?(@parent.bicycle.color === 'red')]", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("Results: %d (expected: 2)\n", len(results2))
	for i, result := range results2 {
		fmt.Printf("  Result %d: %s -> %v\n", i, result.Path, result.Value)
	}
	
	// Test 3: Check what books we get without filter
	fmt.Println("\n3. All books: $.store.book")
	results3, err := jp.Query("$.store.book", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("Results: %d\n", len(results3))
	for i, result := range results3 {
		fmt.Printf("  Result %d: %s -> %v\n", i, result.Path, result.Value)
	}
	
	// Test 4: Individual book elements
	fmt.Println("\n4. Individual book elements: $.store.book[*]")
	results4, err := jp.Query("$.store.book[*]", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("Results: %d\n", len(results4))
	for i, result := range results4 {
		fmt.Printf("  Result %d: %s -> %v\n", i, result.Path, result.Value)
	}
}