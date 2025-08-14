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

	fmt.Println("=== Comparing different wildcard contexts ===")

	// Test what $.store.book looks like when reached via different paths
	fmt.Println("Direct path: $.store.book")
	results1, err := jp.Query("$.store.book", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		for _, r := range results1 {
			fmt.Printf("Path: '%s'\n", r.Path)
		}
	}

	fmt.Println()
	fmt.Println("Recursive path: $..book")
	results2, err := jp.Query("$..book", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		for _, r := range results2 {
			fmt.Printf("Path: '%s'\n", r.Path)
		}
	}

	fmt.Println()
	fmt.Println("The key insight: Both have the same path '$.store.book'")
	fmt.Println("So the wildcard behavior should be the same for both")
	fmt.Println()

	// Test the working case
	fmt.Println("Working: $..book.*")
	results3, err := jp.Query("$..book.*", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Results: %d\n", len(results3))
		for i, r := range results3 {
			fmt.Printf("  [%d] %v (path: %s)\n", i, r.Value, r.Path)
		}
	}

	fmt.Println()

	// Test the broken case
	fmt.Println("Broken: $.store.book[*]")
	results4, err := jp.Query("$.store.book[*]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Results: %d\n", len(results4))
		for i, r := range results4 {
			fmt.Printf("  [%d] %v (type: %T, path: %s)\n", i, r.Value, r.Value, r.Path)
		}
	}
}
