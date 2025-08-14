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

	// Get the actual paths to see what they look like
	fmt.Println("=== Actual paths from query ===")
	results, err := jp.Query("$.store.book[*]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for i, r := range results {
		book := r.Value.(map[string]interface{})
		fmt.Printf("Book %d: %s\n", i, book["title"])
		fmt.Printf("  Path: %s\n", r.Path)

		// Test what bracket conversion would give us
		fmt.Printf("  Expected for comparison: $['store']['book'][%d]\n", i)
		fmt.Println()
	}

	// Test each path individually
	fmt.Println("=== Testing individual paths ===")

	// Test exact match for first book
	exactQuery := "$.store.book[?(@path === \"$['store']['book'][0]\")]"
	fmt.Printf("Query: %s\n", exactQuery)
	results1, err := jp.Query(exactQuery, jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Results: %d\n", len(results1))
	}

	// Test exact match for second book
	exactQuery2 := "$.store.book[?(@path === \"$['store']['book'][1]\")]"
	fmt.Printf("Query: %s\n", exactQuery2)
	results2, err := jp.Query(exactQuery2, jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Results: %d\n", len(results2))
	}
}
