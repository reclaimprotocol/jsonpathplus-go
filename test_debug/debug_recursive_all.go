package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	// Simplified test data
	jsonData := `{
		"a": {
			"b": 1,
			"c": 2
		},
		"d": [3, 4]
	}`

	fmt.Println("=== Testing $..*===")
	
	results, err := jp.Query("$..*", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	
	fmt.Printf("Total results: %d\n", len(results))
	
	// Group by path to identify duplicates
	pathCounts := make(map[string]int)
	pathValues := make(map[string][]interface{})
	
	for _, result := range results {
		pathCounts[result.Path]++
		pathValues[result.Path] = append(pathValues[result.Path], result.Value)
	}
	
	fmt.Println("\nPath frequencies:")
	duplicateFound := false
	for path, count := range pathCounts {
		if count > 1 {
			fmt.Printf("  %s: %d times ⚠️ DUPLICATE - Values: %v\n", path, count, pathValues[path])
			duplicateFound = true
		} else {
			fmt.Printf("  %s: %d time - Value: %v\n", path, count, pathValues[path][0])
		}
	}
	
	if !duplicateFound {
		fmt.Println("\n✅ No duplicates found!")
	} else {
		fmt.Println("\n❌ Duplicates detected!")
	}
}