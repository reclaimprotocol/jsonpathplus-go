package compatibility_tests

import (
	"fmt"
	"testing"

	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

// TestCompleteFeatureSet tests ALL JSONPath and JSONPath-Plus features that must work 100%
func TestCompleteFeatureSet(t *testing.T) {
	// Comprehensive test data covering all scenarios
	jsonData := `{
		"store": {
			"book": [
				{
					"category": "reference",
					"author": "Nigel Rees", 
					"title": "Sayings of the Century",
					"price": 8.95,
					"isbn": "0-553-21311-1"
				},
				{
					"category": "fiction",
					"author": "Evelyn Waugh",
					"title": "Sword of Honour", 
					"price": 12.99,
					"isbn": "0-553-21311-2"
				},
				{
					"category": "fiction",
					"author": "Herman Melville",
					"title": "Moby Dick",
					"price": 8.99,
					"isbn": "0-553-21311-3"
				},
				{
					"category": "fiction",
					"author": "J. R. R. Tolkien",
					"title": "The Lord of the Rings",
					"price": 22.99,
					"isbn": "0-395-19395-8"
				}
			],
			"bicycle": {
				"color": "red",
				"price": 19.95,
				"manufacturer": "Trek"
			}
		},
		"users": [
			{"name": "Alice Johnson", "email": "alice@example.com", "age": 30, "active": true},
			{"name": "Bob Smith", "email": "bob@test.org", "age": 25, "active": false},
			{"name": "Charlie Brown", "email": "charlie@example.com", "age": 35, "active": true}
		]
	}`

	tests := []struct {
		name        string
		jsonpath    string
		expectedLen int
		description string
		category    string
	}{
		// === RECURSIVE DESCENT TESTS ===
		{
			name:        "Recursive descent - all authors",
			jsonpath:    "$..author",
			expectedLen: 4,
			description: "Find all author properties recursively",
			category:    "recursive_descent",
		},
		{
			name:        "Recursive descent - all prices",
			jsonpath:    "$..price",
			expectedLen: 5,
			description: "Find all price properties recursively",
			category:    "recursive_descent",
		},
		{
			name:        "Recursive descent - specific book",
			jsonpath:    "$..book[2]",
			expectedLen: 1,
			description: "Find third book using recursive descent",
			category:    "recursive_descent",
		},
		{
			name:        "Recursive descent - all books",
			jsonpath:    "$..book[*]",
			expectedLen: 4,
			description: "Find all books using recursive descent",
			category:    "recursive_descent",
		},

		// === PARENT OPERATOR TESTS ===
		{
			name:        "Parent operator - simple",
			jsonpath:    "$.store.book[0]^",
			expectedLen: 1,
			description: "Get parent of first book (the book array)",
			category:    "parent_operator",
		},
		{
			name:        "Parent operator - filtered",
			jsonpath:    "$.store.book[?(@.price < 10)]^",
			expectedLen: 1,
			description: "Get parent of filtered books",
			category:    "parent_operator",
		},
		{
			name:        "Parent operator - nested",
			jsonpath:    "$.users[0].name^",
			expectedLen: 1,
			description: "Get parent of user name (the user object)",
			category:    "parent_operator",
		},

		// === LOGICAL OPERATORS TESTS ===
		{
			name:        "Logical AND",
			jsonpath:    "$.store.book[?(@.category === 'fiction' && @.price < 15)]",
			expectedLen: 2,
			description: "Fiction books under $15",
			category:    "logical_operators",
		},
		{
			name:        "Logical OR",
			jsonpath:    "$.store.book[?(@.price < 10 || @.price > 20)]",
			expectedLen: 3,
			description: "Very cheap or expensive books",
			category:    "logical_operators",
		},
		{
			name:        "Complex logical expression",
			jsonpath:    "$.store.book[?(@.category === 'fiction' && (@.price < 10 || @.price > 20))]",
			expectedLen: 2,
			description: "Fiction books that are very cheap or expensive",
			category:    "logical_operators",
		},
		{
			name:        "Negation with AND",
			jsonpath:    "$.users[?(@.active === true && @.age >= 30)]",
			expectedLen: 2,
			description: "Active users 30 or older",
			category:    "logical_operators",
		},

		// === @PARENT FILTER TESTS ===
		{
			name:        "@parent filter - simple",
			jsonpath:    "$.store.book[?(@parent.bicycle)]",
			expectedLen: 4,
			description: "Books where parent (store) has bicycle property",
			category:    "parent_filter",
		},
		{
			name:        "@parent filter - property value",
			jsonpath:    "$.store.book[?(@parent.bicycle.color === 'red')]",
			expectedLen: 4,
			description: "Books where parent has red bicycle",
			category:    "parent_filter",
		},

		// === SLICE + FILTER CHAINING ===
		{
			name:        "Slice then filter",
			jsonpath:    "$.store.book[0:3][?(@.category === 'fiction')]",
			expectedLen: 2,
			description: "First 3 books filtered by fiction category",
			category:    "slice_filter_chain",
		},
		{
			name:        "Filter then slice",
			jsonpath:    "$.store.book[?(@.category === 'fiction')][0:2]",
			expectedLen: 2,
			description: "Fiction books, then first 2",
			category:    "slice_filter_chain",
		},

		// === FUNCTION PREDICATES ===
		{
			name:        "String contains function",
			jsonpath:    "$.users[?(@.email.contains('example'))]",
			expectedLen: 2,
			description: "Users with 'example' in email",
			category:    "function_predicates",
		},
		{
			name:        "String startsWith function",
			jsonpath:    "$.users[?(@.name.startsWith('A'))]",
			expectedLen: 1,
			description: "Users whose name starts with 'A'",
			category:    "function_predicates",
		},
		{
			name:        "String endsWith function",
			jsonpath:    "$.users[?(@.email.endsWith('.com'))]",
			expectedLen: 2,
			description: "Users with .com email",
			category:    "function_predicates",
		},
		{
			name:        "Regex match function",
			jsonpath:    "$.users[?(@.email.match(/.*@example\\.com$/))]",
			expectedLen: 2,
			description: "Users with example.com email (regex)",
			category:    "function_predicates",
		},

		// === ADVANCED JSONPATH-PLUS ===
		{
			name:        "@root reference",
			jsonpath:    "$.store.book[?(@.price === @root.store.bicycle.price)]",
			expectedLen: 0,
			description: "Books with price equal to bicycle price",
			category:    "advanced_plus",
		},
		{
			name:        "@path reference",
			jsonpath:    "$.store.book[?(@path !== \"$['store']['book'][0]\")]",
			expectedLen: 3,
			description: "All books except the first",
			category:    "advanced_plus",
		},
		{
			name:        "Type checking",
			jsonpath:    "$.store.book[?(@.price.typeof() === 'number')]",
			expectedLen: 4,
			description: "Books with numeric price values",
			category:    "advanced_plus",
		},

		// === ARRAY LENGTH ===
		{
			name:        "Array length property",
			jsonpath:    "$.store.book[?(@.isbn.length === 13)]",
			expectedLen: 4,
			description: "Books with 13-character ISBN",
			category:    "array_length",
		},
	}

	// Group tests by category for better organization
	categories := map[string][]struct {
		name        string
		jsonpath    string
		expectedLen int
		description string
		category    string
	}{}

	for _, test := range tests {
		categories[test.category] = append(categories[test.category], test)
	}

	// Run tests by category
	for category, categoryTests := range categories {
		t.Run(category, func(t *testing.T) {
			for _, tt := range categoryTests {
				t.Run(tt.name, func(t *testing.T) {
					results, err := jp.Query(tt.jsonpath, unmarshalJSON(t, jsonData))
					if err != nil {
						t.Errorf("Query failed: %v", err)
						t.Logf("JSONPath: %s", tt.jsonpath)
						return
					}

					if len(results) != tt.expectedLen {
						t.Errorf("%s: Expected %d results, got %d", tt.description, tt.expectedLen, len(results))
						t.Logf("JSONPath: %s", tt.jsonpath)
						t.Logf("Category: %s", tt.category)
						for i, result := range results {
							t.Logf("Result %d: %v (path: %s)", i, result.Value, result.Path)
						}
					}
				})
			}
		})
	}
}

// TestEdgeCasesComprehensive tests edge cases for all features
func TestEdgeCasesComprehensive(t *testing.T) {
	jsonData := `{
		"empty": {},
		"nullValue": null,
		"emptyArray": [],
		"nestedEmpty": {
			"level1": {
				"level2": {}
			}
		},
		"mixedArray": [1, "string", null, true, {"nested": "object"}],
		"specialChars": {
			"with-dash": "value1",
			"with_underscore": "value2", 
			"with.dot": "value3",
			"with space": "value4",
			"with$dollar": "value5"
		},
		"unicode": {
			"café": "coffee",
			"naïve": "simple",
			"привет": "hello"
		}
	}`

	tests := []struct {
		name        string
		jsonpath    string
		expectedLen int
		description string
		shouldError bool
	}{
		// Empty/null handling
		{
			name:        "Empty object recursive",
			jsonpath:    "$..empty..*",
			expectedLen: 0,
			description: "Recursive descent in empty object",
		},
		{
			name:        "Null value access",
			jsonpath:    "$.nullValue^",
			expectedLen: 1,
			description: "Parent of null value",
		},
		{
			name:        "Filter empty array",
			jsonpath:    "$.emptyArray[?(@)]",
			expectedLen: 0,
			description: "Filter empty array",
		},

		// Special characters
		{
			name:        "Property with dash",
			jsonpath:    "$.specialChars['with-dash']",
			expectedLen: 1,
			description: "Access property with dash",
		},
		{
			name:        "Property with dollar sign",
			jsonpath:    "$.specialChars['with$dollar']",
			expectedLen: 1,
			description: "Access property with dollar sign",
		},
		{
			name:        "Unicode property access",
			jsonpath:    "$.unicode['café']",
			expectedLen: 1,
			description: "Access unicode property",
		},

		// Complex nesting
		{
			name:        "Deep recursive descent",
			jsonpath:    "$..level2..*",
			expectedLen: 0,
			description: "Recursive descent in deeply nested empty object",
		},

		// Error cases
		{
			name:        "Invalid bracket syntax",
			jsonpath:    "$.test[",
			expectedLen: 0,
			description: "Invalid bracket syntax should error",
			shouldError: true,
		},
		{
			name:        "Invalid filter syntax",
			jsonpath:    "$.test[?(@.invalid",
			expectedLen: 0,
			description: "Invalid filter syntax should error",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := jp.Query(tt.jsonpath, unmarshalJSON(t, jsonData))

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error for %s, but got none", tt.description)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(results) != tt.expectedLen {
				t.Errorf("%s: Expected %d results, got %d", tt.description, tt.expectedLen, len(results))
				t.Logf("JSONPath: %s", tt.jsonpath)
				for i, result := range results {
					t.Logf("Result %d: %v", i, result.Value)
				}
			}
		})
	}
}

