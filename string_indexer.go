// Package jsonpathplus provides JSONPath querying with string character position tracking.
//
// This file contains the core string position tracking functionality that maps JSON elements
// to their exact character positions in the original JSON string.
package jsonpathplus

import (
	"encoding/json"
	"fmt"
	"strings"
)

// StringPosition represents a position in the original JSON string
type StringPosition struct {
	Start  int // Starting character position
	End    int // Ending character position
	Length int // Length of the element
}

// IndexedValue wraps a value with its string position information
type IndexedValue struct {
	Value    interface{}
	Position StringPosition
}

// StringIndexTracker tracks character positions during JSON parsing
type StringIndexTracker struct {
	originalJSON string
	pathToPos    map[string]StringPosition
}

// NewStringIndexTracker creates a new string index tracker
func NewStringIndexTracker(jsonStr string) *StringIndexTracker {
	return &StringIndexTracker{
		originalJSON: jsonStr,
		pathToPos:    make(map[string]StringPosition),
	}
}

// JSONParseWithIndex parses JSON while preserving string positions
func JSONParseWithIndex(jsonStr string) (*IndexedValue, error) {
	tracker := NewStringIndexTracker(jsonStr)

	// Parse the JSON normally first
	var value interface{}
	if err := json.Unmarshal([]byte(jsonStr), &value); err != nil {
		return nil, err
	}

	// Position tracking is handled in jsonpath_string_index.go
	// No need to build comprehensive position map here

	rootPos := StringPosition{Start: 0, End: len(jsonStr), Length: len(jsonStr)}
	tracker.pathToPos["$"] = rootPos

	return &IndexedValue{
		Value:    value,
		Position: rootPos,
	}, nil
}

// parseValue recursively parses values and tracks their positions
func (t *StringIndexTracker) parseValue(value interface{}, path string, startPos int) error {

	switch v := value.(type) {
	case map[string]interface{}:
		return t.parseObject(v, path, startPos)
	case []interface{}:
		return t.parseArray(v, path, startPos)
	default:
		// For primitive values, find their position in the string
		pos, err := t.findValuePosition(value, startPos)
		if err != nil {
			return err
		}
		t.pathToPos[path] = pos
		return nil
	}
}

// parseObject parses a JSON object and tracks property positions
func (t *StringIndexTracker) parseObject(obj map[string]interface{}, path string, startPos int) error {
	jsonStr := t.originalJSON

	// Find the opening brace
	objStart := t.findChar('{', startPos)
	if objStart == -1 {
		return fmt.Errorf("could not find opening brace for object")
	}

	pos := objStart + 1 // Start after the opening brace

	for key, val := range obj {
		// Find the property key in the JSON string
		keyPos := t.findPropertyKey(key, pos)
		if keyPos == -1 {
			// Fallback: search from current position
			keyPos = strings.Index(jsonStr[pos:], fmt.Sprintf(`"%s"`, key))
			if keyPos != -1 {
				keyPos += pos
			}
		}

		if keyPos != -1 {
			// Store position for this property
			propertyPath := path + "." + key

			// Find the value position (after the colon)
			colonPos := strings.Index(jsonStr[keyPos:], ":")
			if colonPos != -1 {
				valueStart := keyPos + colonPos + 1
				// Skip whitespace
				for valueStart < len(jsonStr) && isWhitespace(jsonStr[valueStart]) {
					valueStart++
				}

				valuePos := StringPosition{Start: keyPos, End: -1, Length: -1}
				t.pathToPos[propertyPath] = valuePos

				// Recursively parse the value
				err := t.parseValue(val, propertyPath, valueStart)
				if err != nil {
					return err
				}

				// Update position for next property
				pos = valueStart
			}
		}
	}

	return nil
}

