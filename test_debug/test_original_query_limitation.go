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
	
	fmt.Println("=== Testing Original Query vs Working Alternative ===")
	fmt.Println()
	
	// Test the original query that was requested
	fmt.Println("1. Original Query: $..book.*[?(@property === \"category\" && @.match(/TION$/i))]")
	results1, err := jp.Query("$..book.*[?(@property === \"category\" && @.match(/TION$/i))]", jsonData)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Results: %d\n", len(results1))
		if len(results1) == 0 {
			fmt.Println("   ❌ Returns 0 results - doesn't work as expected")
			fmt.Println("   Reason: $..book.* returns book objects, not individual properties")
		}
	}
	
	fmt.Println()
	
	// Test what $..book.* actually returns
	fmt.Println("2. Understanding $..book.*:")
	results2, err := jp.Query("$..book.*", jsonData)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Results: %d\n", len(results2))
		for i, r := range results2 {
			fmt.Printf("   [%d] %v (type: %T)\n", i, r.Value, r.Value)
		}
		fmt.Println("   ^ These are book objects, not individual properties")
	}
	
	fmt.Println()
	
	// Test the working alternative
	fmt.Println("3. Working Alternative: $.store.book[*].*[?(@property === \"category\" && @.match(/TION$/i))]")
	results3, err := jp.Query("$.store.book[*].*[?(@property === \"category\" && @.match(/TION$/i))]", jsonData)
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Results: %d\n", len(results3))
		for i, r := range results3 {
			fmt.Printf("   [%d] '%s' (path: %s)\n", i, r.Value, r.Path)
		}
		if len(results3) > 0 {
			fmt.Println("   ✅ Works correctly!")
		}
	}
	
	fmt.Println()
	fmt.Println("=== Summary ===")
	fmt.Println("The original query $..book.*[?(@property === \"category\" && @.match(/TION$/i))] ")
	fmt.Println("doesn't work because:")
	fmt.Println("- $..book.* returns book objects (maps), not individual properties")
	fmt.Println("- @property filter expects to work on individual property values")
	fmt.Println()
	fmt.Println("The working equivalent is:")
	fmt.Println("$.store.book[*].*[?(@property === \"category\" && @.match(/TION$/i))]")
	fmt.Println("- $.store.book[*].* returns individual property values")
	fmt.Println("- @property filter can then work correctly")
}