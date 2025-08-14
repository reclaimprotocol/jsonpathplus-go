package types

import (
	"fmt"
	"strconv"
	"strings"
)

// Result represents a single result from a JSONPath query
type Result struct {
	Value          interface{} // The actual value found
	Path           string      // JSONPath to this value
	Parent         interface{} // Parent object/array containing this value
	ParentProperty string      // Property name or array index in parent
	Index          int         // Index in the result set
	OriginalIndex  int         // Original index in parent array (for array elements)
	// String index fields for character position tracking
	Start  int // Starting character position in original JSON string
	End    int // Ending character position in original JSON string
	Length int // Length of the element in the original JSON string
}

// String returns a string representation of the result
func (r Result) String() string {
	return fmt.Sprintf("Result{Value: %v, Path: %s}", r.Value, r.Path)
}

// Options configures JSONPath query execution
type Options struct {
	Root interface{} // Root object for $ references in filters
}

// AstNode represents a node in the Abstract Syntax Tree for JSONPath expressions
type AstNode struct {
	Type     string     // Node type: "root", "property", "wildcard", "index", "slice", "filter", "recursive", "union"
	Value    string     // Node value (property name, index, filter expression, etc.)
	Children []*AstNode // Child nodes
}

// String returns a string representation of the AST node
func (n *AstNode) String() string {
	if len(n.Children) == 0 {
		return fmt.Sprintf("%s(%s)", n.Type, n.Value)
	}
	return fmt.Sprintf("%s(%s)[%d children]", n.Type, n.Value, len(n.Children))
}

// Context holds evaluation context for advanced JSONPath features
type Context struct {
	Root           interface{} // Root object
	Current        interface{} // Current object being evaluated
	Parent         interface{} // Parent of current object
	ParentProperty string      // Property name or index in parent
	Path           string      // Current JSONPath
	Index          int         // Current index (for arrays)
}

// NewContext creates a new evaluation context
func NewContext(root, current, parent interface{}, parentProperty, path string, index int) *Context {
	return &Context{
		Root:           root,
		Current:        current,
		Parent:         parent,
		ParentProperty: parentProperty,
		Path:           path,
		Index:          index,
	}
}

// GetBracketPath returns the path in bracket notation format
func (ctx *Context) GetBracketPath() string {
	return convertToBracketNotation(ctx.Path)
}

// convertToBracketNotation converts dot notation paths to bracket notation
// e.g. "$.store.book[0]" -> "$['store']['book'][0]"
func convertToBracketNotation(path string) string {
	if path == "$" {
		return "$"
	}

	result := "$"
	remaining := path[1:] // Remove initial $

	for len(remaining) > 0 {
		if remaining[0] == '.' {
			// Handle property access
			remaining = remaining[1:] // Skip the dot

			// Find the end of the property name
			end := 0
			for end < len(remaining) && remaining[end] != '.' && remaining[end] != '[' {
				end++
			}

			if end > 0 {
				property := remaining[:end]
				result += fmt.Sprintf("['%s']", property)
				remaining = remaining[end:]
			}
		} else if remaining[0] == '[' {
			// Handle array access - find the matching closing bracket
			bracketEnd := strings.Index(remaining, "]")
			if bracketEnd != -1 {
				bracket := remaining[:bracketEnd+1]
				result += bracket
				remaining = remaining[bracketEnd+1:]
			} else {
				break
			}
		} else {
			break
		}
	}

	return result
}

// GetPropertyName returns the property name for the current context
func (c *Context) GetPropertyName() string {
	return c.ParentProperty
}

// GetParent returns the parent object
func (c *Context) GetParent() interface{} {
	return c.Parent
}

// IsArrayIndex returns true if the current property is an array index
func (c *Context) IsArrayIndex() bool {
	if c.ParentProperty == "" {
		return false
	}
	_, err := strconv.Atoi(c.ParentProperty)
	return err == nil
}
