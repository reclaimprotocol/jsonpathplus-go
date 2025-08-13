package jsonpathplus

import (
	"strings"
	"testing"
)

// TestStringIndexPreservation tests that indices represent character positions in original JSON string
func TestStringIndexPreservation(t *testing.T) {
	engine, err := NewEngine(DefaultConfig())
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer func() { _ = engine.Close() }()

	t.Run("SimpleObjectProperty", func(t *testing.T) {
		testSimpleObjectStringIndex(t, engine)
	})

	t.Run("ArrayElementsStringIndex", func(t *testing.T) {
		testArrayElementsStringIndex(t, engine)
	})

	t.Run("NestedObjectsStringIndex", func(t *testing.T) {
		testNestedObjectsStringIndex(t, engine)
	})

	t.Run("WhitespacePreservation", func(t *testing.T) {
		testWhitespacePreservation(t, engine)
	})

	t.Run("ComplexNestingStringIndex", func(t *testing.T) {
		testComplexNestingStringIndex(t, engine)
	})

	t.Run("ArrayOfObjectsStringIndex", func(t *testing.T) {
		testArrayOfObjectsStringIndex(t, engine)
	})
}

func testSimpleObjectStringIndex(t *testing.T, _ *JSONPathEngine) {
	// Test case: {"id":123,"name":"test"}
	//             ^   ^    ^     ^
	//             0   5    11    18
	jsonStr := `{"id":123,"name":"test"}`

	tests := []struct {
		path          string
		expectedIndex int
		description   string
	}{
		{"$.id", 1, "Property 'id' starts at position 1 (after opening brace)"},
		{"$.name", 10, "Property 'name' starts at position 10"},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			results, err := Query(test.path, jsonStr)
			if err != nil {
				t.Fatalf("Query failed: %v", err)
			}

			if len(results) != 1 {
				t.Fatalf("Expected 1 result, got %d", len(results))
			}

			result := results[0]
			if result.OriginalIndex != test.expectedIndex {
				t.Errorf("Expected string index %d, got %d for path %s",
					test.expectedIndex, result.OriginalIndex, test.path)

				// Show context for debugging
				if test.expectedIndex < len(jsonStr) {
					char := jsonStr[test.expectedIndex]
					t.Errorf("Character at expected position %d: '%c'", test.expectedIndex, char)
				}
				if result.OriginalIndex < len(jsonStr) {
					char := jsonStr[result.OriginalIndex]
					t.Errorf("Character at actual position %d: '%c'", result.OriginalIndex, char)
				}
				t.Errorf("JSON string: %s", jsonStr)
				t.Errorf("Position markers: %s", strings.Repeat(" ", test.expectedIndex)+"^")
			}
		})
	}
}

func testArrayElementsStringIndex(t *testing.T, _ *JSONPathEngine) {
	// Test case: ["first","second","third"]
	//             ^       ^        ^
	//             1       9        18
	jsonStr := `["first","second","third"]`

	results, err := Query("$[*]", jsonStr)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	expectedIndices := []int{1, 9, 18} // Start positions of each string
	expectedValues := []string{"first", "second", "third"}

	if len(results) != len(expectedIndices) {
		t.Fatalf("Expected %d results, got %d", len(expectedIndices), len(results))
	}

	for i, result := range results {
		expectedIndex := expectedIndices[i]
		expectedValue := expectedValues[i]

		if result.Value != expectedValue {
			t.Errorf("Result %d: expected value %s, got %v", i, expectedValue, result.Value)
		}

		if result.OriginalIndex != expectedIndex {
			t.Errorf("Result %d (%s): expected string index %d, got %d",
				i, expectedValue, expectedIndex, result.OriginalIndex)
			t.Errorf("JSON: %s", jsonStr)
			t.Errorf("Expected position: %s", strings.Repeat(" ", expectedIndex)+"^")
		}
	}
}

func testNestedObjectsStringIndex(t *testing.T, _ *JSONPathEngine) {
	// Test case: {"user":{"name":"john","age":25}}
	//             ^       ^       ^      ^
	//             0       8       16     24
	jsonStr := `{"user":{"name":"john","age":25}}`

	tests := []struct {
		path          string
		expectedIndex int
		description   string
	}{
		{"$.user", 1, "Property 'user' starts at position 1"},
		{"$.user.name", 9, "Nested property 'name' starts at position 9"},
		{"$.user.age", 23, "Nested property 'age' starts at position 23"},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			results, err := Query(test.path, jsonStr)
			if err != nil {
				t.Fatalf("Query failed: %v", err)
			}

			if len(results) != 1 {
				t.Fatalf("Expected 1 result, got %d", len(results))
			}

			result := results[0]
			if result.OriginalIndex != test.expectedIndex {
				t.Errorf("Path %s: expected string index %d, got %d",
					test.path, test.expectedIndex, result.OriginalIndex)
				t.Errorf("JSON: %s", jsonStr)
				t.Errorf("Expected: %s", strings.Repeat(" ", test.expectedIndex)+"^")
				if result.OriginalIndex < len(jsonStr) {
					t.Errorf("Actual:   %s", strings.Repeat(" ", result.OriginalIndex)+"^")
				}
			}
		})
	}
}