// parseArray parses a JSON array and tracks element positions
func (t *StringIndexTracker) parseArray(arr []interface{}, path string, startPos int) error {
	jsonStr := t.originalJSON

	// Find the opening bracket
	arrStart := t.findChar('[', startPos)
	if arrStart == -1 {
		return fmt.Errorf("could not find opening bracket for array")
	}

	pos := arrStart + 1 // Start after the opening bracket

	for i, val := range arr {
		// Skip whitespace
		for pos < len(jsonStr) && isWhitespace(jsonStr[pos]) {
			pos++
		}

		elementPath := fmt.Sprintf("%s[%d]", path, i)

		// Store position for this array element
		elementPos := StringPosition{Start: pos, End: -1, Length: -1}
		t.pathToPos[elementPath] = elementPos

		// Recursively parse the element
		err := t.parseValue(val, elementPath, pos)
		if err != nil {
			return err
		}

		// Find next comma or end of array
		pos = t.findNextElementStart(pos)
	}

	return nil
}

// findPropertyKey finds the position of a property key in the JSON string
func (t *StringIndexTracker) findPropertyKey(key string, startPos int) int {
	jsonStr := t.originalJSON
	searchStr := fmt.Sprintf(`"%s"`, key)

	pos := startPos
	for {
		found := strings.Index(jsonStr[pos:], searchStr)
		if found == -1 {
			return -1
		}

		absolutePos := pos + found

		// Check if this is actually a property key (followed by colon)
		afterKey := absolutePos + len(searchStr)
		for afterKey < len(jsonStr) && isWhitespace(jsonStr[afterKey]) {
			afterKey++
		}

		if afterKey < len(jsonStr) && jsonStr[afterKey] == ':' {
			return absolutePos
		}

		// Not a property key, continue searching
		pos = absolutePos + 1
	}
}

// findValuePosition finds the character position of a primitive value
func (t *StringIndexTracker) findValuePosition(value interface{}, startPos int) (StringPosition, error) {
	jsonStr := t.originalJSON

	var searchStr string
	switch v := value.(type) {
	case string:
		searchStr = fmt.Sprintf(`"%s"`, v)
	case int:
		searchStr = fmt.Sprintf("%d", v)
	case float64:
		searchStr = fmt.Sprintf("%g", v)
	case bool:
		searchStr = fmt.Sprintf("%t", v)
	case nil:
		searchStr = "null"
	default:
		return StringPosition{}, fmt.Errorf("unsupported value type: %T", v)
	}

	pos := strings.Index(jsonStr[startPos:], searchStr)
	if pos == -1 {
		return StringPosition{}, fmt.Errorf("could not find value %v in JSON string", value)
	}

	absolutePos := startPos + pos
	return StringPosition{
		Start:  absolutePos,
		End:    absolutePos + len(searchStr),
		Length: len(searchStr),
	}, nil
}

// findChar finds the next occurrence of a character from the given position
func (t *StringIndexTracker) findChar(char byte, startPos int) int {
	jsonStr := t.originalJSON
	for i := startPos; i < len(jsonStr); i++ {
		if jsonStr[i] == char {
			return i
		}
	}
	return -1
}

// findNextElementStart finds the start position of the next array element
func (t *StringIndexTracker) findNextElementStart(currentPos int) int {
	jsonStr := t.originalJSON

	// Skip current value by counting braces/brackets
	depth := 0
	inString := false
	escape := false

	for i := currentPos; i < len(jsonStr); i++ {
		char := jsonStr[i]

		if escape {
			escape = false
			continue
		}

		if char == '\\' {
			escape = true
			continue
		}

		if char == '"' {
			inString = !inString
			continue
		}

		if inString {
			continue
		}

		switch char {
		case '{', '[':
			depth++
		case '}', ']':
			depth--
			if depth < 0 {
				return i // End of array
			}
		case ',':
			if depth == 0 {
				// Found comma at same level, next element starts after it
				for j := i + 1; j < len(jsonStr) && isWhitespace(jsonStr[j]); j++ {
					i = j
				}
				return i + 1
			}
		}
	}

	return len(jsonStr)
}

// GetStringPositionByPath returns the string position for a given JSONPath
func (t *StringIndexTracker) GetStringPositionByPath(path string) (StringPosition, bool) {
	pos, exists := t.pathToPos[path]
	return pos, exists
}

// GetPositionByPath returns the string position for a given JSONPath
func (t *StringIndexTracker) GetPositionByPath(path string) (StringPosition, bool) {
	pos, exists := t.pathToPos[path]
	return pos, exists
}

// isWhitespace checks if a character is whitespace
func isWhitespace(char byte) bool {
	return char == ' ' || char == '\t' || char == '\n' || char == '\r'
}
