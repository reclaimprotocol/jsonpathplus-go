package main

import (
	"fmt"
	"github.com/reclaimprotocol/jsonpathplus-go/pkg/types"
)

func main() {
	// Test the bracket notation conversion directly
	testPaths := []string{
		"$",
		"$.store",
		"$.store.book",
		"$.store.book[0]",
		"$.store.book[0].title",
		"$.users[1].profile.bio",
	}

	fmt.Println("=== Testing bracket notation conversion ===")
	for _, path := range testPaths {
		// Create a dummy context to test the conversion
		ctx := types.NewContext(nil, nil, nil, "", path, 0)
		bracketPath := ctx.GetBracketPath()
		fmt.Printf("Original: %s\n", path)
		fmt.Printf("Bracket:  %s\n", bracketPath)
		fmt.Println()
	}
}
