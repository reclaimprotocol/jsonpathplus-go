package jsonpathplus

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func debug_simple_propertyMain() {
	data := `{"book":[{"title":"Book0"},{"title":"Book1"}]}`
	
	fmt.Println("=== Testing Simple Array Property Filtering ===")
	
	queries := []string{
		"$.book[*]",                    // All books to see context
		"$.book[?(@property === 0)]",   // Should return book 0
		"$.book[?(@property === 1)]",   // Should return book 1
		"$.book[?(@property !== 0)]",   // Should return book 1
	}
	
	for _, query := range queries {
		fmt.Printf("\nQuery: %s\n", query)
		results, err := jp.Query(query, data)
		if err != nil {
			fmt.Printf("Go Error: %s\n", err.Error())
		} else {
			fmt.Printf("Go Count: %d\n", len(results))
			for i, r := range results {
				fmt.Printf("  [%d] Path: %s, ParentProperty: '%s'\n", i, r.Path, r.ParentProperty)
			}
		}
	}
}