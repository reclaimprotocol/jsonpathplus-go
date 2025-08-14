package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	jsonData := `{
		"store": {
			"book": [
				{"title": "Book 1", "price": 8.95},
				{"title": "Book 2", "price": 12.99},
				{"title": "Book 3", "price": 8.99},
				{"title": "Book 4", "price": 22.99}
			]
		}
	}`
	
	// Test basic access first
	fmt.Println("=== Testing basic price access ===")
	results, err := jp.Query("$.store.book[*].price", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	for i, r := range results {
		fmt.Printf("Price %d: %v (type: %T)\n", i, r.Value, r.Value)
	}
	
	// Test typeof function
	fmt.Println("\n=== Testing typeof function ===")
	results2, err := jp.Query("$.store.book[*].price[?(@.typeof() === 'number')]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Numeric prices: %d results\n", len(results2))
	
	for i, r := range results2 {
		fmt.Printf("  [%d] %v\n", i, r.Value)
	}
}