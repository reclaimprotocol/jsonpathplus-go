package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	jsonData := `{
		"store": {
			"book": [
				{"category": "fiction", "title": "Book 1"},
				{"category": "action", "title": "Book 2"}
			]
		}
	}`
	
	fmt.Println("=== Testing [*] pattern ===")
	
	// Test [*] expansion
	fmt.Println("1. $.store.book[*] (expand array elements)")
	results1, err := jp.Query("$.store.book[*]", jsonData)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Results: %d\n", len(results1))
		for i, r := range results1 {
			fmt.Printf("   [%d] %v (type: %T, path: %s)\n", i, r.Value, r.Value, r.Path)
		}
	}
	
	fmt.Println()
	
	// Test [*].* chaining
	fmt.Println("2. $.store.book[*].* (expand array, then get properties)")
	results2, err := jp.Query("$.store.book[*].*", jsonData)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Results: %d\n", len(results2))
		for i, r := range results2 {
			fmt.Printf("   [%d] %v (type: %T, path: %s)\n", i, r.Value, r.Value, r.Path)
		}
	}
	
	fmt.Println()
	
	// Test direct wildcard
	fmt.Println("3. $.store.book.* (direct wildcard on array)")
	results3, err := jp.Query("$.store.book.*", jsonData)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Results: %d\n", len(results3))
		for i, r := range results3 {
			fmt.Printf("   [%d] %v (type: %T, path: %s)\n", i, r.Value, r.Value, r.Path)
		}
	}
	
	fmt.Println()
	fmt.Println("=== Analysis ===")
	fmt.Println("Both patterns should now return the same results:")
	fmt.Println("- $.store.book[*].* should expand array then get properties")
	fmt.Println("- $.store.book.* should directly get properties from array elements")
}