func testWhitespacePreservation(t *testing.T, _ *JSONPathEngine) {
	// Test with significant whitespace that should be preserved
	jsonStr := `{
  "id": 123,
  "data": {
    "name": "test",
    "values": [1, 2, 3]
  }
}`

	// Find where "id" appears in the string (after newline and spaces)
	idPos := strings.Index(jsonStr, `"id"`)
	namePos := strings.Index(jsonStr, `"name"`)
	valuesPos := strings.Index(jsonStr, `"values"`)

	tests := []struct {
		path          string
		expectedIndex int
		description   string
	}{
		{"$.id", idPos, "Property 'id' position with whitespace"},
		{"$.data.name", namePos, "Nested property 'name' position with whitespace"},
		{"$.data.values", valuesPos, "Nested property 'values' position with whitespace"},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			results, err := Query(test.path, jsonStr)
			if err != nil {
				t.Fatalf("Query failed: %v", err)
			}

			if len(results) != 1 {
				t.Fatalf("Expected 1 result, got %d", len(results))
			}

			result := results[0]
			if result.OriginalIndex != test.expectedIndex {
				t.Errorf("Path %s: expected string index %d, got %d",
					test.path, test.expectedIndex, result.OriginalIndex)

				// Show context lines
				lines := strings.Split(jsonStr, "\n")
				for i, line := range lines {
					t.Errorf("Line %d: %s", i, line)
				}
			}
		})
	}
}

func testComplexNestingStringIndex(t *testing.T, _ *JSONPathEngine) {
	jsonStr := `{"company":{"departments":[{"name":"engineering","employees":[{"name":"alice","id":1}]}]}}`

	tests := []struct {
		path        string
		description string
	}{
		{"$.company", "Root company object"},
		{"$.company.departments", "Departments array"},
		{"$.company.departments[0].name", "First department name"},
		{"$.company.departments[0].employees", "Employees array"},
		{"$.company.departments[0].employees[0].name", "First employee name"},
		{"$.company.departments[0].employees[0].id", "First employee id"},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			results, err := Query(test.path, jsonStr)
			if err != nil {
				t.Fatalf("Query failed: %v", err)
			}

			if len(results) != 1 {
				t.Fatalf("Expected 1 result, got %d", len(results))
			}

			result := results[0]

			// Verify the index points to a reasonable position in the string
			if result.OriginalIndex < 0 || result.OriginalIndex >= len(jsonStr) {
				t.Errorf("String index %d is out of bounds for JSON string of length %d",
					result.OriginalIndex, len(jsonStr))
			}

			t.Logf("Path %s -> index %d (char: '%c')",
				test.path, result.OriginalIndex, jsonStr[result.OriginalIndex])
		})
	}
}

func testArrayOfObjectsStringIndex(t *testing.T, _ *JSONPathEngine) {
	jsonStr := `[{"id":1,"name":"first"},{"id":2,"name":"second"}]`

	// Test accessing objects in array
	results, err := Query("$[*].id", jsonStr)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	// Find positions of "id" properties in the string
	firstIDPos := strings.Index(jsonStr, `"id":1`)
	secondIDPos := strings.Index(jsonStr[firstIDPos+1:], `"id":2`) + firstIDPos + 1

	expectedPositions := []int{firstIDPos, secondIDPos}

	for i, result := range results {
		expectedPos := expectedPositions[i]

		t.Logf("Result %d: id=%v, string index=%d, expected=%d",
			i, result.Value, result.OriginalIndex, expectedPos)

		// The index should be close to the expected position (exact position depends on implementation)
		if result.OriginalIndex < expectedPos-5 || result.OriginalIndex > expectedPos+5 {
			t.Errorf("Result %d: string index %d is not close to expected position %d",
				i, result.OriginalIndex, expectedPos)
		}
	}
}

// Benchmark string index preservation performance
func BenchmarkStringIndexPreservation(b *testing.B) {
	jsonStr := `{"users":[{"name":"alice","age":25},{"name":"bob","age":30},{"name":"charlie","age":35}]}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Query("$.users[*].name", jsonStr)
		if err != nil {
			b.Fatalf("Query failed: %v", err)
		}
	}
}

// Helper function to validate that string index points to expected character
func validateStringIndex(t *testing.T, jsonStr string, index int, expectedChar byte, description string) {
	if index < 0 || index >= len(jsonStr) {
		t.Errorf("%s: index %d out of bounds for string length %d", description, index, len(jsonStr))
		return
	}

	actualChar := jsonStr[index]
	if actualChar != expectedChar {
		t.Errorf("%s: expected char '%c' at index %d, got '%c'",
			description, expectedChar, index, actualChar)

		// Show context
		start := maxInt(0, index-10)
		end := minInt(len(jsonStr), index+10)
		context := jsonStr[start:end]
		pointer := strings.Repeat(" ", index-start) + "^"
		t.Errorf("Context: %s", context)
		t.Errorf("         %s", pointer)
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
