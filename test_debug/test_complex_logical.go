package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	jsonData := `{
		"store": {
			"book": [
				{"category": "fiction", "price": 8.95, "title": "Book 1"},
				{"category": "fiction", "price": 12.99, "title": "Book 2"},
				{"category": "fiction", "price": 22.99, "title": "Book 3"}
			]
		},
		"users": [
			{"name": "Alice", "active": true, "age": 30},
			{"name": "Bob", "active": true, "age": 35},
			{"name": "Charlie", "active": false, "age": 40}
		]
	}`

	// Test parts of complex expression
	fmt.Println("=== Testing parts of complex expression ===")

	// Test the OR part alone
	results1, err := jp.Query("$.store.book[?(@.price < 10 || @.price > 20)]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Books cheap or expensive: %d results\n", len(results1))
	for i, r := range results1 {
		book := r.Value.(map[string]interface{})
		fmt.Printf("  [%d] %s ($%.2f)\n", i, book["title"], book["price"])
	}

	// Test with parentheses
	results2, err := jp.Query("$.store.book[?(@.price < 10 || @.price > 20)]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Books cheap or expensive (with parens): %d results\n", len(results2))

	// Test the full complex expression
	fmt.Println("\n=== Testing full complex expression ===")
	results3, err := jp.Query("$.store.book[?(@.category === 'fiction' && (@.price < 10 || @.price > 20))]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Fiction books cheap or expensive: %d results\n", len(results3))
	for i, r := range results3 {
		book := r.Value.(map[string]interface{})
		fmt.Printf("  [%d] %s ($%.2f)\n", i, book["title"], book["price"])
	}

	// Test the users case
	fmt.Println("\n=== Testing users case ===")
	results4, err := jp.Query("$.users[?(@.active === true)]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Active users: %d results\n", len(results4))

	results5, err := jp.Query("$.users[?(@.age > 30)]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Users over 30: %d results\n", len(results5))

	results6, err := jp.Query("$.users[?(@.active === true && @.age > 30)]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Active users over 30: %d results\n", len(results6))
	for i, r := range results6 {
		user := r.Value.(map[string]interface{})
		fmt.Printf("  [%d] %s (age: %.0f, active: %v)\n", i, user["name"], user["age"], user["active"])
	}
}
