package compatibility_tests

import (
	"testing"

	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

// TestStringIndexFunctionality tests comprehensive string index finding capabilities
func TestStringIndexFunctionality(t *testing.T) {
	jsonData := `{
		"usernames": ["alice", "bob", "charlie", "alice", "diana", "bob"],
		"tags": ["go", "json", "path", "go", "programming", "json"],
		"mixed": [1, "test", 3, "test", 5, true, "test"],
		"people": [
			{"name": "Alice", "age": 30},
			{"name": "Bob", "age": 25},
			{"name": "Alice", "age": 35},
			{"name": "Charlie", "age": 40}
		],
		"nested": {
			"level1": {
				"items": ["first", "second", "first", "third"]
			}
		}
	}`

	tests := []struct {
		name        string
		jsonpath    string
		expectedLen int
		description string
		validate    func([]jp.Result) bool
	}{
		{
			name:        "Find all string occurrences",
			jsonpath:    "$.usernames[?(@==='alice')]",
			expectedLen: 2,
			description: "Find all occurrences of 'alice'",
			validate: func(results []jp.Result) bool {
				for _, r := range results {
					if r.Value.(string) != "alice" {
						return false
					}
				}
				return true
			},
		},
		{
			name:        "Find string in mixed array",
			jsonpath:    "$.mixed[?(@==='test')]",
			expectedLen: 3,
			description: "Find all 'test' strings in mixed array",
			validate: func(results []jp.Result) bool {
				for _, r := range results {
					if r.Value.(string) != "test" {
						return false
					}
				}
				return true
			},
		},
		{
			name:        "Find by object property value",
			jsonpath:    "$.people[?(@.name==='Alice')]",
			expectedLen: 2,
			description: "Find people named Alice",
			validate: func(results []jp.Result) bool {
				for _, r := range results {
					person := r.Value.(map[string]interface{})
					if person["name"].(string) != "Alice" {
						return false
					}
				}
				return true
			},
		},
		{
			name:        "Find first occurrence by index",
			jsonpath:    "$.usernames[0]",
			expectedLen: 1,
			description: "Get first username",
			validate: func(results []jp.Result) bool {
				return results[0].Value.(string) == "alice"
			},
		},
		{
			name:        "Find last occurrence by negative index",
			jsonpath:    "$.usernames[-1]",
			expectedLen: 1,
			description: "Get last username",
			validate: func(results []jp.Result) bool {
				return results[0].Value.(string) == "bob"
			},
		},
		{
			name:        "Find specific index occurrences",
			jsonpath:    "$.usernames[1,3,5]",
			expectedLen: 3,
			description: "Get usernames at specific indices",
			validate: func(results []jp.Result) bool {
				expected := []string{"bob", "alice", "bob"}
				for i, r := range results {
					if r.Value.(string) != expected[i] {
						return false
					}
				}
				return true
			},
		},
		{
			name:        "Find range of string indices",
			jsonpath:    "$.tags[1:4]",
			expectedLen: 3,
			description: "Get tags from index 1 to 3",
			validate: func(results []jp.Result) bool {
				expected := []string{"json", "path", "go"}
				for i, r := range results {
					if r.Value.(string) != expected[i] {
						return false
					}
				}
				return true
			},
		},
		{
			name:        "Find nested string occurrences",
			jsonpath:    "$.nested.level1.items[?(@==='first')]",
			expectedLen: 2,
			description: "Find 'first' in nested array",
			validate: func(results []jp.Result) bool {
				for _, r := range results {
					if r.Value.(string) != "first" {
						return false
					}
				}
				return true
			},
		},
		{
			name:        "Find strings with recursive descent",
			jsonpath:    "$..items[?(@==='first')]",
			expectedLen: 2,
			description: "Recursively find 'first' strings",
			validate: func(results []jp.Result) bool {
				for _, r := range results {
					if r.Value.(string) != "first" {
						return false
					}
				}
				return true
			},
		},
		{
			name:        "Count string occurrences via length",
			jsonpath:    "$.usernames[?(@==='bob')]",
			expectedLen: 2,
			description: "Count occurrences of 'bob'",
			validate: func(results []jp.Result) bool {
				return len(results) == 2
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := jp.Query(tt.jsonpath, unmarshalJSON(t, jsonData))
			if err != nil {
				t.Errorf("Query failed: %v", err)
				return
			}

			if len(results) != tt.expectedLen {
				t.Errorf("%s: Expected %d results, got %d", tt.description, tt.expectedLen, len(results))
				t.Logf("JSONPath: %s", tt.jsonpath)
				for i, result := range results {
					t.Logf("Result %d: %v (path: %s)", i, result.Value, result.Path)
				}
				return
			}

			if tt.validate != nil && !tt.validate(results) {
				t.Errorf("%s: Validation failed", tt.description)
				for i, result := range results {
					t.Logf("Result %d: %v", i, result.Value)
				}
			}
		})
	}
}

// TestStringIndexWithMetadata tests string finding with result metadata
func TestStringIndexWithMetadata(t *testing.T) {
	jsonData := `{
		"documents": [
			{"id": 1, "content": "hello world"},
			{"id": 2, "content": "hello universe"},
			{"id": 3, "content": "goodbye world"}
		]
	}`

	tests := []struct {
		name        string
		jsonpath    string
		expectedLen int
		description string
		checkPaths  []string
	}{
		// Note: indexOf functionality would go here when implemented
		{
			name:        "Find by index with path",
			jsonpath:    "$.documents[0,2]",
			expectedLen: 2,
			description: "Get first and third documents",
			checkPaths:  []string{"$.documents[0]", "$.documents[2]"},
		},
		{
			name:        "Find all content strings",
			jsonpath:    "$.documents[*].content",
			expectedLen: 3,
			description: "Get all content strings",
			checkPaths:  []string{"$.documents[0].content", "$.documents[1].content", "$.documents[2].content"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			results, err := jp.Query(tt.jsonpath, unmarshalJSON(t, jsonData))
			if err != nil {
				t.Errorf("Query failed: %v", err)
				return
			}

			if len(results) != tt.expectedLen {
				t.Errorf("%s: Expected %d results, got %d", tt.description, tt.expectedLen, len(results))
				return
			}

			if tt.checkPaths != nil {
				for i, result := range results {
					if i < len(tt.checkPaths) && result.Path != tt.checkPaths[i] {
						t.Errorf("Result %d path: expected %s, got %s", i, tt.checkPaths[i], result.Path)
					}
				}
			}
		})
	}
}

// TestStringIndexPerformance tests performance characteristics of string finding
func TestStringIndexPerformance(t *testing.T) {
	// Generate test data with many string occurrences
	largeArray := make([]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		if i%10 == 0 {
			largeArray[i] = "target"
		} else {
			largeArray[i] = "other"
		}
	}

	jsonStr, err := jp.JSONParse(`{"large_array": []}`)
	if err != nil {
		t.Fatalf("Failed to create test data: %v", err)
	}

	// Manually set the large array
	data := jsonStr.(map[string]interface{})
	data["large_array"] = largeArray

	tests := []struct {
		name        string
		jsonpath    string
		expectedLen int
		description string
	}{
		{
			name:        "Find in large array",
			jsonpath:    "$.large_array[?(@==='target')]",
			expectedLen: 100,
			description: "Find all 'target' strings in large array",
		},
		{
			name:        "Find first occurrence in large array",
			jsonpath:    "$.large_array[0]",
			expectedLen: 1,
			description: "Get first element of large array",
		},
		{
			name:        "Find last occurrence in large array",
			jsonpath:    "$.large_array[-1]",
			expectedLen: 1,
			description: "Get last element of large array",
		},
		{
			name:        "Slice large array",
			jsonpath:    "$.large_array[0:10]",
			expectedLen: 10,
			description: "Get first 10 elements",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := jp.Query(tt.jsonpath, data)
			if err != nil {
				t.Errorf("Query failed: %v", err)
				return
			}

			if len(results) != tt.expectedLen {
				t.Errorf("%s: Expected %d results, got %d", tt.description, tt.expectedLen, len(results))
			}
		})
	}
}

