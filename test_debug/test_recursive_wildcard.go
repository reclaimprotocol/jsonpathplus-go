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
					"category": "fiction",
					"title": "Book 1"
				},
				{
					"category": "action", 
					"title": "Book 2"
				}
			]
		}
	}`
	
	fmt.Println("=== Testing $..book.* (recursive wildcard) ===")
	results1, err := jp.Query("$..book.*", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Found %d results:\n", len(results1))
		for i, r := range results1 {
			fmt.Printf("  [%d] %v (type: %T, path: %s)\n", i, r.Value, r.Value, r.Path)
		}
	}
	
	fmt.Println("\n=== Testing $.store.book[*].* (alternative) ===")
	results2, err := jp.Query("$.store.book[*].*", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Found %d results:\n", len(results2))
		for i, r := range results2 {
			fmt.Printf("  [%d] %v (type: %T, path: %s)\n", i, r.Value, r.Value, r.Path)
		}
	}
	
	fmt.Println("\n=== Testing $.store.book[*].*[?(@property === \"category\")] ===")
	results3, err := jp.Query("$.store.book[*].*[?(@property === \"category\")]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Found %d category properties:\n", len(results3))
		for i, r := range results3 {
			fmt.Printf("  [%d] %v (path: %s)\n", i, r.Value, r.Path)
		}
	}
}