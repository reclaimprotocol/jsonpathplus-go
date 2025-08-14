package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	jsonData := `{"store": {"book": [{"title": "Book 1"}, {"title": "Book 2"}]}}`

	results, err := jp.Query("$..book", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("$..book found %d results:\n", len(results))
	for i, r := range results {
		fmt.Printf("  [%d] %v (path: %s)\n", i, r.Value, r.Path)
	}

	fmt.Println()

	results2, err := jp.Query("$..book[*]", jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("$..book[*] found %d results:\n", len(results2))
	for i, r := range results2 {
		fmt.Printf("  [%d] %v (path: %s)\n", i, r.Value, r.Path)
	}
}
