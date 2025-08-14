package main

import (
	"fmt"
	
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	data := `{
		"orders": [
			{"id": "ORD001", "items": [{"product": "laptop", "price": 999.99}, {"product": "mouse", "price": 29.99}]},
			{"id": "ORD002", "items": [{"product": "keyboard", "price": 79.99}]},
			{"id": "ORD003", "items": [{"product": "laptop", "price": 1299.99}]}
		]
	}`
	
	query := `$.orders[?(@.items[?(@.product === 'laptop')])]`
	fmt.Printf("Testing: %s\n", query)
	
	results, err := jp.Query(query, data)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Count: %d\n", len(results))
		for i, result := range results {
			if order, ok := result.Value.(map[string]interface{}); ok {
				fmt.Printf("[%d] ID: %v\n", i, order["id"])
			}
		}
	}
}