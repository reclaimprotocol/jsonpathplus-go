package main

import (
	"fmt"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
	"github.com/reclaimprotocol/jsonpathplus-go/pkg/types"
)

func printAST(node *types.AstNode, indent string) {
	if node == nil {
		return
	}
	fmt.Printf("%sType: %s, Value: %s\n", indent, node.Type, node.Value)
	for i, child := range node.Children {
		fmt.Printf("%sChild %d:\n", indent, i)
		printAST(child, indent+"  ")
	}
}

func main() {
	// Test different wildcard expressions to see how they parse
	paths := []string{
		"$.*",   // Property wildcard
		"$[*]",  // Index wildcard  
	}
	
	for _, path := range paths {
		fmt.Printf("=== AST for: %s ===\n", path)
		ast, err := jp.Parse(path)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		printAST(ast, "")
		fmt.Println()
		
		// Test on an array
		jsonStr := `[{"id":1,"name":"first"},{"id":2,"name":"second"}]`
		fmt.Printf("Testing %s on array:\n", path)
		results, err := jp.Query(path, jsonStr)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Results: %d\n", len(results))
			for i, result := range results {
				fmt.Printf("  [%d] Value: %v, Path: %s\n", i, result.Value, result.Path)
			}
		}
		fmt.Println()
	}
}