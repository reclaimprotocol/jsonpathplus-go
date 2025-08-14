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
					"category": "action",
					"author": "Herman Melville",
					"title": "Moby Dick",
					"price": 8.99
				},
				{
					"category": "science",
					"author": "J. R. R. Tolkien",
					"title": "The Lord of the Rings",
					"price": 22.99
				}
			]
		}
	}`
	
	fmt.Println("=== Testing the working complex query ===")
	fmt.Println("Query: $.store.book[*].*[?(@property === \"category\" && @.match(/TION$/i))]")
	fmt.Println("Expected: 'fiction', 'action' (2 results)")
	
	results, err := jp.Query("$.store.book[*].*[?(@property === \"category\" && @.match(/TION$/i))]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("Results: %d\n", len(results))
	for i, r := range results {
		fmt.Printf("  [%d] '%s' (path: %s)\n", i, r.Value, r.Path)
	}
	
	fmt.Println("\n=== Testing the original (problematic) query ===")
	fmt.Println("Query: $..book.*[?(@property === \"category\" && @.match(/TION$/i))]")
	fmt.Println("This returns 0 because $..book.* returns book objects, not properties")
	
	results2, err := jp.Query("$..book.*[?(@property === \"category\" && @.match(/TION$/i))]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("Results: %d\n", len(results2))
	for i, r := range results2 {
		fmt.Printf("  [%d] '%s' (path: %s)\n", i, r.Value, r.Path)
	}
	
	fmt.Println("\n=== Summary ===")
	fmt.Println("The original query $..book.*[?(@property === \"category\" && @.match(/TION$/i))] ")
	fmt.Println("doesn't work because $..book.* returns book objects, not individual properties.")
	fmt.Println("The working equivalent is: $.store.book[*].*[?(@property === \"category\" && @.match(/TION$/i))]")
}