package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	jsonData := `{
		"store": {
			"book": [
				{
					"category": "reference",
					"author": "Nigel Rees",
					"title": "Sayings of the Century",
					"price": 8.95
				},
				{
					"category": "fiction",
					"author": "Evelyn Waugh",
					"title": "Sword of Honour",
					"price": 12.99
				},
				{
					"category": "fiction",
					"author": "Herman Melville",
					"title": "Moby Dick",
					"isbn": "0-553-21311-3",
					"price": 8.99
				},
				{
					"category": "fiction",
					"author": "J. R. R. Tolkien",
					"title": "The Lord of the Rings",
					"isbn": "0-395-19395-8",
					"price": 22.99
				}
			],
			"bicycle": {
				"color": "red",
				"price": 19.95
			}
		}
	}`

	fmt.Println("=== Debugging @property filter ===")
	fmt.Printf("JSONPath: $..*[?(@property === 'price' && @ !== 8.95)]\n")
	fmt.Println("Expected: 4 price values not equal to 8.95")
	fmt.Println("Prices in data: 8.95, 12.99, 8.99, 22.99, 19.95")
	fmt.Println("Should match: 12.99, 8.99, 22.99, 19.95 (4 values)")
	
	results, err := jp.Query("$..*[?(@property === 'price' && @ !== 8.95)]", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("\nActual results: %d\n", len(results))
	for i, result := range results {
		fmt.Printf("  [%d] %v (path: %s)\n", i, result.Value, result.Path)
	}
	
	// Test simpler version
	fmt.Println("\n=== Testing simpler version ===")
	fmt.Printf("JSONPath: $..*[?(@property === 'price')]\n")
	
	results2, err := jp.Query("$..*[?(@property === 'price')]", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("Results: %d\n", len(results2))
	for i, result := range results2 {
		fmt.Printf("  [%d] %v (path: %s)\n", i, result.Value, result.Path)
	}
}