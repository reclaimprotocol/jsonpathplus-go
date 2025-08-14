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
	paths := []string{
		"$[*]",
		"$[*].id",
		"$.id",
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
	}
}