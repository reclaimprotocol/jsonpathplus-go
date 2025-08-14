package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
	// Test the complex JSONPath expression from the original user request
	jsonStr := `{
  "store": {
    "book": [
      {
        "category": "reference",
        "author": "Nigel Rees",
        "title": "Sayings of the Century",
        "price": 8.95
      },
      {
        "category": "fiction",
        "author": "Evelyn Waugh", 
        "title": "Sword of Honour",
        "price": 12.99
      },
      {
        "category": "FICTION",
        "author": "Herman Melville",
        "title": "Moby Dick",
        "isbn": "0-553-21311-3",
        "price": 8.99
      },
      {
        "category": "ACTION",
        "author": "J. R. R. Tolkien",
        "title": "The Lord of the Rings",
        "isbn": "0-395-19395-8",
        "price": 22.99
      }
    ]
  }
}`

	fmt.Println("=== Testing Complex JSONPath Expression ===")
	fmt.Printf("Expression: $..book.*[?(@property === \"category\" && @.match(/TION$/i))]\n")
	fmt.Printf("Description: All categories of books which end in 'TION' (case insensitive)\n\n")

	results, err := jp.Query(`$..book.*[?(@property === "category" && @.match(/TION$/i))]`, jsonStr)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("Results: %d\n", len(results))
	for i, result := range results {
		fmt.Printf("  [%d] Value: %v\n", i, result.Value)
		fmt.Printf("      Path: %s\n", result.Path)
		fmt.Printf("      String Position: Start=%d, End=%d, Length=%d\n", 
			result.Start, result.End, result.Length)
		
		if result.Start > 0 && result.End > result.Start && result.End <= len(jsonStr) {
			extracted := jsonStr[result.Start:result.End]
			fmt.Printf("      Extracted: %q\n", extracted)
		}
		fmt.Println()
	}

	fmt.Println("Expected: Should match 'FICTION' and 'ACTION' categories")
}