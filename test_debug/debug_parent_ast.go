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
	fmt.Printf("%s%s: %s (children: %d)\n", indent, node.Type, node.Value, len(node.Children))
	for i, child := range node.Children {
		fmt.Printf("%s  child[%d]:\n", indent, i)
		printAST(child, indent+"    ")
	}
}

func main() {
	fmt.Println("=== AST Analysis for parent filter ===")
	
	// Test 1: Simple parent filter
	fmt.Println("\n1. $.store.book[?(@parent.bicycle)]")
	ast1, err := jp.Parse("$.store.book[?(@parent.bicycle)]")
	if err != nil {
		fmt.Printf("❌ Parse error: %v\n", err)
		return
	}
	printAST(ast1, "")
	
	// Test 2: Index wildcard for comparison
	fmt.Println("\n2. $.store.book[*] for comparison")
	ast2, err := jp.Parse("$.store.book[*]")
	if err != nil {
		fmt.Printf("❌ Parse error: %v\n", err)
		return
	}
	printAST(ast2, "")
}