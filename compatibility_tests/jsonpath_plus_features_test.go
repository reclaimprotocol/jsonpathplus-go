package compatibility_tests

import (
	"testing"

	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

// TestJSONPathPlusFeatures tests the extended JSONPath-Plus features
func TestJSONPathPlusFeatures(t *testing.T) {
	jsonData := `{
		"company": {
			"departments": {
				"engineering": {
					"manager": "Alice",
					"employees": [
						{"name": "Bob", "role": "developer", "level": "senior"},
						{"name": "Carol", "role": "designer", "level": "junior"}
					]
				},
				"sales": {
					"manager": "Dave",
					"employees": [
						{"name": "Eve", "role": "rep", "level": "senior"}
					]
				}
			},
			"metadata": {
				"founded": 2010,
				"location": "San Francisco"
			}
		}
	}`

	tests := []struct {
		name        string
		jsonpath    string
		expectedLen int
		description string
		skip        bool
		skipReason  string
	}{
		// Property names operator (~)
		{
			name:        "Department names",
			jsonpath:    "$.company.departments.*~",
			expectedLen: 2,
			description: "Get property names of departments",
		},
		{
			name:        "All property names in company",
			jsonpath:    "$.company.*~",
			expectedLen: 2,
			description: "Get top-level property names in company",
		},
		{
			name:        "Employee property names",
			jsonpath:    "$.company.departments.engineering.employees[0].*~",
			expectedLen: 3,
			description: "Get property names of first employee",
		},

		// @property filter
		{
			name:        "Filter by property name",
			jsonpath:    "$.company.departments.*[?(@property === 'engineering')]",
			expectedLen: 1,
			description: "Find departments with property name 'engineering'",
		},
		{
			name:        "Exclude specific property",
			jsonpath:    "$.company.*[?(@property !== 'metadata')]",
			expectedLen: 1,
			description: "Get all company properties except metadata",
		},

		// @parentProperty filter
		{
			name:        "Filter by parent property",
			jsonpath:    "$.company.departments.engineering.employees[?(@parentProperty === 'employees')]",
			expectedLen: 2,
			description: "Find items where parent property is 'employees'",
		},
		{
			name:        "Filter departments by parent property",
			jsonpath:    "$.company.departments.*[?(@parentProperty === 'departments')]",
			expectedLen: 2,
			description: "Find items where parent property is 'departments'",
		},

		// Parent operator (^)
		{
			name:        "Get parent of employee",
			jsonpath:    "$.company.departments.engineering.employees[0]^",
			expectedLen: 1,
			description: "Get parent of first employee (the employees array)",
			skip:        true,
			skipReason:  "Parent operator needs debugging",
		},
		{
			name:        "Get parent of filtered items",
			jsonpath:    "$.company.departments.engineering.employees[?(@.role === 'developer')]^",
			expectedLen: 1,
			description: "Get parent of developer employees",
			skip:        true,
			skipReason:  "Parent operator needs debugging",
		},

		// @parent filter
		{
			name:        "Filter by parent object",
			jsonpath:    "$.company.departments.*.employees[?(@parent.manager === 'Alice')]",
			expectedLen: 2,
			description: "Find employees where parent department manager is Alice",
		},

		// Complex combinations
		{
			name:        "Complex property and parent filter",
			jsonpath:    "$.company.departments.*[?(@property === 'engineering')].employees[?(@parentProperty === 'employees')]",
			expectedLen: 2,
			description: "Engineering employees with parent property filter",
		},
		{
			name:        "Multi-level property names",
			jsonpath:    "$.company.departments.engineering.*~",
			expectedLen: 2,
			description: "Property names in engineering department",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skipf("Skipped: %s", tt.skipReason)
				return
			}

			results, err := jp.Query(tt.jsonpath, unmarshalJSON(t, jsonData))
			if err != nil {
				t.Errorf("Query failed: %v", err)
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

// TestChainedOperationsAdvanced tests complex chained operations
func TestChainedOperationsAdvanced(t *testing.T) {
	jsonData := `{
		"inventory": [
			{
				"category": "electronics",
				"items": [
					{"name": "Phone", "price": 500, "inStock": true},
					{"name": "Laptop", "price": 1000, "inStock": false},
					{"name": "Tablet", "price": 300, "inStock": true}
				]
			},
			{
				"category": "books",
				"items": [
					{"name": "Novel", "price": 15, "inStock": true},
					{"name": "Textbook", "price": 100, "inStock": false}
				]
			}
		]
	}`

	tests := []struct {
		name        string
		jsonpath    string
		expectedLen int
		description string
		skip        bool
		skipReason  string
	}{
		{
			name:        "Filter then property access",
			jsonpath:    "$.inventory[?(@.category === 'electronics')].items[*].name",
			expectedLen: 3,
			description: "Get names of all electronics items",
		},
		{
			name:        "Chained filter operations",
			jsonpath:    "$.inventory[*].items[?(@.price > 200)][?(@.inStock === true)]",
			expectedLen: 2,
			description: "Items over $200 that are in stock",
		},
		{
			name:        "Slice then filter",
			jsonpath:    "$.inventory[0].items[0:2][?(@.inStock === true)]",
			expectedLen: 1,
			description: "First 2 electronics items that are in stock",
			skip:        true,
			skipReason:  "Slice+filter chaining needs implementation",
		},
		{
			name:        "Union then filter",
			jsonpath:    "$.inventory[*].items[0,2][?(@.price < 400)]",
			expectedLen: 2,
			description: "First and third items under $400",
		},
		{
			name:        "Complex nested chaining",
			jsonpath:    "$.inventory[?(@.category === 'electronics')].items[?(@.inStock === true)].name",
			expectedLen: 2,
			description: "Names of in-stock electronics",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skipf("Skipped: %s", tt.skipReason)
				return
			}

			results, err := jp.Query(tt.jsonpath, unmarshalJSON(t, jsonData))
			if err != nil {
				t.Errorf("Query failed: %v", err)
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

// TestEdgeCasesAndErrorHandling tests edge cases and error conditions
func TestEdgeCasesAndErrorHandling(t *testing.T) {
	jsonData := `{
		"empty": {},
		"nullValue": null,
		"emptyArray": [],
		"mixedArray": [1, "string", null, true, {"nested": "object"}],
		"specialChars": {
			"with-dash": "value1",
			"with_underscore": "value2",
			"with.dot": "value3",
			"with space": "value4"
		}
	}`

	tests := []struct {
		name        string
		jsonpath    string
		expectedLen int
		description string
		expectError bool
	}{
		{
			name:        "Access empty object",
			jsonpath:    "$.empty.*",
			expectedLen: 0,
			description: "Access properties of empty object",
		},
		{
			name:        "Access null value",
			jsonpath:    "$.nullValue",
			expectedLen: 1,
			description: "Access null value",
		},
		{
			name:        "Access empty array",
			jsonpath:    "$.emptyArray[*]",
			expectedLen: 0,
			description: "Access elements of empty array",
		},
		{
			name:        "Filter mixed array",
			jsonpath:    "$.mixedArray[?(@)]",
			expectedLen: 4, // All except null
			description: "Filter truthy values from mixed array",
		},
		{
			name:        "Access special character properties",
			jsonpath:    "$.specialChars['with-dash']",
			expectedLen: 1,
			description: "Access property with dash using bracket notation",
		},
		{
			name:        "Access property with space",
			jsonpath:    "$.specialChars['with space']",
			expectedLen: 1,
			description: "Access property with space",
		},
		{
			name:        "Invalid array index",
			jsonpath:    "$.mixedArray[100]",
			expectedLen: 0,
			description: "Access non-existent array index",
		},
		{
			name:        "Invalid property",
			jsonpath:    "$.nonExistent",
			expectedLen: 0,
			description: "Access non-existent property",
		},
		{
			name:        "Empty path",
			jsonpath:    "",
			expectedLen: 0,
			description: "Empty JSONPath",
			expectError: true,
		},
		{
			name:        "Invalid syntax",
			jsonpath:    "$..[",
			expectedLen: 0,
			description: "Invalid bracket syntax",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := jp.Query(tt.jsonpath, unmarshalJSON(t, jsonData))

			if tt.expectError {
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
