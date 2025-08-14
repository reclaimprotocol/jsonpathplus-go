package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/reclaimprotocol/jsonpathplus-go/pkg/types"
)

// Parser handles JSONPath expression parsing
type Parser struct{}

// NewParser creates a new parser instance
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses a JSONPath expression into an AST
func (p *Parser) Parse(path string) (*types.AstNode, error) {
	if path == "" {
		return nil, fmt.Errorf("empty path")
	}

	// Remove leading/trailing whitespace
	path = strings.TrimSpace(path)

	// Handle root-only path
	if path == "$" {
		return &types.AstNode{Type: "root", Value: "$"}, nil
	}

	// Parse the path starting from root
	if !strings.HasPrefix(path, "$") {
		return nil, fmt.Errorf("path must start with $")
	}

	return p.parseFromRoot(path[1:]) // Remove the $ prefix
}

// parseFromRoot parses the path after the root $
func (p *Parser) parseFromRoot(path string) (*types.AstNode, error) {
	root := &types.AstNode{Type: "root", Value: "$"}

	if path == "" {
		return root, nil
	}

	// Parse the rest of the path
	current := root
	remaining := path

	for remaining != "" {
		node, nextPos, err := p.parseNextSegment(remaining)
		if err != nil {
			return nil, err
		}

		current.Children = append(current.Children, node)
		current = node
		remaining = remaining[nextPos:]
	}

	return root, nil
}

// parseNextSegment parses the next segment of the path
func (p *Parser) parseNextSegment(path string) (*types.AstNode, int, error) {
	if path == "" {
		return nil, 0, fmt.Errorf("empty segment")
	}

	switch path[0] {
	case '.':
		return p.parseDotSegment(path)
	case '[':
		return p.parseBracketSegment(path)
	default:
		return nil, 0, fmt.Errorf("unexpected character: %c at position 0", path[0])
	}
}

// parseDotSegment parses segments starting with '.'
func (p *Parser) parseDotSegment(path string) (*types.AstNode, int, error) {
	if len(path) < 2 {
		return nil, 0, fmt.Errorf("incomplete dot segment")
	}

	// Handle recursive descent (..)
	if path[1] == '.' {
		recursiveNode := &types.AstNode{Type: "recursive", Value: ".."}
		// If next is a wildcard, attach it as a child and consume
		if len(path) > 2 && path[2] == '*' {
			wildcardNode := &types.AstNode{Type: "wildcard", Value: "*"}
			recursiveNode.Children = append(recursiveNode.Children, wildcardNode)
			return recursiveNode, 3, nil
		}

		// Check if there's more after the ..
		if len(path) > 2 {
			// If the next character is not . or [, it's a direct property after ..
			if path[2] != '.' && path[2] != '[' {
				// Find the end of the property name
				end := p.findPropertyEnd(path, 2)
				if end > 2 {
					property := path[2:end]
					propertyNode := &types.AstNode{Type: "property", Value: property}

					// Check if there are more segments after the property (like [*])
					if end < len(path) {
						// Parse the remaining path and attach as children to the property
						remaining := path[end:]
						currentNode := propertyNode

						for remaining != "" {
							nextNode, nextPos, err := p.parseNextSegment(remaining)
							if err != nil {
								break
							}
							currentNode.Children = append(currentNode.Children, nextNode)
							currentNode = nextNode
							remaining = remaining[nextPos:]
						}
					}

					recursiveNode.Children = append(recursiveNode.Children, propertyNode)
					return recursiveNode, len(path), nil
				}
			}
		}

		return recursiveNode, 2, nil
	}

	// Handle property access (.property) or wildcard (.*)
	pos := 1
	if path[1] == '*' {
		// Check for operators after wildcard
		if len(path) > 2 {
			if path[2] == '~' {
				return &types.AstNode{Type: "property_names", Value: "*"}, 3, nil
			}
			if path[2] == '^' {
				// Return parent of all wildcard children by first selecting wildcard, then parent
				wildcardNode := &types.AstNode{Type: "wildcard", Value: "*"}
				parentNode := &types.AstNode{Type: "parent", Value: "*"}
				parentNode.Children = []*types.AstNode{wildcardNode}
				return parentNode, 3, nil
			}
		}
		return &types.AstNode{Type: "wildcard", Value: "*"}, 2, nil
	}

	// Extract property name
	end := p.findPropertyEnd(path, pos)
	if end == pos {
		return nil, 0, fmt.Errorf("empty property name")
	}

	property := path[pos:end]

	// Check for operators after property
	if end < len(path) {
		if path[end] == '~' {
			return &types.AstNode{Type: "property_names", Value: property}, end + 1, nil
		}
		if path[end] == '^' {
			// Build a parent node whose child selects the property first
			propNode := &types.AstNode{Type: "property", Value: property}
			parentNode := &types.AstNode{Type: "parent", Value: property}
			parentNode.Children = []*types.AstNode{propNode}
			return parentNode, end + 1, nil
		}
	}

	return &types.AstNode{Type: "property", Value: property}, end, nil
}

