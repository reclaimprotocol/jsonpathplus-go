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

	// Test simple @path exact match
	fmt.Println("=== Testing exact @path match ===")
	results1, err := jp.Query("$.store.book[?(@path === \"$['store']['book'][0]\")]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Exact match for first book: %d results\n", len(results1))
	for i, r := range results1 {
		book := r.Value.(map[string]interface{})
		fmt.Printf("  [%d] %s (path: %s)\n", i, book["title"], r.Path)
	}

	// Test @path inequality
	fmt.Println("\n=== Testing @path inequality ===")
	results2, err := jp.Query("$.store.book[?(@path !== \"$['store']['book'][0]\")]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Not first book: %d results\n", len(results2))
	for i, r := range results2 {
		book := r.Value.(map[string]interface{})
		fmt.Printf("  [%d] %s (path: %s)\n", i, book["title"], r.Path)
	}

	// Test if @path is recognized at all
	fmt.Println("\n=== Testing if @path is working ===")
	results3, err := jp.Query("$.store.book[?(@path)]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Books with any @path: %d results\n", len(results3))
}
