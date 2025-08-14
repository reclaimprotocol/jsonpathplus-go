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

	fmt.Println("=== Expected Results ===")
	fmt.Println("Categories ending with 'tion': reference, fiction, action")
	fmt.Println("Query should find these 3 values")

	fmt.Println("\n=== Testing simpler approach ===")
	fmt.Println("Query: $.store.book[*].category[?(@.match(/TION$/i))]")

	results1, err := jp.Query("$.store.book[*].category[?(@.match(/TION$/i))]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Results: %d\n", len(results1))
		for i, r := range results1 {
			fmt.Printf("  [%d] '%s'\n", i, r.Value)
		}
	}

	fmt.Println("\n=== Testing case-sensitive version ===")
	fmt.Println("Query: $.store.book[*].category[?(@.match(/tion$/))]")

	results2, err := jp.Query("$.store.book[*].category[?(@.match(/tion$/))]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Results: %d\n", len(results2))
		for i, r := range results2 {
			fmt.Printf("  [%d] '%s'\n", i, r.Value)
		}
	}

	fmt.Println("\n=== Testing original complex query ===")
	fmt.Println("Query: $..book.*[?(@property === \"category\" && @.match(/TION$/i))]")

	results3, err := jp.Query("$..book.*[?(@property === \"category\" && @.match(/TION$/i))]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Results: %d\n", len(results3))
		for i, r := range results3 {
			fmt.Printf("  [%d] '%s' (path: %s)\n", i, r.Value, r.Path)
		}
	}
}
