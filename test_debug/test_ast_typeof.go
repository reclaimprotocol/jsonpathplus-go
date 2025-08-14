package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	// Test what $.store.book[*].price returns
	jsonData := `{
		"store": {
			"book": [
				{"title": "Book 1", "price": 8.95},
				{"title": "Book 2", "price": 12.99}
			]
		}
	}`

	fmt.Println("=== Testing $.store.book[*].price ===")
	results, err := jp.Query("$.store.book[*].price", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for i, r := range results {
		fmt.Printf("Price %d: %v (type: %T) (path: %s)\n", i, r.Value, r.Value, r.Path)
	}

	// Test direct filter on these prices
	fmt.Println("\n=== Testing direct filter on prices ===")
	results2, err := jp.Query("$.store.book[*].price[?(@ === 8.95)]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Direct filter results: %d\n", len(results2))
}
