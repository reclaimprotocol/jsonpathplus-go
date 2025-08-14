package main

import (
	"encoding/json"
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	jsonData := `{
		"store": {
			"book": [
				{"price": 8.95, "title": "Book 1"},
				{"price": 12.99, "title": "Book 2"},
				{"price": 8.99, "title": "Book 3"}
			],
			"bicycle": {
				"price": 19.95,
				"color": "red"
			}
		}
	}`

	fmt.Println("=== Testing recursive descent for prices ===")
	
	// Test 1: Get all prices using recursive descent
	fmt.Println("\n1. Using $..price")
	results1, err := jp.Query("$..price", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("Results: %d\n", len(results1))
		for i, result := range results1 {
			fmt.Printf("  [%d] %v (path: %s)\n", i, result.Value, result.Path)
		}
	}
	
	// Test 2: Get all elements with recursive descent
	fmt.Println("\n2. Using $..*")
	results2, err := jp.Query("$..*", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("Results: %d total elements\n", len(results2))
		// Count how many are actually price properties
		priceCount := 0
		for _, result := range results2 {
			// Check if this is a price value
			if num, ok := result.Value.(float64); ok {
				// Check if it matches one of our price values
				if num == 8.95 || num == 12.99 || num == 8.99 || num == 19.95 {
					priceCount++
					fmt.Printf("  Price found: %v (path: %s)\n", num, result.Path)
				}
			}
		}
		fmt.Printf("Total price values found: %d\n", priceCount)
	}
	
	// Test 3: Try the @property filter
	fmt.Println("\n3. Using $..*[?(@property === 'price')]")
	results3, err := jp.Query("$..*[?(@property === 'price')]", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("Results: %d\n", len(results3))
		// Show only actual price values
		for i, result := range results3 {
			if i < 10 { // Limit output
				marshaled, _ := json.Marshal(result.Value)
				fmt.Printf("  [%d] %s (path: %s)\n", i, string(marshaled), result.Path)
			}
		}
	}
	
	// Test 4: Alternative approach - get all price properties directly
	fmt.Println("\n4. What we SHOULD get with @property === 'price':")
	fmt.Println("  Expected: only the 4 price values (8.95, 12.99, 8.99, 19.95)")
}