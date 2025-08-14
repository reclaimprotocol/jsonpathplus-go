package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	jsonData := `{
		"store": {
			"book": [
				{"isbn": "0-553-21311-12", "title": "Book 1"},
				{"isbn": "0-553-21311-123", "title": "Book 2"},
				{"isbn": "0-553-21311-1234", "title": "Book 3"}
			]
		}
	}`
	
	// Test basic access first
	fmt.Println("=== Testing basic access ===")
	results, err := jp.Query("$.store.book[*].isbn", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	for i, r := range results {
		fmt.Printf("ISBN %d: %s (length: %d)\n", i, r.Value, len(r.Value.(string)))
	}
	
	// Test length function
	fmt.Println("\n=== Testing length function ===")
	results2, err := jp.Query("$.store.book[?(@.isbn.length === 13)]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Books with 13-char ISBN: %d results\n", len(results2))
	
	// Test alternative syntax
	results3, err := jp.Query("$.store.book[?(@.isbn.length == 13)]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Books with 13-char ISBN (==): %d results\n", len(results3))
	
	// Test with 14-character length
	results4, err := jp.Query("$.store.book[?(@.isbn.length === 14)]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Books with 14-char ISBN: %d results\n", len(results4))
}