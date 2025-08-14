package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	jsonData := `{
		"products": [
			{"id": 1, "name": "Widget", "price": 10.0, "inStock": true},
			{"id": 2, "name": "Gadget", "price": 25.0, "inStock": false},
			{"id": 3, "name": "Tool", "price": 15.0, "inStock": true},
			{"id": 4, "name": "Device", "price": 30.0, "inStock": false}
		]
	}`
	
	// Test boolean existence filter
	fmt.Println("=== Testing boolean existence filter ===")
	results, err := jp.Query("$.products[?(@.inStock)]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Products in stock: %d results\n", len(results))
	
	for i, r := range results {
		product := r.Value.(map[string]interface{})
		fmt.Printf("  [%d] %s (inStock: %v)\n", i, product["name"], product["inStock"])
	}
	
	// Test negated boolean existence filter
	fmt.Println("\n=== Testing negated boolean existence filter ===")
	results2, err := jp.Query("$.products[?(!@.inStock)]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Products not in stock: %d results\n", len(results2))
	
	for i, r := range results2 {
		product := r.Value.(map[string]interface{})
		fmt.Printf("  [%d] %s (inStock: %v)\n", i, product["name"], product["inStock"])
	}
}