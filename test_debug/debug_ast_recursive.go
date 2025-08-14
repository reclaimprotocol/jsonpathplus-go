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
	fmt.Println("=== AST Analysis for $..*===")
	
	ast, err := jp.Parse("$..*")
	if err != nil {
		fmt.Printf("❌ Parse error: %v\n", err)
		return
	}
	
	printAST(ast, "")
}