package jsonpathplus

import (
	"encoding/json"
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func debug_go_indexMain() {
	data := `{"store":{"book":[{"title":"Book0"},{"title":"Book1"},{"title":"Book2"}]}}`

	fmt.Println("=== Testing Go String Index Access ===")

	// Test different index access patterns
	tests := []string{
		"$.store.book[0]",         // Direct numeric index
		"$.store.book['0']",       // String index
		"$['store']['book'][0]",   // Bracket notation with numeric
		"$['store']['book']['0']", // Bracket notation with string
		"$.store['book'][0]",      // Mixed notation
	}

	for _, query := range tests {
		fmt.Printf("\nQuery: %s\n", query)
		results, err := jp.Query(query, data)
		if err != nil {
			fmt.Printf("Go Error: %s\n", err.Error())
		} else {
			fmt.Printf("Go Count: %d\n", len(results))
			if len(results) > 0 {
				valueStr, _ := json.Marshal(results[0].Value)
				fmt.Printf("Go Values: [%s]\n", string(valueStr))
				fmt.Printf("Go Paths: [\"%s\"]\n", results[0].Path)
			}
		}
	}

	fmt.Println("\n=== Testing @property context ===")
	// Let's also test what @property values we get

	fmt.Printf("\nQuery: $.store.book[*]\n")
	results, err := jp.Query("$.store.book[*]", data)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	} else {
		fmt.Printf("Count: %d\n", len(results))
		for i, r := range results {
			fmt.Printf("  [%d] Path: %s, ParentProperty: '%s'\n", i, r.Path, r.ParentProperty)
		}
	}
}
