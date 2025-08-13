// Package jsonpathplus provides JSONPath querying with string character position tracking.
//
// This module extends the core JSONPath functionality to return the exact character positions
// of elements in the original JSON string, enabling precise location tracking for JSON editors,
// error reporting, and other applications requiring character-level precision.
package jsonpathplus

import (
	"fmt"
	"strconv"
	"strings"
)

// QueryWithStringIndex executes a JSONPath query and returns results with string indices
func QueryWithStringIndex(path string, jsonStr string) ([]Result, error) {
	// Parse JSON with string index tracking
	indexedData, err := JSONParseWithIndex(jsonStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Execute normal query using the legacy function
	results, err := QueryData(path, indexedData.Value)
	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}

	// Debug: check if we get any results
	if len(results) == 0 {
		return nil, fmt.Errorf("query '%s' returned no results", path)
	}

	// Update results with string index positions
	for i, result := range results {
		// Find string position for this result
		stringPos := findStringPositionForValue(result, jsonStr, nil)
		results[i].OriginalIndex = stringPos.Start
		results[i].Length = stringPos.Length
	}

	return results, nil
}

// findStringPositionForValue finds the character position of a value in the JSON string
func findStringPositionForValue(result Result, jsonStr string, tracker *StringIndexTracker) StringPosition {
	// Handle property access in different contexts
	if result.Path != "$" && result.ParentProperty != "" {
		// Check if this is an array element with property access (e.g., $[0].id, $[1].id)
		if strings.Contains(result.Path, "[") && strings.Contains(result.Path, "].") {
			return findPropertyInArrayElement(result, jsonStr)
		}

		// Regular property access like $.id
		if strings.Contains(result.Path, ".") {
			parts := strings.Split(result.Path, ".")
			if len(parts) > 1 {
				propName := parts[len(parts)-1]
				return findPropertyKeyPosition(propName, jsonStr)
			}
		}
	}

	// For array elements, find based on array index
	if result.Parent != nil {
		if _, isArray := result.Parent.([]interface{}); isArray {
			if idx, err := strconv.Atoi(result.ParentProperty); err == nil {
				return findArrayElementPosition(idx, result.Path, jsonStr)
			}
		}
	}

	// Try to get position from tracker using path as fallback
	if tracker != nil {
		if pos, exists := tracker.GetPositionByPath(result.Path); exists {
			return pos
		}
	}

	return StringPosition{Start: 0, End: 0, Length: 0}
}

// findPropertyInArrayElement finds a property within a specific array element
func findPropertyInArrayElement(result Result, jsonStr string) StringPosition {
	// Parse the path to extract array index and property name
	// Example: $[0].id -> array index 0, property "id"
	path := result.Path

	// Find the array index
	start := strings.Index(path, "[")
	end := strings.Index(path, "]")
	if start == -1 || end == -1 || end <= start {
		return StringPosition{Start: 0, End: 0, Length: 0}
	}

	indexStr := path[start+1 : end]
	arrayIndex, err := strconv.Atoi(indexStr)
	if err != nil {
		return StringPosition{Start: 0, End: 0, Length: 0}
	}

	// Find the property name
	dotIndex := strings.Index(path[end:], ".")
	if dotIndex == -1 || end+dotIndex+1 >= len(path) {
		return StringPosition{Start: 0, End: 0, Length: 0}
	}

	propName := path[end+dotIndex+1:]

	// Find the start of the target array element
	elementStart := findArrayElementStart(jsonStr, arrayIndex)
	if elementStart == -1 {
		return StringPosition{Start: 0, End: 0, Length: 0}
	}

	// Find the end of the target array element
	elementEnd := findArrayElementEnd(jsonStr, elementStart)
	if elementEnd == -1 {
		elementEnd = len(jsonStr)
	}

	// Search for the property within this element
	elementContent := jsonStr[elementStart:elementEnd]
	searchStr := fmt.Sprintf(`"%s"`, propName)

	pos := strings.Index(elementContent, searchStr)
	if pos == -1 {
		return StringPosition{Start: 0, End: 0, Length: 0}
	}

	absolutePos := elementStart + pos
	return StringPosition{
		Start:  absolutePos,
		End:    absolutePos + len(searchStr),
		Length: len(searchStr),
	}
}

