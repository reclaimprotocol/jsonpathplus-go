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
	Root                   interface{} // Root object
	Current                interface{} // Current object being evaluated
	Parent                 interface{} // Parent of current object
	ParentProperty         string      // Property name or index in parent
	Path                   string      // Current JSONPath
	Index                  int         // Current index (for arrays)
	ParentOfParentProperty string      // Property that led to the parent (for @parentProperty)
	ActualParentArray      interface{} // For array elements, the actual array (for @property type detection)
}

// NewContext creates a new evaluation context
func NewContext(root, current, parent interface{}, parentProperty, path string, index int) *Context {
	return &Context{
		Root:                   root,
		Current:                current,
		Parent:                 parent,
		ParentProperty:         parentProperty,
		Path:                   path,
		Index:                  index,
		ParentOfParentProperty: "",  // Default empty
		ActualParentArray:      nil, // Default nil
	}
}

// NewArrayElementContext creates a context for array elements with proper parent tracking
func NewArrayElementContext(root, current, parent interface{}, parentProperty, path string, index int, actualArray interface{}) *Context {
	return &Context{
		Root:                   root,
		Current:                current,
		Parent:                 parent,
		ParentProperty:         parentProperty,
		Path:                   path,
		Index:                  index,
		ParentOfParentProperty: "",          // Default empty
		ActualParentArray:      actualArray, // The array that contains this element
	}
}

// NewContextWithParentProperty creates a new evaluation context with parent property tracking
func NewContextWithParentProperty(root, current, parent interface{}, parentProperty, path string, index int, parentOfParentProperty string) *Context {
	return &Context{
		Root:                   root,
		Current:                current,
		Parent:                 parent,
		ParentProperty:         parentProperty,
		Path:                   path,
		Index:                  index,
		ParentOfParentProperty: parentOfParentProperty,
		ActualParentArray:      nil, // Default nil
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

// GetParentPropertyName returns the property name that led to the parent (for @parentProperty)
func (c *Context) GetParentPropertyName() string {
	// If explicitly set, use it
	if c.ParentOfParentProperty != "" {
		return c.ParentOfParentProperty
	}

	// Otherwise, derive from path
	// For path like "$.users.1['name']", @parentProperty should be "users"
	// For path like "$.store.book[0]['title']", @parentProperty should be "book"
	return extractParentPropertyFromPath(c.Path)
}

// extractParentPropertyFromPath extracts the parent property name from a JSONPath
func extractParentPropertyFromPath(path string) string {
	// For @parentProperty, we want the property that led to the parent of the current item
	// Examples:
	// "$.users.1['name']" -> "users" (parent of parent of 'name')
	// "$.store.book[0]['title']" -> "store" (parent of parent of 'title')
	// "$.users['1']" -> "" (filtering users object directly)

	if strings.Contains(path, "['") {
		// Path has bracket notation for current property like $.users.1['name']
		// Remove the bracket part to get the parent path
		lastBracket := strings.LastIndex(path, "['")
		if lastBracket > 0 {
			parentPath := path[:lastBracket] // "$.users.1"

			// Now get the parent of this parent
			// Find the second-to-last property
			if strings.Contains(parentPath, ".") {
				parts := strings.Split(parentPath, ".")
				if len(parts) >= 3 { // $, intermediate, property
					// Get the second-to-last part (parent of parent)
					return parts[len(parts)-2]
				}
			}
		}
	}

	return "" // Default for cases we can't parse
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

// IsParentArray returns true if the parent is an array (not an object)
func (c *Context) IsParentArray() bool {
	// For array elements, we track the actual array separately
	if c.ActualParentArray != nil {
		_, isArray := c.ActualParentArray.([]interface{})
		return isArray
	}

	// Otherwise check the Parent field
	if c.Parent == nil {
		return false
	}
	_, isArray := c.Parent.([]interface{})
	return isArray
}

// GetPropertyValue returns the property as the appropriate type (number for array indices, string for object keys)
func (c *Context) GetPropertyValue() interface{} {
	if c.IsParentArray() && c.IsArrayIndex() {
		// For array indices, return as number
		if idx, err := strconv.Atoi(c.ParentProperty); err == nil {
			return idx
		}
	}
	// For object keys, return as string
	return c.ParentProperty
}
