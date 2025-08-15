package jsonpathplus

import (
	"encoding/json"
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func debug_propertyMain() {
	data := `{"store":{"book":[{"category":"reference","author":"Nigel Rees","title":"Sayings of the Century","price":8.95},{"category":"fiction","author":"Evelyn Waugh","title":"Sword of Honour","price":12.99},{"category":"fiction","author":"Herman Melville","title":"Moby Dick","isbn":"0-553-21311-3","price":8.99},{"category":"fiction","author":"J. R. R. Tolkien","title":"The Lord of the Rings","isbn":"0-395-19395-8","price":22.99}],"bicycle":{"color":"red","price":19.95}}}`

	fmt.Println("=== Testing Property Filter Queries ===")

	// Test different property filter queries
	queries := []string{
		"$..book[*]",                                   // All books to see @property values
		"$..book[?(@property !== 0)]",                  // The failing query
		"$..book[?(@property != '0')]",                 // Alternative with string
		"$..book[?(@property > 0)]",                    // Alternative comparison
		"$..*[?(@property === 'price' && @ !== 8.95)]", // The other failing query
	}

	for _, query := range queries {
		fmt.Printf("\nQuery: %s\n", query)
		results, err := jp.Query(query, data)
		if err != nil {
			fmt.Printf("Go Error: %s\n", err.Error())
		} else {
			fmt.Printf("Go Count: %d\n", len(results))
			fmt.Printf("Go Paths: [")
			for i, r := range results {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("\"%s\"", r.Path)
			}
			fmt.Println("]")

			if len(results) > 0 && len(results) <= 5 {
				fmt.Printf("Go Values: [")
				for i, r := range results {
					if i > 0 {
						fmt.Print(", ")
					}
					valueStr, _ := json.Marshal(r.Value)
					fmt.Print(string(valueStr))
				}
				fmt.Println("]")
			}

			// For the book wildcard query, show @property values
			if query == "$..book[*]" {
				for i, r := range results {
					fmt.Printf("  [%d] Path: %s, ParentProperty: '%s'\n", i, r.Path, r.ParentProperty)
				}
			}
		}
	}
}
