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
				"items": [
					{"product": "laptop", "price": 999.99, "quantity": 1},
					{"product": "mouse", "price": 29.99, "quantity": 2}
				]
			},
			{
				"id": "ORD002", 
				"items": [
					{"product": "keyboard", "price": 79.99, "quantity": 1}
				]
			},
			{
				"id": "ORD003",
				"items": [
					{"product": "monitor", "price": 299.99, "quantity": 1},
					{"product": "cable", "price": 19.99, "quantity": 3}
				]
			}
		]
	}`

	fmt.Println("=== Debugging Wildcard Issue ===")
	fmt.Printf("JSONPath: $.orders[*].items[*]\n")
	fmt.Println("Expected: 6 items total (2+1+3)")
	fmt.Println("- ORD001: laptop, mouse (2 items)")
	fmt.Println("- ORD002: keyboard (1 item)") 
	fmt.Println("- ORD003: monitor, cable (2 items)")
	fmt.Println("Total: 5 items, but test expects 6")
	
	results, err := jp.Query("$.orders[*].items[*]", jsonData)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("Actual results: %d\n", len(results))
	for i, result := range results {
		fmt.Printf("  [%d] %v\n", i, result.Value)
	}
}