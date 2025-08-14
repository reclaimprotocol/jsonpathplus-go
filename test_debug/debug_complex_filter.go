package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	jsonData := `{
		"orders": [
			{
				"id": "ORD001",
				"customer": {"name": "Alice", "type": "premium"},
				"items": [
					{"product": "laptop", "price": 999.99, "quantity": 1},
					{"product": "mouse", "price": 29.99, "quantity": 2}
				],
				"status": "shipped",
				"total": 1059.97
			},
			{
				"id": "ORD002", 
				"customer": {"name": "Bob", "type": "regular"},
				"items": [
					{"product": "keyboard", "price": 79.99, "quantity": 1}
				],
				"status": "pending",
				"total": 79.99
			},
			{
				"id": "ORD003",
				"customer": {"name": "Charlie", "type": "premium"},
				"items": [
					{"product": "monitor", "price": 299.99, "quantity": 1},
					{"product": "cable", "price": 19.99, "quantity": 3}
				],
				"status": "shipped",
				"total": 359.97
			}
		]
	}`

	fmt.Println("=== Debugging Complex Filter ===")
	fmt.Printf("JSONPath: $.orders[?(@.customer.type === 'premium' && @.status === 'shipped')]\n")
	fmt.Println("Expected: 2 results (ORD001 and ORD003 - both premium and shipped)")
	fmt.Println("- ORD001: customer.type=premium, status=shipped ✓")
	fmt.Println("- ORD002: customer.type=regular, status=pending ✗")
	fmt.Println("- ORD003: customer.type=premium, status=shipped ✓")
	
	results, err := jp.Query("$.orders[?(@.customer.type === 'premium' && @.status === 'shipped')]", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("Actual results: %d\n", len(results))
	for i, result := range results {
		fmt.Printf("  [%d] %v\n", i, result.Value)
	}
	
	// Test individual parts
	fmt.Println("\n=== Testing individual parts ===")
	
	// Test customer.type === 'premium'
	results1, err := jp.Query("$.orders[?(@.customer.type === 'premium')]", jsonData)
	if err != nil {
		fmt.Printf("❌ Error testing customer.type: %v\n", err)
	} else {
		fmt.Printf("customer.type === 'premium': %d results\n", len(results1))
	}
	
	// Test status === 'shipped'
	results2, err := jp.Query("$.orders[?(@.status === 'shipped')]", jsonData)
	if err != nil {
		fmt.Printf("❌ Error testing status: %v\n", err)
	} else {
		fmt.Printf("status === 'shipped': %d results\n", len(results2))
	}
}