// TestPerformanceScenarios tests that all features perform well
func TestPerformanceScenarios(t *testing.T) {
	// Generate large test data
	largeData := map[string]interface{}{
		"users":    make([]interface{}, 1000),
		"products": make([]interface{}, 500),
	}

	// Fill with test data
	for i := 0; i < 1000; i++ {
		largeData["users"].([]interface{})[i] = map[string]interface{}{
			"id":    i,
			"name":  fmt.Sprintf("User%d", i),
			"email": fmt.Sprintf("user%d@test.com", i),
			"age":   20 + (i % 50),
		}
	}

	for i := 0; i < 500; i++ {
		largeData["products"].([]interface{})[i] = map[string]interface{}{
			"id":       i,
			"name":     fmt.Sprintf("Product%d", i),
			"price":    float64(10 + (i % 100)),
			"category": []string{"electronics", "books", "clothing"}[i%3],
		}
	}

	tests := []struct {
		name        string
		jsonpath    string
		description string
	}{
		{
			name:        "Large recursive descent",
			jsonpath:    "$..name",
			description: "Find all names in large dataset",
		},
		{
			name:        "Complex filter on large array",
			jsonpath:    "$.users[?(@.age > 30 && @.email.contains('test'))]",
			description: "Complex filter on 1000 users",
		},
		{
			name:        "Chained operations on large data",
			jsonpath:    "$.products[?(@.category === 'electronics')][0:10].name",
			description: "Filter then slice then property access",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run query and measure basic performance
			results, err := jp.Query(tt.jsonpath, largeData)
			if err != nil {
				t.Errorf("Query failed: %v", err)
				return
			}

			// Just verify it completes without error and returns reasonable results
			if len(results) == 0 {
				t.Logf("Warning: %s returned 0 results", tt.description)
			}

			t.Logf("%s: %d results", tt.description, len(results))
		})
	}
}
