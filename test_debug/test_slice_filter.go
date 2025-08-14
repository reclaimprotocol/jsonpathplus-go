package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	jsonData := `{
		"store": {
			"book": [
				{"category": "reference", "title": "Book 1"},
				{"category": "fiction", "title": "Book 2"},
				{"category": "fiction", "title": "Book 3"},
				{"category": "fiction", "title": "Book 4"}
			]
		}
	}`
	
	// Test individual operations first
	fmt.Println("=== Testing individual operations ===")
	
	// Test slice
	results1, err := jp.Query("$.store.book[0:3]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("First 3 books: %d results\n", len(results1))
	for i, r := range results1 {
		book := r.Value.(map[string]interface{})
		fmt.Printf("  [%d] %s (%s)\n", i, book["title"], book["category"])
	}
	
	// Test filter
	results2, err := jp.Query("$.store.book[?(@.category === 'fiction')]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Fiction books: %d results\n", len(results2))
	for i, r := range results2 {
		book := r.Value.(map[string]interface{})
		fmt.Printf("  [%d] %s (%s)\n", i, book["title"], book["category"])
	}
	
	// Test slice then filter
	fmt.Println("\n=== Testing slice then filter ===")
	results3, err := jp.Query("$.store.book[0:3][?(@.category === 'fiction')]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("First 3 books, then filter fiction: %d results\n", len(results3))
	for i, r := range results3 {
		book := r.Value.(map[string]interface{})
		fmt.Printf("  [%d] %s (%s)\n", i, book["title"], book["category"])
	}
	
	// Test filter then slice
	fmt.Println("\n=== Testing filter then slice ===")
	results4, err := jp.Query("$.store.book[?(@.category === 'fiction')][0:2]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Fiction books, then first 2: %d results\n", len(results4))
	for i, r := range results4 {
		book := r.Value.(map[string]interface{})
		fmt.Printf("  [%d] %s (%s)\n", i, book["title"], book["category"])
	}
}