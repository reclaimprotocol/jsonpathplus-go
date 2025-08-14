package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	// Let's check what AST structure is created for different queries
	queries := []string{
		"$..book",
		"$..book.*",
		"$.store.book",
		"$.store.book.*",
		"$.store.book[*].*",
	}

	jsonData := `{
		"store": {
			"book": [
				{"category": "fiction", "title": "Book 1"},
				{"category": "action", "title": "Book 2"}
			]
		}
	}`

	for _, query := range queries {
		fmt.Printf("=== Query: %s ===\n", query)

		results, err := jp.Query(query, jsonData)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Results: %d\n", len(results))
			for i, r := range results {
				fmt.Printf("  [%d] %v (type: %T, path: %s)\n", i, r.Value, r.Value, r.Path)
			}
		}
		fmt.Println()
	}

	fmt.Println("=== Analysis ===")
	fmt.Println("$..book should find the book array")
	fmt.Println("$..book.* should find all properties of all books (individual values)")
	fmt.Println("Currently $..book.* returns book objects instead of properties")
	fmt.Println("Compare with $.store.book[*].* which works correctly")
}
