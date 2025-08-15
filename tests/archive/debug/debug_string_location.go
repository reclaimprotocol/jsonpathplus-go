package jsonpathplus

import (
	"encoding/json"
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
	"strings"
)

func debug_string_locationMain() {
	data := `{"store":{"book":[{"title":"Book0"},{"title":"Book1"},{"title":"Book2"}]}}`

	fmt.Println("=== Testing Go String Index Locations ===")
	fmt.Printf("JSON Data: %s\n", data)
	fmt.Printf("Data Length: %d characters\n\n", len(data))

	// Test different queries to see string locations
	tests := []string{
		"$.store.book[0]",
		"$.store.book[1]",
		"$.store.book[2]",
		"$.store.book[*]",
		"$.store.book[0].title",
		"$..title",
	}

	for _, query := range tests {
		fmt.Printf("Query: %s\n", query)
		results, err := jp.Query(query, data)
		if err != nil {
			fmt.Printf("  Error: %s\n", err.Error())
		} else {
			fmt.Printf("  Count: %d\n", len(results))
			for i, r := range results {
				valueStr, _ := json.Marshal(r.Value)
				fmt.Printf("  [%d] Path: %s\n", i, r.Path)
				fmt.Printf("      Value: %s\n", string(valueStr))

				// Try to find the string location of this value in the original JSON
				valueInJson := strings.Index(data, string(valueStr))
				if valueInJson != -1 {
					fmt.Printf("      String Location: character %d-%d\n", valueInJson, valueInJson+len(valueStr)-1)
				} else {
					fmt.Printf("      String Location: not found in original JSON\n")
				}

				// Show context around the location
				if valueInJson != -1 && len(valueStr) < 50 {
					start := max(0, valueInJson-10)
					end := min(len(data), valueInJson+len(valueStr)+10)
					context := data[start:end]
					fmt.Printf("      Context: ...%s...\n", context)

					// Show pointer to the exact location
					pointer := strings.Repeat(" ", valueInJson-start) + strings.Repeat("^", len(valueStr))
					fmt.Printf("      Pointer: ...%s...\n", pointer)
				}
			}
		}
		fmt.Println()
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
