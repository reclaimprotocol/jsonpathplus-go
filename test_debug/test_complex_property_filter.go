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
					"price": 8.95,
					"isbn": "0-553-21311-1"
				},
				{
					"category": "fiction",
					"author": "Evelyn Waugh",
					"title": "Sword of Honour", 
					"price": 12.99,
					"isbn": "0-553-21311-2"
				},
				{
					"category": "action",
					"author": "Herman Melville",
					"title": "Moby Dick",
					"price": 8.99,
					"isbn": "0-553-21311-3"
				},
				{
					"category": "science",
					"author": "J. R. R. Tolkien",
					"title": "The Lord of the Rings",
					"price": 22.99,
					"isbn": "0-395-19395-8"
				}
			]
		}
	}`
	
	fmt.Println("=== Test Data Categories ===")
	// First, let's see all categories
	results, err := jp.Query("$..book[*].category", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	for i, r := range results {
		fmt.Printf("Category %d: '%s'\n", i, r.Value)
	}
	
	fmt.Println("\n=== Testing complex property filter ===")
	fmt.Println("Query: $..book.*[?(@property === \"category\" && @.match(/TION$/i))]")
	fmt.Println("Should find: 'reference', 'fiction', 'action' (all ending with 'tion')")
	
	// Test the complex expression
	results2, err := jp.Query("$..book.*[?(@property === \"category\" && @.match(/TION$/i))]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("Results: %d\n", len(results2))
	for i, r := range results2 {
		fmt.Printf("  [%d] '%s' (path: %s)\n", i, r.Value, r.Path)
	}
	
	fmt.Println("\n=== Testing individual components ===")
	
	// Test @property filter alone
	fmt.Println("1. Testing @property === \"category\":")
	results3, err := jp.Query("$..book.*[?(@property === \"category\")]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("   Found %d category properties\n", len(results3))
		for i, r := range results3 {
			fmt.Printf("   [%d] '%s'\n", i, r.Value)
		}
	}
	
	// Test regex match alone on categories
	fmt.Println("\n2. Testing @.match(/TION$/i) on categories:")
	results4, err := jp.Query("$..book[*].category[?(@.match(/TION$/i))]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("   Found %d categories ending with 'tion'\n", len(results4))
		for i, r := range results4 {
			fmt.Printf("   [%d] '%s'\n", i, r.Value)
		}
	}
}