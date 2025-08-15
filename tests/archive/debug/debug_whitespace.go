package jsonpathplus

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
	"strings"
)

func debug_whitespaceMain() {
	testData := `{
		"users": [
			{"name": "Alice", "age": 25, "active": true},
			{"name": "Bob", "age": 30, "active": false}, 
			{"name": "Charlie", "age": 35, "active": true}
		]
	}`

	queries := []string{
		`$.users[?(@.name === "Alice")]`,   // No spaces
		`$.users[?(@.name==="Alice")]`,     // No spaces at all
		`$.users[?(@.name === "Alice")]`,   // Normal spaces
		`$.users[?(@.name  ===  "Alice")]`, // Extra spaces
		`$.users[?( @.name === "Alice" )]`, // Spaces around expression
		`$.users[?(@.name === 'Alice')]`,   // Single quotes
	}

	fmt.Println("Testing whitespace sensitivity in Go JSONPath implementation:")
	fmt.Println(strings.Repeat("=", 60))

	for i, query := range queries {
		fmt.Printf("\nTest %d: %s\n", i+1, query)
		results, err := jp.Query(query, testData)
		if err != nil {
			fmt.Printf("✗ Error: %s\n", err.Error())
		} else {
			name := "none"
			if len(results) > 0 {
				if user, ok := results[0].Value.(map[string]interface{}); ok {
					if n, exists := user["name"]; exists {
						name = fmt.Sprintf("%v", n)
					}
				}
			}
			fmt.Printf("✓ Results: %d (%s)\n", len(results), name)
		}
	}

	// Test the specific failing nested filter
	fmt.Printf("\n%s\n", strings.Repeat("=", 60))
	fmt.Println("Testing nested filter:")
	fmt.Printf("%s\n", strings.Repeat("=", 60))

	orderData := `{
		"orders": [
			{"id": "ORD001", "items": [{"product": "laptop", "price": 999.99}, {"product": "mouse", "price": 29.99}]},
			{"id": "ORD002", "items": [{"product": "keyboard", "price": 79.99}]},
			{"id": "ORD003", "items": [{"product": "laptop", "price": 1299.99}]}
		]
	}`

	nestedQuery := `$.orders[?(@.items[?(@.product === "laptop")])]`
	fmt.Printf("\nNested Test: %s\n", nestedQuery)
	results, err := jp.Query(nestedQuery, orderData)
	if err != nil {
		fmt.Printf("✗ Error: %s\n", err.Error())
	} else {
		fmt.Printf("✓ Results: %d orders\n", len(results))
		if len(results) > 0 {
			if order, ok := results[0].Value.(map[string]interface{}); ok {
				if id, exists := order["id"]; exists {
					fmt.Printf("  First result ID: %v\n", id)
				}
			}
		}
	}
}
