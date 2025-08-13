// Package jsonpathplus provides JSONPath querying with string character position tracking.
//
// This file contains the core string position tracking functionality that maps JSON elements
// to their exact character positions in the original JSON string.
package jsonpathplus

import (
	"encoding/json"
)

// StringPosition represents a position in the original JSON string.
type StringPosition struct {
	Start  int // Starting character position
	End    int // Ending character position
	Length int // Length of the element
}

// IndexedValue wraps a value with its string position information.
type IndexedValue struct {
	Value    interface{}
	Position StringPosition
}

// StringIndexTracker tracks character positions during JSON parsing.
type StringIndexTracker struct {
	originalJSON string
	pathToPos    map[string]StringPosition
}

// NewStringIndexTracker creates a new string index tracker.
func NewStringIndexTracker(jsonStr string) *StringIndexTracker {
	return &StringIndexTracker{
		originalJSON: jsonStr,
		pathToPos:    make(map[string]StringPosition),
	}
}

// JSONParseWithIndex parses JSON while preserving string positions.
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

// GetStringPositionByPath returns the string position for a given JSONPath.
func (t *StringIndexTracker) GetStringPositionByPath(path string) (StringPosition, bool) {
	pos, exists := t.pathToPos[path]
	return pos, exists
}

// GetPositionByPath returns the string position for a given JSONPath.
func (t *StringIndexTracker) GetPositionByPath(path string) StringPosition {
	if pos, exists := t.pathToPos[path]; exists {
		return pos
	}
	return StringPosition{}
}

// isWhitespace checks if a character is whitespace.
func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}