// parseBracketSegment parses segments starting with '['
func (p *Parser) parseBracketSegment(path string) (*types.AstNode, int, error) {
	end := p.findMatchingBracket(path, 0)
	if end == -1 {
		return nil, 0, fmt.Errorf("unmatched bracket")
	}

	content := path[1:end]

	// Check for chained operations (multiple bracket segments)
	nextPos := end + 1
	var chainedNodes []*types.AstNode

	// Parse the main bracket content
	node, err := p.parseBracketContent(content)
	if err != nil {
		return nil, 0, err
	}
	chainedNodes = append(chainedNodes, node)

	// Look for additional bracket segments
	remaining := path[nextPos:]
	for strings.HasPrefix(remaining, "[") {
		chainEnd := p.findMatchingBracket(remaining, 0)
		if chainEnd == -1 {
			break
		}

		chainContent := remaining[1:chainEnd]
		chainNode, err := p.parseBracketContent(chainContent)
		if err != nil {
			break
		}

		chainedNodes = append(chainedNodes, chainNode)
		nextPos += chainEnd + 1
		remaining = path[nextPos:]
	}

	// Check for operators after bracket segments
	finalNode := node
	if len(chainedNodes) > 1 {
		chainNode := &types.AstNode{
			Type:  "chain",
			Value: "chained_operations",
		}
		chainNode.Children = chainedNodes
		finalNode = chainNode
	}

	// Check for operators after the bracket expression
	remaining = path[nextPos:]
	if len(remaining) > 0 {
		if remaining[0] == '~' {
			operatorNode := &types.AstNode{Type: "property_names", Value: ""}
			operatorNode.Children = []*types.AstNode{finalNode}
			return operatorNode, nextPos + 1, nil
		}
		if remaining[0] == '^' {
			operatorNode := &types.AstNode{Type: "parent", Value: ""}
			operatorNode.Children = []*types.AstNode{finalNode}
			return operatorNode, nextPos + 1, nil
		}
	}

	return finalNode, nextPos, nil
}

// parseBracketContent parses the content inside brackets
func (p *Parser) parseBracketContent(content string) (*types.AstNode, error) {
	content = strings.TrimSpace(content)

	if content == "" {
		return nil, fmt.Errorf("empty bracket content")
	}

	// Handle wildcard
	if content == "*" {
		return &types.AstNode{Type: "index_wildcard", Value: "*"}, nil
	}

	// Handle filter expressions
	if strings.HasPrefix(content, "?") {
		return &types.AstNode{Type: "filter", Value: content}, nil
	}

	// Handle quoted property names
	if (strings.HasPrefix(content, "'") && strings.HasSuffix(content, "'")) ||
		(strings.HasPrefix(content, "\"") && strings.HasSuffix(content, "\"")) {
		property := content[1 : len(content)-1]

		// Check for special operators
		if strings.HasSuffix(property, "~") {
			return &types.AstNode{Type: "property_names", Value: strings.TrimSuffix(property, "~")}, nil
		}
		if strings.HasSuffix(property, "^") {
			return &types.AstNode{Type: "parent", Value: strings.TrimSuffix(property, "^")}, nil
		}

		return &types.AstNode{Type: "property", Value: property}, nil
	}

	// Handle property names operator
	if strings.HasSuffix(content, "~") {
		return &types.AstNode{Type: "property_names", Value: strings.TrimSuffix(content, "~")}, nil
	}

	// Handle parent operator
	if strings.HasSuffix(content, "^") {
		return &types.AstNode{Type: "parent", Value: strings.TrimSuffix(content, "^")}, nil
	}

	// Handle union (comma-separated values)
	if strings.Contains(content, ",") {
		parts := p.splitUnion(content)
		unionNode := &types.AstNode{Type: "union", Value: "union"}

		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}

			// Parse each union part
			childNode, err := p.parseUnionPart(part)
			if err != nil {
				return nil, err
			}
			unionNode.Children = append(unionNode.Children, childNode)
		}

		return unionNode, nil
	}

	// Handle slice notation
	if strings.Contains(content, ":") {
		return &types.AstNode{Type: "slice", Value: content}, nil
	}

	// Handle array index
	if idx, err := strconv.Atoi(content); err == nil {
		return &types.AstNode{Type: "index", Value: strconv.Itoa(idx)}, nil
	}

	// Handle negative array index
	if strings.HasPrefix(content, "-") {
		if idx, err := strconv.Atoi(content); err == nil {
			return &types.AstNode{Type: "index", Value: strconv.Itoa(idx)}, nil
		}
	}

	// Default to property
	return &types.AstNode{Type: "property", Value: content}, nil
}

