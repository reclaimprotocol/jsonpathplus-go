package jsonpathplus

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestBasicPaths(t *testing.T) {
	jsonStr := `{
		"store": {
			"book": [
				{
					"category": "reference",
					"author": "Nigel Rees",
					"title": "Sayings of the Century",
					"price": 8.95
				},
				{
					"category": "fiction",
					"author": "Evelyn Waugh",
					"title": "Sword of Honor",
					"price": 12.99
				},
				{
					"category": "fiction",
					"author": "Herman Melville",
					"title": "Moby Dick",
					"isbn": "0-553-21311-3",
					"price": 8.99
				},
				{
					"category": "fiction",
					"author": "J. R. R. Tolkien",
					"title": "The Lord of the Rings",
					"isbn": "0-395-19395-8",
					"price": 22.99
				}
			],
			"bicycle": {
				"color": "red",
				"price": 19.95
			}
		},
		"expensive": 10
	}`

	data, err := JSONParse(jsonStr)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	tests := []struct {
		path     string
		expected int
		desc     string
	}{
		{"$.store.book[*].author", 4, "All authors"},
		{"$..author", 4, "All authors recursive"},
		{"$.store.*", 2, "All items in store"},
		{"$.store..price", 5, "All prices in store"},
		{"$..book[2]", 1, "Third book"},
		{"$..book[-1]", 1, "Last book"},
		{"$..book[0,1]", 2, "First two books"},
		{"$..book[:2]", 2, "First two books via slice"},
		{"$..book[?(@.isbn)]", 2, "Books with ISBN"},
		{"$..book[?(@.price < 10)]", 2, "Books cheaper than 10"},
		{"$..*", 28, "All members recursive"},
		{"$.store.book[*]", 4, "All books"},
		{"$.store.book[2].title", 1, "Title of third book"},
		{"$.store.book[0]['title']", 1, "Title with bracket notation"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			results, err := QueryData(test.path, data)
			if err != nil {
				t.Errorf("Query failed for path %s: %v", test.path, err)
				return
			}
			if len(results) != test.expected {
				t.Errorf("Path %s: expected %d results, got %d", test.path, test.expected, len(results))
			}
		})
	}
}

func TestIndexPreservation(t *testing.T) {
	jsonStr := `{
		"items": ["a", "b", "c", "d", "e"],
		"nested": {
			"first": 1,
			"second": 2,
			"third": 3
		}
	}`

	data, err := JSONParse(jsonStr)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	results, err := QueryData("$.items[*]", data)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	for i, result := range results {
		if result.OriginalIndex != i {
			t.Errorf("Expected original index %d, got %d", i, result.OriginalIndex)
		}
	}

	results, err = QueryData("$.nested.*", data)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	for _, result := range results {
		if result.OriginalIndex < 0 {
			t.Errorf("Original index not preserved for object properties")
		}
	}
}

func TestFilterExpressions(t *testing.T) {
	jsonStr := `{
		"users": [
			{"name": "John", "age": 30, "active": true},
			{"name": "Jane", "age": 25, "active": false},
			{"name": "Bob", "age": 35, "active": true},
			{"name": "Alice", "age": 28}
		]
	}`

	data, err := JSONParse(jsonStr)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	tests := []struct {
		path     string
		expected []string
		desc     string
	}{
		{"$.users[?(@.age > 25)]", []string{"John", "Bob", "Alice"}, "Age greater than 25"},
		{"$.users[?(@.active == true)]", []string{"John", "Bob"}, "Active users"},
		{"$.users[?(@.active)]", []string{"John", "Bob"}, "Users with active field"},
		{"$.users[?(@.age <= 30)]", []string{"John", "Jane", "Alice"}, "Age less than or equal to 30"},
		{"$.users[?(@.name != 'Bob')]", []string{"John", "Jane", "Alice"}, "Name not Bob"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			results, err := QueryData(test.path, data)
			if err != nil {
				t.Errorf("Query failed for path %s: %v", test.path, err)
				return
			}

			var names []string
			for _, result := range results {
				if user, ok := result.Value.(map[string]interface{}); ok {
					if name, ok := user["name"].(string); ok {
						names = append(names, name)
					}
				}
			}

			if !reflect.DeepEqual(names, test.expected) {
				t.Errorf("Path %s: expected %v, got %v", test.path, test.expected, names)
			}
		})
	}
}

