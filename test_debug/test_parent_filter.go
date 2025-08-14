package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	jsonData := `{
		"store": {
			"book": [
				{"category": "reference", "author": "Nigel Rees", "title": "Sayings of the Century", "price": 8.95},
				{"category": "fiction", "author": "Evelyn Waugh", "title": "Sword of Honour", "price": 12.99}
			],
			"bicycle": {
				"color": "red",
				"price": 19.95,
				"manufacturer": "Trek"
			}
		}
	}`
	
	// Test basic array access first
	fmt.Println("=== Testing basic access ===")
	results1, err := jp.Query("$.store.book[*]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("All books: %d results\n", len(results1))
	for i, r := range results1 {
		book := r.Value.(map[string]interface{})
		fmt.Printf("  [%d] %s (parent: %v)\n", i, book["title"], r.Parent != nil)
	}
	
	// Test simple filter first
	fmt.Println("\n=== Testing simple filter ===")
	results2, err := jp.Query("$.store.book[?(@.price < 15)]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Cheap books: %d results\n", len(results2))
	for i, r := range results2 {
		book := r.Value.(map[string]interface{})
		fmt.Printf("  [%d] %s\n", i, book["title"])
	}
	
	// Test @parent existence
	fmt.Println("\n=== Testing @parent filter ===")
	results3, err := jp.Query("$.store.book[?(@parent)]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Books with parent: %d results\n", len(results3))
	
	// Test @parent.bicycle
	fmt.Println("\n=== Testing @parent.bicycle ===")
	results4, err := jp.Query("$.store.book[?(@parent.bicycle)]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Books where parent has bicycle: %d results\n", len(results4))
	
	// Test @parent.bicycle.color
	fmt.Println("\n=== Testing @parent.bicycle.color ===")
	results5, err := jp.Query("$.store.book[?(@parent.bicycle.color === 'red')]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Books where parent has red bicycle: %d results\n", len(results5))
}