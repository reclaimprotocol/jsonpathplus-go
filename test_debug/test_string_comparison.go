package main

import (
	"fmt"

	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	jsonData := `{
		"store": {
			"book": [
				{"category": "fiction", "price": 8.95, "title": "Book 1"},
				{"category": "fiction", "price": 12.99, "title": "Book 2"},
				{"category": "reference", "price": 22.99, "title": "Book 3"}
			]
		}
	}`

	// Test different string comparison operators
	operators := []string{"===", "==", "!=", "!=="}

	for _, op := range operators {
		query := fmt.Sprintf("$.store.book[?(@.category %s 'fiction')]", op)
		fmt.Printf("Testing: %s\n", query)

		results, err := jp.Query(query, jsonData)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
		} else {
			fmt.Printf("  Results: %d\n", len(results))
			for i, r := range results {
				book := r.Value.(map[string]interface{})
				fmt.Printf("    [%d] %s (%s)\n", i, book["title"], book["category"])
			}
		}
		fmt.Println()
	}
}