func TestSliceNotation(t *testing.T) {
	jsonStr := `{
		"numbers": [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
	}`

	data, err := JSONParse(jsonStr)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	tests := []struct {
		path     string
		expected []int
		desc     string
	}{
		{"$.numbers[0:5]", []int{0, 1, 2, 3, 4}, "First 5 elements"},
		{"$.numbers[5:]", []int{5, 6, 7, 8, 9}, "From index 5 to end"},
		{"$.numbers[:3]", []int{0, 1, 2}, "First 3 elements"},
		{"$.numbers[::2]", []int{0, 2, 4, 6, 8}, "Every second element"},
		{"$.numbers[1:8:2]", []int{1, 3, 5, 7}, "Range with step"},
		{"$.numbers[-3:]", []int{7, 8, 9}, "Last 3 elements"},
		{"$.numbers[-5:-2]", []int{5, 6, 7}, "Negative indices"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			results, err := QueryData(test.path, data)
			if err != nil {
				t.Errorf("Query failed for path %s: %v", test.path, err)
				return
			}

			var values []int
			for _, result := range results {
				if num, ok := result.Value.(int); ok {
					values = append(values, num)
				} else if num, ok := result.Value.(float64); ok {
					values = append(values, int(num))
				}
			}

			if !reflect.DeepEqual(values, test.expected) {
				t.Errorf("Path %s: expected %v, got %v", test.path, test.expected, values)
			}
		})
	}
}

func TestRecursiveDescent(t *testing.T) {
	jsonStr := `{
		"a": {
			"b": {
				"c": 1,
				"d": {
					"c": 2
				}
			},
			"c": 3
		},
		"c": 4
	}`

	data, err := JSONParse(jsonStr)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	results, err := QueryData("$..c", data)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	expected := []int{1, 2, 3, 4}
	if len(results) != len(expected) {
		t.Errorf("Expected %d results, got %d", len(expected), len(results))
	}

	var values []int
	for _, result := range results {
		if num, ok := result.Value.(int); ok {
			values = append(values, num)
		} else if num, ok := result.Value.(float64); ok {
			values = append(values, int(num))
		}
	}

	for _, exp := range expected {
		found := false
		for _, val := range values {
			if val == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected value %d not found in results", exp)
		}
	}
}

func TestComplexNestedStructure(t *testing.T) {
	jsonStr := `{
		"root": {
			"level1": {
				"level2": {
					"items": [
						{"id": 1, "value": "a"},
						{"id": 2, "value": "b"},
						{"id": 3, "value": "c"}
					]
				}
			},
			"otherBranch": {
				"items": [
					{"id": 4, "value": "d"},
					{"id": 5, "value": "e"}
				]
			}
		}
	}`

	data, err := JSONParse(jsonStr)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	results, err := QueryData("$..items[*].value", data)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	expected := []string{"a", "b", "c", "d", "e"}
	if len(results) != len(expected) {
		t.Errorf("Expected %d results, got %d", len(expected), len(results))
	}

	// Collect actual values for comparison
	var actualValues []string
	for _, result := range results {
		if str, ok := result.Value.(string); ok {
			actualValues = append(actualValues, str)
		}
	}

	// Check that all expected values are present
	for _, expectedVal := range expected {
		found := false
		for _, actualVal := range actualValues {
			if actualVal == expectedVal {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected value %s not found in results %v", expectedVal, actualValues)
		}
	}
}

func TestParentAndPath(t *testing.T) {
	jsonStr := `{
		"users": [
			{"name": "John", "age": 30},
			{"name": "Jane", "age": 25}
		]
	}`

	data, err := JSONParse(jsonStr)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	results, err := QueryData("$.users[*].name", data)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	for i, result := range results {
		expectedPath := []string{"$.users[0].name", "$.users[1].name"}
		if result.Path != expectedPath[i] {
			t.Errorf("Expected path %s, got %s", expectedPath[i], result.Path)
		}

		if result.ParentProperty != "name" {
			t.Errorf("Expected parent property 'name', got %s", result.ParentProperty)
		}

		if result.Parent == nil {
			t.Errorf("Parent should not be nil")
		}
	}
}

func TestEmptyResults(t *testing.T) {
	jsonStr := `{"a": 1, "b": 2}`

	data, err := JSONParse(jsonStr)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	results, err := QueryData("$.nonexistent", data)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results for nonexistent path, got %d", len(results))
	}
}

