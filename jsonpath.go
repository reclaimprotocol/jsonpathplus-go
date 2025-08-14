package jsonpathplus

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/reclaimprotocol/jsonpathplus-go/internal/evaluator"
	"github.com/reclaimprotocol/jsonpathplus-go/internal/parser"
	"github.com/reclaimprotocol/jsonpathplus-go/pkg/types"
	"github.com/reclaimprotocol/jsonpathplus-go/pkg/utils"
)

// Result represents a JSONPath query result (alias for types.Result for backward compatibility)
type Result = types.Result

// Options represents JSONPath options (alias for types.Options for backward compatibility)
type Options = types.Options

// JSONPathEngine is the main engine for JSONPath operations
type JSONPathEngine struct {
	parser    *parser.Parser
	evaluator *evaluator.Evaluator
}

// NewJSONPathEngine creates a new JSONPath engine
func NewJSONPathEngine() *JSONPathEngine {
	return &JSONPathEngine{
		parser:    parser.NewParser(),
		evaluator: evaluator.NewEvaluator(),
	}
}

// JSONPath represents a compiled JSONPath expression using the new architecture
type JSONPath struct {
	path   string
	ast    *types.AstNode
	engine *JSONPathEngine
}

// New creates a new JSONPath instance
func New(path string) (*JSONPath, error) {
	engine := NewJSONPathEngine()

	ast, err := engine.parser.Parse(path)
	if err != nil {
		return nil, err
	}

	return &JSONPath{
		path:   path,
		ast:    ast,
		engine: engine,
	}, nil
}

// Execute executes the JSONPath against the given data
func (jp *JSONPath) Execute(data interface{}) (results []Result, err error) {
	options := &types.Options{}
	
	// Catch panics from JavaScript compatibility errors (like null.length)
	defer func() {
		if r := recover(); r != nil {
			if str, ok := r.(string); ok && strings.Contains(str, "Cannot read properties of null") {
				// JavaScript compatibility: propagate the error instead of returning empty results
				results = []Result{}
				err = fmt.Errorf("jsonPath: %s", str)
				return
			}
			// Re-panic for other errors
			panic(r)
		}
	}()
	
	results = jp.engine.evaluator.Evaluate(jp.ast, data, options)
	err = nil
	return
}

// ExecuteWithOptions executes the JSONPath with custom options
func (jp *JSONPath) ExecuteWithOptions(data interface{}, options *Options) ([]Result, error) {
	if options == nil {
		options = &Options{}
	}
	return jp.engine.evaluator.Evaluate(jp.ast, data, options), nil
}

// Path returns the JSONPath expression
func (jp *JSONPath) Path() string {
	return jp.path
}

// AST returns the parsed AST
func (jp *JSONPath) AST() *types.AstNode {
	return jp.ast
}

// Convenience functions for the refactored API

// JSONParse parses a JSON string into an interface{} with preserved property order
func JSONParse(jsonStr string) (interface{}, error) {
	return utils.ParseOrderedJSON([]byte(jsonStr))
}

// Query executes a JSONPath query against JSON string or data
func Query(path string, input interface{}) ([]Result, error) {
	var data interface{}
	var jsonStr string
	var isStringInput bool

	// Check if input is a string (JSON) or already parsed data
	if str, ok := input.(string); ok {
		jsonStr = str
		isStringInput = true
		var err error
		data, err = utils.ParseOrderedJSON([]byte(jsonStr))
		if err != nil {
			return nil, err
		}
	} else {
		data = input
		isStringInput = false
	}

	jp, err := New(path)
	if err != nil {
		return nil, err
	}

	results, err := jp.Execute(data)
	if err != nil {
		return nil, err
	}

	// If input was a JSON string, calculate string indices for each result
	if isStringInput {
		for i := range results {
			stringPos := findStringPositionForResult(results[i], jsonStr)
			results[i].Start = stringPos.Start
			results[i].End = stringPos.End
			results[i].Length = stringPos.Length
			// For backward compatibility, also set OriginalIndex to the start position
			results[i].OriginalIndex = stringPos.Start
		}
	}

	return results, nil
}

// Parse parses a JSONPath expression and returns the AST
func Parse(path string) (*types.AstNode, error) {
	p := parser.NewParser()
	return p.Parse(path)
}

// Validate validates a JSONPath expression
func Validate(path string) error {
	p := parser.NewParser()
	return p.ValidatePath(path)
}

// Additional JSONPathEngine methods for backward compatibility

// Query executes a JSONPath query using the engine
func (engine *JSONPathEngine) Query(path string, input interface{}) ([]Result, error) {
	return Query(path, input)
}

// Close closes the engine (no-op for compatibility)
func (engine *JSONPathEngine) Close() error {
	return nil
}

// NewEngine creates a new JSONPath engine (backward compatibility)
func NewEngine() *JSONPathEngine {
	return NewJSONPathEngine()
}

// String position tracking functionality

// StringPosition represents a position in the original JSON string
type StringPosition struct {
	Start  int // Starting character position
	End    int // Ending character position
	Length int // Length of the element
}

// findStringPositionForResult finds the character position of a result in the JSON string
func findStringPositionForResult(result Result, jsonStr string) StringPosition {
	// Handle different types of JSONPath results
	path := result.Path

	// Root element
	if path == "$" {
		return StringPosition{Start: 0, End: len(jsonStr), Length: len(jsonStr)}
	}

	// Handle property access in objects
	if result.ParentProperty != "" && result.Parent != nil {
		// Check if parent is an array (array element access)
		if _, isArray := result.Parent.([]interface{}); isArray {
			if idx, err := strconv.Atoi(result.ParentProperty); err == nil {
				return findArrayElementPosition(idx, jsonStr)
			}
		}

		// Property access in objects
		if _, isObject := result.Parent.(map[string]interface{}); isObject {
			return findPropertyValuePosition(result.ParentProperty, jsonStr, path)
		}
	}

	// Fallback: try to find by parsing the path
	return findPositionByPath(path, jsonStr)
}