// parseUnionPart parses a single part of a union expression
func (p *Parser) parseUnionPart(part string) (*types.AstNode, error) {
	// Handle quoted strings
	if (strings.HasPrefix(part, "'") && strings.HasSuffix(part, "'")) ||
		(strings.HasPrefix(part, "\"") && strings.HasSuffix(part, "\"")) {
		property := part[1 : len(part)-1]
		return &types.AstNode{Type: "property", Value: property}, nil
	}

	// Handle array indices
	if idx, err := strconv.Atoi(part); err == nil {
		return &types.AstNode{Type: "index", Value: strconv.Itoa(idx)}, nil
	}

	// Handle negative indices
	if strings.HasPrefix(part, "-") {
		if idx, err := strconv.Atoi(part); err == nil {
			return &types.AstNode{Type: "index", Value: strconv.Itoa(idx)}, nil
		}
	}

	// Default to property
	return &types.AstNode{Type: "property", Value: part}, nil
}

// Helper functions

func (p *Parser) findPropertyEnd(path string, start int) int {
	for i := start; i < len(path); i++ {
		ch := path[i]
		if ch == '.' || ch == '[' || ch == '~' || ch == '^' {
			return i
		}
	}
	return len(path)
}

func (p *Parser) findMatchingBracket(path string, start int) int {
	if start >= len(path) || path[start] != '[' {
		return -1
	}

	depth := 0
	inQuotes := false
	quoteChar := byte(0)

	for i := start; i < len(path); i++ {
		ch := path[i]

		if !inQuotes {
			if ch == '\'' || ch == '"' {
				inQuotes = true
				quoteChar = ch
			} else if ch == '[' {
				depth++
			} else if ch == ']' {
				depth--
				if depth == 0 {
					return i
				}
			}
		} else {
			if ch == quoteChar && (i == 0 || path[i-1] != '\\') {
				inQuotes = false
				quoteChar = 0
			}
		}
	}

	return -1
}

func (p *Parser) splitUnion(content string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false
	quoteChar := byte(0)
	depth := 0

	for i := 0; i < len(content); i++ {
		ch := content[i]

		if !inQuotes {
			if ch == '\'' || ch == '"' {
				inQuotes = true
				quoteChar = ch
				current.WriteByte(ch)
			} else if ch == '[' {
				depth++
				current.WriteByte(ch)
			} else if ch == ']' {
				depth--
				current.WriteByte(ch)
			} else if ch == ',' && depth == 0 {
				parts = append(parts, current.String())
				current.Reset()
			} else {
				current.WriteByte(ch)
			}
		} else {
			current.WriteByte(ch)
			if ch == quoteChar && (i == 0 || content[i-1] != '\\') {
				inQuotes = false
				quoteChar = 0
			}
		}
	}

	parts = append(parts, current.String())
	return parts
}

// ValidatePath validates a JSONPath expression syntax
func (p *Parser) ValidatePath(path string) error {
	_, err := p.Parse(path)
	return err
}
