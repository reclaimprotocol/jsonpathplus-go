package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	jsonData := `{
		"store": {
			"book": [
				{"title": "Book 1", "price": 8.95}
			]
		}
	}`
	
	// Test simpler typeof function
	fmt.Println("=== Testing simpler typeof function ===")
	results, err := jp.Query("$.store.book[0][?(@.price.typeof() === 'number')]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Simple typeof results: %d\n", len(results))
}