// TestStringIndexEdgeCases tests edge cases in string index finding
func TestStringIndexEdgeCases(t *testing.T) {
	jsonData := `{
		"empty_strings": ["", "hello", "", "world", ""],
		"whitespace": [" ", "  ", "\t", "\n", "text"],
		"special_chars": ["@", "#", "$", "%", "&"],
		"unicode": ["café", "naïve", "résumé", "привет"],
		"numbers_as_strings": ["1", "2", "3", "1", "2"],
		"mixed_types": [1, "1", true, "true", null, "null"]
	}`

	tests := []struct {
		name        string
		jsonpath    string
		expectedLen int
		description string
	}{
		{
			name:        "Find empty strings",
			jsonpath:    "$.empty_strings[?(@==='')]",
			expectedLen: 3,
			description: "Find all empty strings",
		},
		{
			name:        "Find whitespace strings",
			jsonpath:    "$.whitespace[?(@===' ')]",
			expectedLen: 1,
			description: "Find single space strings",
		},
		{
			name:        "Find special character strings",
			jsonpath:    "$.special_chars[?(@==='$')]",
			expectedLen: 1,
			description: "Find dollar sign string",
		},
		{
			name:        "Find unicode strings",
			jsonpath:    "$.unicode[?(@==='café')]",
			expectedLen: 1,
			description: "Find unicode string",
		},
		{
			name:        "Find number-like strings",
			jsonpath:    "$.numbers_as_strings[?(@==='1')]",
			expectedLen: 2,
			description: "Find string '1' occurrences",
		},
		{
			name:        "Find string vs number distinction",
			jsonpath:    "$.mixed_types[?(@==='1')]",
			expectedLen: 1,
			description: "Find string '1' (not number 1)",
		},
		{
			name:        "Find string vs boolean distinction",
			jsonpath:    "$.mixed_types[?(@==='true')]",
			expectedLen: 1,
			description: "Find string 'true' (not boolean true)",
		},
		{
			name:        "Find string vs null distinction",
			jsonpath:    "$.mixed_types[?(@==='null')]",
			expectedLen: 1,
			description: "Find string 'null' (not null value)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := jp.Query(tt.jsonpath, unmarshalJSON(t, jsonData))
			if err != nil {
				t.Errorf("Query failed: %v", err)
				return
			}

			if len(results) != tt.expectedLen {
				t.Errorf("%s: Expected %d results, got %d", tt.description, tt.expectedLen, len(results))
				t.Logf("JSONPath: %s", tt.jsonpath)
				for i, result := range results {
					t.Logf("Result %d: %v (type: %T)", i, result.Value, result.Value)
				}
			}
		})
	}
}
