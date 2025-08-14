package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	// Use the exact test data from complete_feature_test.go
	jsonData := `{
		"users": [
			{"name": "Alice Johnson", "email": "alice@example.com", "age": 30, "active": true},
			{"name": "Bob Smith", "email": "bob@test.org", "age": 25, "active": false},
			{"name": "Charlie Brown", "email": "charlie@example.com", "age": 35, "active": true}
		]
	}`
	
	fmt.Println("=== Users data ===")
	results, err := jp.Query("$.users[*]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	for i, r := range results {
		user := r.Value.(map[string]interface{})
		fmt.Printf("  [%d] %s: age=%.0f, active=%v\n", i, user["name"], user["age"], user["active"])
	}
	
	fmt.Println("\n=== Individual conditions ===")
	
	results1, err := jp.Query("$.users[?(@.active === true)]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Active users: %d\n", len(results1))
	for i, r := range results1 {
		user := r.Value.(map[string]interface{})
		fmt.Printf("  [%d] %s\n", i, user["name"])
	}
	
	results2, err := jp.Query("$.users[?(@.age > 30)]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Users over 30: %d\n", len(results2))
	for i, r := range results2 {
		user := r.Value.(map[string]interface{})
		fmt.Printf("  [%d] %s\n", i, user["name"])
	}
	
	fmt.Println("\n=== Combined condition ===")
	results3, err := jp.Query("$.users[?(@.active === true && @.age > 30)]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Active users over 30: %d\n", len(results3))
	for i, r := range results3 {
		user := r.Value.(map[string]interface{})
		fmt.Printf("  [%d] %s\n", i, user["name"])
	}
	
	// Maybe the test wants >= 30 instead of > 30?
	fmt.Println("\n=== Testing >= 30 ===")
	results4, err := jp.Query("$.users[?(@.active === true && @.age >= 30)]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Active users 30 or older: %d\n", len(results4))
	for i, r := range results4 {
		user := r.Value.(map[string]interface{})
		fmt.Printf("  [%d] %s\n", i, user["name"])
	}
}