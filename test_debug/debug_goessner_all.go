package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	// Exact test data from Goessner spec
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

	fmt.Println("=== Testing $..*===")
	
	results, err := jp.Query("$..*", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	
	fmt.Printf("Total results: %d (expected: 27)\n", len(results))
	
	// Group by path to identify duplicates
	pathCounts := make(map[string]int)
	pathValues := make(map[string][]interface{})
	
	for _, result := range results {
		pathCounts[result.Path]++
		pathValues[result.Path] = append(pathValues[result.Path], result.Value)
	}
	
	fmt.Println("\nDuplicate paths:")
	duplicateCount := 0
	for path, count := range pathCounts {
		if count > 1 {
			fmt.Printf("  %s: %d times - Values: %v\n", path, count, pathValues[path])
			duplicateCount += count - 1 // Extra occurrences
		}
	}
	
	if duplicateCount == 0 {
		fmt.Println("  None found!")
	} else {
		fmt.Printf("\nTotal extra duplicates: %d\n", duplicateCount)
		fmt.Printf("Results without duplicates would be: %d\n", len(results) - duplicateCount)
	}
}