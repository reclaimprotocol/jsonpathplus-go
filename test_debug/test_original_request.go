package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	// Test the original requested query
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
	
	fmt.Println("=== ORIGINAL REQUESTED QUERY TEST ===")
	fmt.Println()
	fmt.Println("Query: $..book.*[?(@property === \"category\" && @.match(/TION$/i))]")
	fmt.Println("Expected: Categories ending with 'tion' (case-insensitive)")
	fmt.Println("Should find: 'fiction', 'action'")
	fmt.Println("Should NOT find: 'reference' (ends with 'ence'), 'science' (ends with 'ence')")
	fmt.Println()
	
	results, err := jp.Query("$..book.*[?(@property === \"category\" && @.match(/TION$/i))]", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	
	fmt.Printf("📊 Results: %d\n", len(results))
	
	expectedResults := []string{"fiction", "action"}
	actualResults := make([]string, 0, len(results))
	
	for i, r := range results {
		value := r.Value.(string)
		actualResults = append(actualResults, value)
		fmt.Printf("  [%d] '%s' (path: %s)\n", i, value, r.Path)
	}
	
	fmt.Println()
	
	// Verify results
	success := len(results) == len(expectedResults)
	if success {
		for _, expected := range expectedResults {
			found := false
			for _, actual := range actualResults {
				if actual == expected {
					found = true
					break
				}
			}
			if !found {
				success = false
				break
			}
		}
	}
	
	if success {
		fmt.Println("✅ SUCCESS: Original query works perfectly!")
		fmt.Println("✅ Found exactly the expected results")
		fmt.Println("✅ Complex property regex filter is working")
		fmt.Println("✅ @property filter is working")
		fmt.Println("✅ Case-insensitive regex with /TION$/i is working")
		fmt.Println("✅ Logical AND operator (&&) is working")
		fmt.Println("✅ Recursive descent with property wildcard (..*) is working")
	} else {
		fmt.Println("❌ FAILED: Results don't match expected")
	}
}