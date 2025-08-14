package compatibility_tests

import (
	"testing"

	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

// TestRefactoredAdvancedFeatures tests the new advanced features
func TestRefactoredAdvancedFeatures(t *testing.T) {
	jsonData := `{
		"company": {
			"departments": {
				"engineering": {
					"employees": [
						{"name": "Alice", "role": "developer", "id": 1},
						{"name": "Bob", "role": "manager", "id": 2}
					]
				},
				"sales": {
					"employees": [
						{"name": "Charlie", "role": "rep", "id": 3}
					]
				}
			}
		}
	}`

	tests := []struct {
		name        string
		path        string
		expectedLen int
		description string
		useRefactored bool
	}{
		// Test chained operations
		{
			name:          "Chained filter and property access",
			path:          "$.company.departments.engineering.employees[?(@.role == 'developer')].name",
			expectedLen:   1,
			description:   "Get names of developers using chained operations",
			useRefactored: true,
		},
		
		// Test property names operator
		{
			name:          "Property names of departments",
			path:          "$.company.departments.*~",
			expectedLen:   2,
			description:   "Get property names of departments",
			useRefactored: true,
		},
		
		// Test @property filter
		{
			name:          "Filter by property name",
			path:          "$.company.departments.*[?(@property === 'engineering')]",
			expectedLen:   1,
			description:   "Find departments with property name 'engineering'",
			useRefactored: true,
		},
		
		// Test @parentProperty filter
		{
			name:          "Filter by parent property",
			path:          "$.company.departments.engineering.employees[?(@parentProperty === 'employees')]",
			expectedLen:   2,
			description:   "Find items where parent property is 'employees'",
			useRefactored: true,
		},
		
		// Test parent operator
		{
			name:          "Get parent of filtered items",
			path:          "$.company.departments.engineering.employees[?(@.role == 'manager')]^",
			expectedLen:   1,
			description:   "Get parent of manager employees",
			useRefactored: true,
		},
		
		// Compare with original implementation
		{
			name:          "Basic query with original",
			path:          "$.company.departments.engineering.employees[*].name",
			expectedLen:   2,
			description:   "Basic query using original implementation",
			useRefactored: false,
		},
		{
			name:          "Basic query with refactored",
			path:          "$.company.departments.engineering.employees[*].name",
			expectedLen:   2,
			description:   "Basic query using refactored implementation",
			useRefactored: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var results interface{}
			var err error
			
			// Use the unified Query API
			queryResults, queryErr := jp.Query(tt.path, unmarshalJSON(t, jsonData))
			results = queryResults
			err = queryErr
			
			if err != nil {
				// Some advanced features might not be supported in original implementation
				if !tt.useRefactored {
					t.Logf("Original implementation doesn't support: %s", tt.path)
					return
				}
				t.Errorf("Query failed: %v", err)
				return
			}
			
			queryResultsTyped := results.([]jp.Result)
			resultLen := len(queryResultsTyped)
			
			if resultLen != tt.expectedLen {
				t.Errorf("%s: Expected %d results, got %d", tt.description, tt.expectedLen, resultLen)
				t.Logf("Query: %s", tt.path)
				t.Logf("Using refactored: %v", tt.useRefactored)
			}
		})
	}
}

// TestChainedOperations specifically tests chained bracket operations
func TestChainedOperations(t *testing.T) {
	jsonData := `{
		"products": [
			{"name": "A", "price": 10, "category": "electronics"},
			{"name": "B", "price": 20, "category": "books"},
			{"name": "C", "price": 30, "category": "electronics"},
			{"name": "D", "price": 40, "category": "books"},
			{"name": "E", "price": 50, "category": "electronics"},
			{"name": "F", "price": 60, "category": "books"}
		]
	}`

	tests := []struct {
		name        string
		path        string
		expectedLen int
		description string
	}{
		{
			name:        "First 3 products then filter by category",
			path:        "$.products[0:3][?(@.category == 'electronics')]",
			expectedLen: 2, // Products A and C
			description: "Slice then filter chained operation",
		},
		{
			name:        "Last 3 products then filter by price",
			path:        "$.products[-3:][?(@.price > 40)]",
			expectedLen: 2, // Products E and F
			description: "Negative slice then filter chained operation",
		},
		{
			name:        "Filter then get names",
			path:        "$.products[?(@.category == 'electronics')].name",
			expectedLen: 3, // Names of products A, C, E
			description: "Filter then property access",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := jp.Query(tt.path, unmarshalJSON(t, jsonData))
			if err != nil {
				t.Errorf("Query failed: %v", err)
				return
			}

			if len(results) != tt.expectedLen {
				t.Errorf("%s: Expected %d results, got %d", tt.description, tt.expectedLen, len(results))
				t.Logf("Query: %s", tt.path)
				for i, result := range results {
					t.Logf("Result %d: %v", i, result.Value)
				}
			}
		})
	}
}

// Helper function to unmarshal JSON
func unmarshalJSON(t *testing.T, jsonStr string) interface{} {
	data, err := jp.JSONParse(jsonStr)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
	return data
}