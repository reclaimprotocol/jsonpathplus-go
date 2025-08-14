package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	jsonData := `{
		"store": {
			"book": [
				{"category": "fiction", "price": 8.95, "title": "Book 1"},
				{"category": "fiction", "price": 12.99, "title": "Book 2"},
				{"category": "reference", "price": 22.99, "title": "Book 3"}
			]
		}
	}`
	
	// Test individual conditions first
	fmt.Println("=== Testing individual conditions ===")
	
	results1, err := jp.Query("$.store.book[?(@.category === 'fiction')]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Fiction books: %d results\n", len(results1))
	for i, r := range results1 {
		fmt.Printf("  [%d] %v\n", i, r.Value)
	}
	
	results2, err := jp.Query("$.store.book[?(@.price < 15)]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Books under $15: %d results\n", len(results2))
	for i, r := range results2 {
		fmt.Printf("  [%d] %v\n", i, r.Value)
	}
	
	// Test AND condition
	fmt.Println("\n=== Testing AND condition ===")
	results3, err := jp.Query("$.store.book[?(@.category === 'fiction' && @.price < 15)]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Fiction books under $15: %d results\n", len(results3))
	for i, r := range results3 {
		fmt.Printf("  [%d] %v\n", i, r.Value)
	}
	
	// Test OR condition
	fmt.Println("\n=== Testing OR condition ===")
	results4, err := jp.Query("$.store.book[?(@.price < 10 || @.price > 20)]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Very cheap or expensive books: %d results\n", len(results4))
	for i, r := range results4 {
		fmt.Printf("  [%d] %v\n", i, r.Value)
	}
}