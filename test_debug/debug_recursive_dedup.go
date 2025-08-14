package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	jsonData := `{
		"store": {
			"book": [
				{"title": "Book1", "price": 8.95}
			],
			"bicycle": {
				"color": "red",
				"price": 19.95
			}
		}
	}`

	fmt.Println("=== Debugging Recursive Descent Deduplication ===")
	
	// Test 1: Simple recursive descent to see duplicates
	fmt.Println("\n1. Recursive descent without filter: $..*")
	results1, err := jp.Query("$..*", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("Total results: %d\n", len(results1))
	
	pathCounts := make(map[string]int)
	for _, result := range results1 {
		pathCounts[result.Path]++
	}
	
	fmt.Println("Path frequencies:")
	for path, count := range pathCounts {
		if count > 1 {
			fmt.Printf("  %s: %d times ⚠️ DUPLICATE\n", path, count)
		} else {
			fmt.Printf("  %s: %d time\n", path, count)
		}
	}
	
	// Test 2: Recursive descent with filter
	fmt.Println("\n2. Recursive descent with price filter: $..*[?(@property === 'price')]")
	results2, err := jp.Query("$..*[?(@property === 'price')]", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	fmt.Printf("Total results: %d\n", len(results2))
	
	pathCounts2 := make(map[string]int)
	for _, result := range results2 {
		pathCounts2[result.Path]++
	}
	
	fmt.Println("Filtered path frequencies:")
	for path, count := range pathCounts2 {
		if count > 1 {
			fmt.Printf("  %s: %d times ⚠️ DUPLICATE\n", path, count)
		} else {
			fmt.Printf("  %s: %d time\n", path, count)
		}
	}
	
	// Expected: 2 price values
	fmt.Println("\nExpected paths: $.store.book[0].price, $.store.bicycle.price")
}