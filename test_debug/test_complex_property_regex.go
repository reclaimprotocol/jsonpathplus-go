package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	// Test for: $..book.*[?(@property === "category" && @.match(/TION$/i))]
	// This should find all book properties where:
	// 1. The property name is "category" 
	// 2. The property value matches the regex /TION$/i (ends with "tion", case-insensitive)
	
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
	
	fmt.Println("=== Complex Property Regex Filter Test ===")
	fmt.Println("Query: $.store.book[*].*[?(@property === \"category\" && @.match(/TION$/i))]")
	fmt.Println()
	fmt.Println("This query finds all book properties where:")
	fmt.Println("1. Property name equals 'category'")
	fmt.Println("2. Property value ends with 'tion' (case-insensitive)")
	fmt.Println()
	fmt.Println("Expected results:")
	fmt.Println("- 'fiction' (ends with 'tion')")
	fmt.Println("- 'action' (ends with 'tion')")
	fmt.Println("- NOT 'reference' (ends with 'ence')")
	fmt.Println("- NOT 'science' (ends with 'ence')")
	fmt.Println()
	
	// Test the working query
	results, err := jp.Query("$.store.book[*].*[?(@property === \"category\" && @.match(/TION$/i))]", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	
	fmt.Printf("📊 Results: %d matches\n", len(results))
	
	if len(results) == 2 {
		fmt.Println("✅ Correct number of results!")
	} else {
		fmt.Printf("❌ Expected 2 results, got %d\n", len(results))
	}
	
	expectedValues := map[string]bool{"fiction": false, "action": false}
	
	for i, r := range results {
		value := r.Value.(string)
		fmt.Printf("  [%d] '%s' (path: %s)\n", i, value, r.Path)
		
		if _, expected := expectedValues[value]; expected {
			expectedValues[value] = true
			fmt.Printf("      ✅ Expected value\n")
		} else {
			fmt.Printf("      ❌ Unexpected value\n")
		}
	}
	
	// Check if all expected values were found
	fmt.Println()
	allFound := true
	for value, found := range expectedValues {
		if !found {
			fmt.Printf("❌ Missing expected value: '%s'\n", value)
			allFound = false
		}
	}
	
	if allFound {
		fmt.Println("✅ All expected values found!")
		fmt.Println("✅ Complex property regex filter test PASSED!")
	} else {
		fmt.Println("❌ Some expected values missing!")
		fmt.Println("❌ Complex property regex filter test FAILED!")
	}
}