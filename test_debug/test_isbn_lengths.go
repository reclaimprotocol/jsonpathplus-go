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
					"price": 8.95,
					"isbn": "0-553-21311-1"
				},
				{
					"category": "fiction",
					"author": "Evelyn Waugh",
					"title": "Sword of Honour", 
					"price": 12.99,
					"isbn": "0-553-21311-2"
				},
				{
					"category": "fiction",
					"author": "Herman Melville",
					"title": "Moby Dick",
					"price": 8.99,
					"isbn": "0-553-21311-3"
				},
				{
					"category": "fiction",
					"author": "J. R. R. Tolkien",
					"title": "The Lord of the Rings",
					"price": 22.99,
					"isbn": "0-395-19395-8"
				}
			]
		}
	}`

	// Check all ISBN lengths
	fmt.Println("=== Checking all ISBN lengths ===")
	results, err := jp.Query("$.store.book[*].isbn", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for i, r := range results {
		isbn := r.Value.(string)
		fmt.Printf("ISBN %d: '%s' (length: %d)\n", i, isbn, len(isbn))
	}

	// Test the filter
	fmt.Println("\n=== Testing filter ===")
	results2, err := jp.Query("$.store.book[?(@.isbn.length === 13)]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Books with 13-char ISBN: %d results\n", len(results2))
	for i, r := range results2 {
		book := r.Value.(map[string]interface{})
		fmt.Printf("  [%d] %s (ISBN: %s)\n", i, book["title"], book["isbn"])
	}
}