// findArrayElementStart finds the start position of the nth array element
func findArrayElementStart(jsonStr string, index int) int {
	// Find the array start
	arrayStart := strings.Index(jsonStr, "[")
	if arrayStart == -1 {
		return -1
	}

	pos := arrayStart + 1
	currentIndex := 0

	// Skip whitespace
	for pos < len(jsonStr) && isWhitespace(jsonStr[pos]) {
		pos++
	}

	if index == 0 {
		return pos
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

	return pos
}

// findArrayElementEnd finds the end position of an array element starting at the given position
func findArrayElementEnd(jsonStr string, start int) int {
	return skipJSONElement(jsonStr, start)
}

// findPropertyKeyPosition finds the position of a property key in JSON string
func findPropertyKeyPosition(propName string, jsonStr string) StringPosition {
	searchStr := fmt.Sprintf(`"%s"`, propName)

	pos := 0
	for {
		idx := strings.Index(jsonStr[pos:], searchStr)
		if idx == -1 {
			break
		}

		absolutePos := pos + idx

		// Check if this is a property key (followed by colon after optional whitespace)
		afterKey := absolutePos + len(searchStr)
		for afterKey < len(jsonStr) && isWhitespace(jsonStr[afterKey]) {
			afterKey++
		}

		if afterKey < len(jsonStr) && jsonStr[afterKey] == ':' {
			return StringPosition{
				Start:  absolutePos,
				End:    absolutePos + len(searchStr),
				Length: len(searchStr),
			}
		}

		pos = absolutePos + 1
	}

	return StringPosition{Start: 0, End: 0, Length: 0}
}

// findArrayElementPosition finds the position of an array element
func findArrayElementPosition(index int, _ string, jsonStr string) StringPosition {
	// Find the array in the JSON string
	arrayStart := strings.Index(jsonStr, "[")
	if arrayStart == -1 {
		return StringPosition{Start: 0, End: 0, Length: 0}
	}

	// Skip to after the opening bracket
	pos := arrayStart + 1
	currentIndex := 0

	// Skip whitespace
	for pos < len(jsonStr) && isWhitespace(jsonStr[pos]) {
		pos++
	}

	// If this is index 0, we're at the first element
	if index == 0 {
		elementLength := calculateElementLength(jsonStr, pos)
		return StringPosition{Start: pos, End: pos + elementLength, Length: elementLength}
	}

	// Find the target element by skipping over previous elements
	for currentIndex < index && pos < len(jsonStr) {
		// Skip current element
		pos = skipJSONElement(jsonStr, pos)

		// Skip whitespace and comma
		for pos < len(jsonStr) && (isWhitespace(jsonStr[pos]) || jsonStr[pos] == ',') {
			pos++
		}
		currentIndex++

		if currentIndex == index {
			elementLength := calculateElementLength(jsonStr, pos)
			return StringPosition{Start: pos, End: pos + elementLength, Length: elementLength}
		}
	}

	elementLength := calculateElementLength(jsonStr, pos)
	return StringPosition{Start: pos, End: pos + elementLength, Length: elementLength}
}

// skipJSONElement skips over a complete JSON element
func skipJSONElement(jsonStr string, start int) int {
	pos := start
	if pos >= len(jsonStr) {
		return pos
	}

	char := jsonStr[pos]
	switch char {
	case '"':
		// Skip string
		pos++ // Skip opening quote
		for pos < len(jsonStr) {
			if jsonStr[pos] == '"' && (pos == start+1 || jsonStr[pos-1] != '\\') {
				pos++ // Skip closing quote
				return pos
			}
			pos++
		}
	case '{':
		// Skip object
		return skipBracedElement(jsonStr, pos, '{', '}')
	case '[':
		// Skip array
		return skipBracedElement(jsonStr, pos, '[', ']')
	default:
		// Skip primitive (number, boolean, null)
		for pos < len(jsonStr) && !isWhitespace(jsonStr[pos]) &&
			jsonStr[pos] != ',' && jsonStr[pos] != '}' && jsonStr[pos] != ']' {
			pos++
		}
	}

	return pos
}

// skipBracedElement skips over a complete braced element (object or array)
func skipBracedElement(jsonStr string, start int, openBrace, closeBrace byte) int {
	pos := start + 1 // Skip opening brace
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

// calculateElementLength calculates the length of a JSON element starting at the given position
func calculateElementLength(jsonStr string, start int) int {
	if start >= len(jsonStr) {
		return 0
	}

	end := skipJSONElement(jsonStr, start)
	return end - start
}
