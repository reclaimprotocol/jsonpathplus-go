package jsonpathplus

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func debug_nested_detailedMain() {
	orderData := `{
		"orders": [
			{"id": "ORD001", "items": [{"product": "laptop", "price": 999.99}, {"product": "mouse", "price": 29.99}]},
			{"id": "ORD002", "items": [{"product": "keyboard", "price": 79.99}]},
			{"id": "ORD003", "items": [{"product": "laptop", "price": 1299.99}]}
		]
	}`
	
	// First, let's test a simpler nested approach to verify the data
	fmt.Println("=== Testing data structure ===")
	simpleQuery := `$.orders[*].items[*].product`
	results, err := jp.Query(simpleQuery, orderData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("All products: %d results\n", len(results))
		for i, result := range results {
			fmt.Printf("  %d: %v\n", i+1, result.Value)
		}
	}
	
	// Test array wildcard filter (should work)
	fmt.Println("\n=== Testing array wildcard filter ===")
	wildcardQuery := `$.orders[?(@.items[*].product === "laptop")]`
	results, err = jp.Query(wildcardQuery, orderData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Wildcard results: %d orders\n", len(results))
		for i, result := range results {
			if order, ok := result.Value.(map[string]interface{}); ok {
				fmt.Printf("  %d: %v\n", i+1, order["id"])
			}
		}
	}
	
	// Test the nested filter we're trying to fix
	fmt.Println("\n=== Testing nested filter ===")
	nestedQuery := `$.orders[?(@.items[?(@.product === "laptop")])]`
	results, err = jp.Query(nestedQuery, orderData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Nested results: %d orders\n", len(results))
		for i, result := range results {
			if order, ok := result.Value.(map[string]interface{}); ok {
				fmt.Printf("  %d: %v\n", i+1, order["id"])
			}
		}
	}
	
	// Let's also test the JavaScript version for comparison
	fmt.Println("\n=== Expected JavaScript behavior ===")
	fmt.Println("JavaScript should return 2 orders: ORD001 and ORD003")
	fmt.Println("Both contain items where product === 'laptop'")
}