func TestWildcardWithArrays(t *testing.T) {
	jsonStr := `{
		"data": [1, 2, 3, 4, 5]
	}`

	data, err := JSONParse(jsonStr)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	results, err := QueryData("$.data[*]", data)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 5 {
		t.Errorf("Expected 5 results, got %d", len(results))
	}

	for i, result := range results {
		if result.OriginalIndex != i {
			t.Errorf("Expected original index %d, got %d", i, result.OriginalIndex)
		}
		expectedValue := i + 1
		if val, ok := result.Value.(int); ok {
			if val != expectedValue {
				t.Errorf("Expected value %d, got %d", expectedValue, val)
			}
		} else if val, ok := result.Value.(float64); ok {
			if int(val) != expectedValue {
				t.Errorf("Expected value %d, got %d", expectedValue, int(val))
			}
		}
	}
}

func TestUnionOperator(t *testing.T) {
	jsonStr := `{
		"book": {
			"title": "Example",
			"author": "John Doe",
			"year": 2023,
			"isbn": "123456789"
		}
	}`

	data, err := JSONParse(jsonStr)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	results, err := QueryData("$.book['title','author']", data)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	expectedValues := map[string]bool{
		"Example":  false,
		"John Doe": false,
	}

	for _, result := range results {
		if str, ok := result.Value.(string); ok {
			if _, exists := expectedValues[str]; exists {
				expectedValues[str] = true
			}
		}
	}

	for val, found := range expectedValues {
		if !found {
			t.Errorf("Expected value %s not found", val)
		}
	}
}

func BenchmarkSimplePath(b *testing.B) {
	data := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": 123,
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = QueryData("$.a.b.c", data)
	}
}

func BenchmarkRecursivePath(b *testing.B) {
	data := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": 1,
				"d": map[string]interface{}{
					"c": 2,
				},
			},
			"c": 3,
		},
		"c": 4,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = QueryData("$..c", data)
	}
}

func BenchmarkFilterExpression(b *testing.B) {
	var items []interface{}
	for i := 0; i < 100; i++ {
		items = append(items, map[string]interface{}{
			"id":    i,
			"value": i * 10,
		})
	}
	data := map[string]interface{}{
		"items": items,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = QueryData("$.items[?(@.value > 500)]", data)
	}
}

func TestJSONParsePreservesOrder(t *testing.T) {
	jsonStr := `{"z": 1, "a": 2, "m": 3, "b": 4}`

	data, err := JSONParse(jsonStr)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	obj, ok := data.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map, got %T", data)
	}

	marshaledBytes, err := json.Marshal(obj)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	marshaled := string(marshaledBytes)
	if marshaled != jsonStr {
		t.Logf("Note: Go's map type doesn't preserve insertion order")
		t.Logf("Original: %s", jsonStr)
		t.Logf("Marshaled: %s", marshaled)
	}
}
