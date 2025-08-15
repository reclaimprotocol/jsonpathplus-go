package jsonpathplus

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func debug_array_wildcardMain() {
	// Test the array wildcard specifically
	orderData := `{
		"orders": [
			{"id": "ORD001", "items": [{"product": "laptop", "price": 999.99}, {"product": "mouse", "price": 29.99}]},
			{"id": "ORD002", "items": [{"product": "keyboard", "price": 79.99}]},
			{"id": "ORD003", "items": [{"product": "laptop", "price": 1299.99}]}
		]
	}`

	queries := []string{
		// Array wildcard variations
		`$.orders[?(@.items[*].product === "laptop")]`,
		`$.orders[?(@.items[*].product == "laptop")]`,
		`$.orders[?(@.items[0].product === "laptop")]`, // Test specific index

		// These should also work
		`$.orders[?(@.items[*].price > 1000)]`,
		`$.orders[?(@.items[0].price > 500)]`,
	}

	for i, query := range queries {
		fmt.Printf("Test %d: %s\n", i+1, query)
		results, err := jp.Query(query, orderData)
		if err != nil {
			fmt.Printf("✗ Error: %v\n", err)
		} else {
			fmt.Printf("✓ Results: %d orders", len(results))
			if len(results) > 0 {
				fmt.Printf(" (")
				for j, result := range results {
					if order, ok := result.Value.(map[string]interface{}); ok {
						if j > 0 {
							fmt.Printf(", ")
						}
						fmt.Printf("%v", order["id"])
					}
				}
				fmt.Printf(")")
			}
			fmt.Printf("\n")
		}
		fmt.Println()
	}
}