// findArrayElementPosition finds the position of an array element by index
func findArrayElementPosition(index int, jsonStr string) StringPosition {
	// Find the first array in the JSON
	arrayStart := strings.Index(jsonStr, "[")
	if arrayStart == -1 {
		return StringPosition{}
	}

	pos := arrayStart + 1
	currentIndex := 0

	// Skip whitespace
	for pos < len(jsonStr) && isWhitespace(jsonStr[pos]) {
		pos++
	}

	// If we want index 0, we're at the first element
	if index == 0 {
		elementEnd := skipJSONElement(jsonStr, pos)
		return StringPosition{Start: pos, End: elementEnd, Length: elementEnd - pos}
	}

	// Skip elements until we reach the target index
	for currentIndex < index && pos < len(jsonStr) {
		pos = skipJSONElement(jsonStr, pos)

		// Skip whitespace and comma
		for pos < len(jsonStr) && (isWhitespace(jsonStr[pos]) || jsonStr[pos] == ',') {
			pos++
		}
		currentIndex++
	}

	if pos < len(jsonStr) {
		elementEnd := skipJSONElement(jsonStr, pos)
		return StringPosition{Start: pos, End: elementEnd, Length: elementEnd - pos}
	}

	return StringPosition{}
}

// findPropertyValuePosition finds the position of a property key (not value)
func findPropertyValuePosition(propName, jsonStr, path string) StringPosition {
	// Extract the array index from the path if present
	arrayIndex := -1
	if strings.Contains(path, "[") && strings.Contains(path, "]") {
		// Extract array index from path like $[0].id or $.arr[1].prop
		start := strings.LastIndex(path, "[")
		end := strings.Index(path[start:], "]")
		if start != -1 && end != -1 {
			indexStr := path[start+1 : start+end]
			if idx, err := strconv.Atoi(indexStr); err == nil {
				arrayIndex = idx
			}
		}
	}

	// Search for the property key
	searchStr := fmt.Sprintf(`"%s"`, propName)

	pos := 0
	occurrenceCount := 0
	for {
		idx := strings.Index(jsonStr[pos:], searchStr)
		if idx == -1 {
			break
		}

		absolutePos := pos + idx

		// Check if this is a property key (followed by colon)
		afterKey := absolutePos + len(searchStr)
		for afterKey < len(jsonStr) && isWhitespace(jsonStr[afterKey]) {
			afterKey++
		}

		if afterKey < len(jsonStr) && jsonStr[afterKey] == ':' {
			// If we have an array index, return the occurrence that matches
			if arrayIndex == -1 || occurrenceCount == arrayIndex {
				return StringPosition{Start: absolutePos, End: absolutePos + len(searchStr), Length: len(searchStr)}
			}
			occurrenceCount++
		}

		pos = absolutePos + 1
	}

	return StringPosition{}
}

// findPositionByPath attempts to find position by parsing the JSONPath
func findPositionByPath(path, jsonStr string) StringPosition {
	// This is a simplified implementation
	// For complex paths, we'd need more sophisticated parsing

	// Handle simple property access like $.property
	if strings.HasPrefix(path, "$.") && !strings.Contains(path[2:], ".") && !strings.Contains(path, "[") {
		propName := path[2:]
		return findPropertyValuePosition(propName, jsonStr, path)
	}

	// Handle array access like $[0], $[1], etc.
	if strings.HasPrefix(path, "$[") && strings.HasSuffix(path, "]") {
		indexStr := path[2 : len(path)-1]
		if idx, err := strconv.Atoi(indexStr); err == nil {
			return findArrayElementPosition(idx, jsonStr)
		}
	}

	// For complex paths, return empty position
	return StringPosition{}
}

// Helper functions for JSON parsing

// skipJSONElement skips over a complete JSON element and returns the position after it
func skipJSONElement(jsonStr string, start int) int {
	if start >= len(jsonStr) {
		return start
	}

	char := jsonStr[start]
	switch char {
	case '"':
		// Skip string
		pos := start + 1
		for pos < len(jsonStr) {
			if jsonStr[pos] == '"' && (pos == start+1 || jsonStr[pos-1] != '\\') {
				return pos + 1
			}
			pos++
		}
		return pos
	case '{':
		// Skip object
		return skipBracedElement(jsonStr, start, '{', '}')
	case '[':
		// Skip array
		return skipBracedElement(jsonStr, start, '[', ']')
	default:
		// Skip primitive (number, boolean, null)
		pos := start
		for pos < len(jsonStr) && !isWhitespace(jsonStr[pos]) &&
			jsonStr[pos] != ',' && jsonStr[pos] != '}' && jsonStr[pos] != ']' {
			pos++
		}
		return pos
	}
}

// skipBracedElement skips over a complete braced element (object or array)
func skipBracedElement(jsonStr string, start int, openBrace, closeBrace byte) int {
	pos := start + 1
	depth := 1
	inString := false

	for pos < len(jsonStr) && depth > 0 {
		char := jsonStr[pos]

		if char == '"' && (pos == 0 || jsonStr[pos-1] != '\\') {
			inString = !inString
		} else if !inString {
			if char == openBrace {
				depth++
			} else if char == closeBrace {
				depth--
			}
		}
		pos++
	}

	return pos
}

// isWhitespace checks if a character is whitespace
func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}
