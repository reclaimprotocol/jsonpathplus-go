package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	jsonData := `{"store": {"book": [{"title": "Book 1", "price": 8.95}, {"title": "Book 2", "price": 12.99}]}}`

	// Test simple property access first
	results, err := jp.Query("$.store.book[0]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("$.store.book[0] found %d results:\n", len(results))
	for i, r := range results {
		fmt.Printf("  [%d] %v (path: %s, parent: %v)\n", i, r.Value, r.Path, r.Parent != nil)
	}

	fmt.Println()

	// Test parent operator
	results2, err := jp.Query("$.store.book[0]^", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("$.store.book[0]^ found %d results:\n", len(results2))
	for i, r := range results2 {
		fmt.Printf("  [%d] %v (path: %s)\n", i, r.Value, r.Path)
	}
}
