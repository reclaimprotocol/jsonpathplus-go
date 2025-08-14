package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	// Load test data
	data, err := os.ReadFile("data/goessner_spec_data.json")
	if err != nil {
		log.Fatal(err)
	}

	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		log.Fatal(err)
	}

	// Test the specific query
	query := "$..*[?(@property === 'price' && @ !== 8.95)]"
	fmt.Printf("Testing query: %s\n", query)

	results, err := jsonpathplus.Query(query, jsonData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Go Results: %d\n", len(results))
	for i, result := range results {
		fmt.Printf("%d. Path: %s, Value: %v, Property: %s\n", 
			i+1, result.Path, result.Value, result.ParentProperty)
	}

	// Let's also test just the recursive descent without filter
	fmt.Printf("\n=== Testing $..*  ===\n")
	allResults, err := jsonpathplus.Query("$..*", jsonData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("All recursive results: %d\n", len(allResults))
	
	// Find all price properties
	fmt.Printf("\n=== Finding price properties manually ===\n")
	for i, result := range allResults {
		if result.ParentProperty == "price" {
			fmt.Printf("%d. Path: %s, Value: %v, Property: %s\n", 
				i+1, result.Path, result.Value, result.ParentProperty)
		}
	}
}