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
				{"title": "Book 2", "price": "12.99"},
				{"title": "Book 3", "price": 8.99},
				{"title": "Book 4", "price": true}
			]
		}
	}`
	
	// Test filter books where price is a number
	fmt.Println("=== Testing typeof function with proper usage ===")
	results, err := jp.Query("$.store.book[?(@.price.typeof() === 'number')]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Books with numeric price: %d results\n", len(results))
	
	for i, r := range results {
		book := r.Value.(map[string]interface{})
		fmt.Printf("  [%d] %s (price: %v)\n", i, book["title"], book["price"])
	}
}