package main

import (
	"fmt"
	"strings"
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

	fmt.Println("=== Testing @property filter step by step ===")
	
	// Step 1: Test recursive descent without filter
	fmt.Println("\n1. Recursive descent without filter: $..*")
	results1, err := jp.Query("$..*", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("Total results: %d\n", len(results1))
	
	// Show only price-related paths
	priceCount := 0
	for _, result := range results1 {
		if strings.Contains(result.Path, "price") {
			priceCount++
			fmt.Printf("  Price path: %s = %v\n", result.Path, result.Value)
		}
	}
	fmt.Printf("Price paths found: %d\n", priceCount)
	
	// Step 2: Test simple property filter on a direct path
	fmt.Println("\n2. Simple property filter on direct path: $.store.book[*][?(@property === 'price')]")
	results2, err := jp.Query("$.store.book[*][?(@property === 'price')]", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("Results: %d\n", len(results2))
	for i, result := range results2 {
		fmt.Printf("  [%d] %v (path: %s)\n", i, result.Value, result.Path)
	}
	
	// Step 3: Test the problematic case
	fmt.Println("\n3. Problematic case: $..*[?(@property === 'price')]")
	results3, err := jp.Query("$..*[?(@property === 'price')]", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("Results: %d (should be 3 price values)\n", len(results3))
	
	// Show first few results to understand what's happening
	for i, result := range results3 {
		if i < 10 {
			fmt.Printf("  [%d] %v (path: %s)\n", i, result.Value, result.Path)
		}
	}
	if len(results3) > 10 {
		fmt.Printf("  ... and %d more\n", len(results3)-10)
	}
	
	// Step 4: Test other property names to confirm filter logic
	fmt.Println("\n4. Testing with different property: $..*[?(@property === 'title')]")
	results4, err := jp.Query("$..*[?(@property === 'title')]", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("Results: %d (should be 2 title values)\n", len(results4))
	for i, result := range results4 {
		fmt.Printf("  [%d] %v (path: %s)\n", i, result.Value, result.Path)
	}
}