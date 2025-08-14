package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	jsonData := `{
		"store": {
			"book": [
				{"title": "Book 1"},
				{"title": "Book 2"}
			]
		}
	}`
	
	// Test current path format
	fmt.Println("=== Current path formats ===")
	results, err := jp.Query("$.store.book[*]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	for i, r := range results {
		fmt.Printf("Book %d path: %s\n", i, r.Path)
	}
	
	// Test @path filter with current format
	fmt.Println("\n=== Testing @path with current format ===")
	results2, err := jp.Query("$.store.book[?(@path !== '$.store.book[0]')]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Books excluding first (current format): %d results\n", len(results2))
	
	// Test @path filter with bracket format
	fmt.Println("\n=== Testing @path with bracket format ===")
	results3, err := jp.Query("$.store.book[?(@path !== \"$['store']['book'][0]\")]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Books excluding first (bracket format): %d results\n", len(results3))
}