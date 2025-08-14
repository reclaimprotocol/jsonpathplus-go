package main

import (
	"fmt"
	
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	data := `{"matrix":[[1,2,3],[4,5,6],[7,8,9]]}`
	
	queries := []string{
		"$.matrix[-1]",
		"$.matrix[-2]",
		"$.matrix[2]",
	}
	
	for _, query := range queries {
		fmt.Printf("\nTesting: %s\n", query)
		results, err := jp.Query(query, data)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			continue
		}
		
		fmt.Printf("Count: %d\n", len(results))
		for i, result := range results {
			fmt.Printf("[%d] Value: %v\n", i, result.Value)
		}
	}
}