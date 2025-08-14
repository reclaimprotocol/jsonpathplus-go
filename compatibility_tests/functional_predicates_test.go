package compatibility_tests

import (
	"testing"

	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

// TestFunctionalPredicates tests function-based filtering capabilities
func TestFunctionalPredicates(t *testing.T) {
	jsonData := `{
		"users": [
			{"name": "Alice Johnson", "email": "alice@example.com", "age": 30},
			{"name": "Bob Smith", "email": "bob@test.org", "age": 25},
			{"name": "Charlie Brown", "email": "charlie@example.com", "age": 35},
			{"name": "Diana Prince", "email": "diana@hero.gov", "age": 28}
		],
		"products": [
			{"title": "JavaScript Guide", "category": "programming"},
			{"title": "Python Cookbook", "category": "programming"},
			{"title": "Design Patterns", "category": "architecture"},
			{"title": "Clean Code", "category": "programming"}
		],
		"tags": ["frontend", "backend", "mobile", "web", "api"],
		"descriptions": [
			"This is a great product",
			"Excellent quality item",
			"Good value for money",
			"Outstanding performance"
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
		// String matching functions
		{
			name:        "Match by string contains",
			jsonpath:    "$.users[?(@.email.contains('example'))]",
			expectedLen: 2,
			description: "Find users with 'example' in email",
			skip:        false,
			skipReason:  "",
		},
		{
			name:        "Match by string starts with",
			jsonpath:    "$.users[?(@.name.startsWith('A'))]",
			expectedLen: 1,
			description: "Find users whose name starts with 'A'",
			skip:        false,
			skipReason:  "",
		},
		{
			name:        "Match by string ends with",
			jsonpath:    "$.users[?(@.email.endsWith('.com'))]",
			expectedLen: 2,
			description: "Find users with .com email addresses",
			skip:        false,
			skipReason:  "",
		},
		{
			name:        "Match by regex pattern",
			jsonpath:    "$.users[?(@.email.match(/.*@example\\.com$/))]",
			expectedLen: 2,
			description: "Find users with example.com email using regex",
			skip:        false,
			skipReason:  "",
		},

		// Array functions
		{
			name:        "Check array length",
			jsonpath:    "$.tags[?(@.length === 5)]",
			expectedLen: 0,
			description: "Check if tags array has 5 elements (applied to individual elements)",
			skip:        true,
			skipReason:  "Array length function not implemented",
		},
		{
			name:        "Array includes check",
			jsonpath:    "$.products[?(@.category.includes('prog'))]",
			expectedLen: 0,
			description: "Products in categories containing 'prog'",
			skip:        true,
			skipReason:  "String includes function not implemented",
		},

		// Type checking functions
		{
			name:        "Check string type",
			jsonpath:    "$.users[*].name[?(@.typeof() === 'string')]",
			expectedLen: 4,
			description: "Find all string-type names",
			skip:        false,
			skipReason:  "",
		},
		{
			name:        "Check number type",
			jsonpath:    "$.users[*].age[?(@.typeof() === 'number')]",
			expectedLen: 4,
			description: "Find all number-type ages",
			skip:        false,
			skipReason:  "",
		},

		// Mathematical functions
		{
			name:        "Math operations in filter",
			jsonpath:    "$.users[?(@.age.floor() > 25)]",
			expectedLen: 3,
			description: "Users with floor(age) > 25",
			skip:        false,
			skipReason:  "",
		},
		{
			name:        "Round function",
			jsonpath:    "$.users[?(@.age.round() === 25)]",
			expectedLen: 1,
			description: "Users with rounded age of 25",
			skip:        false,
			skipReason:  "",
		},

		// Currently working basic filters (for comparison)
		{
			name:        "Basic string equality",
			jsonpath:    "$.users[?(@.name === 'Alice Johnson')]",
			expectedLen: 1,
			description: "Find user by exact name match",
		},
		{
			name:        "Basic numeric comparison",
			jsonpath:    "$.users[?(@.age > 30)]",
			expectedLen: 1,
			description: "Find users older than 30",
		},
		{
			name:        "Basic string inequality",
			jsonpath:    "$.products[?(@.category !== 'programming')]",
			expectedLen: 1,
			description: "Find non-programming products",
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

// TestLogicalOperators tests complex logical operations in filters
func TestLogicalOperators(t *testing.T) {
	jsonData := `{
		"employees": [
			{"name": "Alice", "age": 30, "department": "engineering", "salary": 75000},
			{"name": "Bob", "age": 25, "department": "sales", "salary": 50000},
			{"name": "Charlie", "age": 35, "department": "engineering", "salary": 85000},
			{"name": "Diana", "age": 28, "department": "marketing", "salary": 60000},
			{"name": "Eve", "age": 32, "department": "engineering", "salary": 80000}
		]
	}`

	tests := []struct {
		name        string
		jsonpath    string
		expectedLen int
		description string
	}{
		{
			name:        "AND operation",
			jsonpath:    "$.employees[?(@.department === 'engineering' && @.age > 30)]",
			expectedLen: 2,
			description: "Engineers older than 30",
		},
		{
			name:        "OR operation",
			jsonpath:    "$.employees[?(@.department === 'sales' || @.department === 'marketing')]",
			expectedLen: 2,
			description: "Sales or marketing employees",
		},
		{
			name:        "Complex AND/OR combination",
			jsonpath:    "$.employees[?(@.age > 30 && (@.department === 'engineering' || @.salary > 70000))]",
			expectedLen: 2,
			description: "Older employees in engineering or with high salary",
		},
		{
			name:        "Negation with AND",
			jsonpath:    "$.employees[?(@.department !== 'engineering' && @.age < 30)]",
			expectedLen: 2,
			description: "Young non-engineers",
		},
		{
			name:        "Multiple conditions",
			jsonpath:    "$.employees[?(@.age >= 25 && @.age <= 32 && @.salary >= 60000)]",
			expectedLen: 3,
			description: "Employees aged 25-32 with good salary",
		},
		{
			name:        "Parentheses grouping",
			jsonpath:    "$.employees[?(@.age > 30 && (@.department === 'engineering' || @.department === 'marketing'))]",
			expectedLen: 2,
			description: "Older employees in engineering or marketing",
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
					t.Logf("Result %d: %v", i, result.Value)
				}
			}
		})
	}
}

// TestComparisonOperators tests all comparison operators
func TestComparisonOperators(t *testing.T) {
	jsonData := `{
		"numbers": [1, 5, 10, 15, 20, 25, 30],
		"items": [
			{"id": 1, "price": 9.99, "rating": 4.5},
			{"id": 2, "price": 19.99, "rating": 3.8},
			{"id": 3, "price": 29.99, "rating": 4.2},
			{"id": 4, "price": 9.99, "rating": 4.7},
			{"id": 5, "price": 39.99, "rating": 4.0}
		]
	}`

	tests := []struct {
		name        string
		jsonpath    string
		expectedLen int
		description string
	}{
		{
			name:        "Equals operator",
			jsonpath:    "$.items[?(@.price === 9.99)]",
			expectedLen: 2,
			description: "Items with exact price 9.99",
		},
		{
			name:        "Not equals operator",
			jsonpath:    "$.items[?(@.price !== 9.99)]",
			expectedLen: 3,
			description: "Items not priced at 9.99",
		},
		{
			name:        "Greater than operator",
			jsonpath:    "$.items[?(@.rating > 4.0)]",
			expectedLen: 3,
			description: "Items with rating above 4.0",
		},
		{
			name:        "Greater than or equal operator",
			jsonpath:    "$.items[?(@.rating >= 4.0)]",
			expectedLen: 4,
			description: "Items with rating 4.0 or above",
		},
		{
			name:        "Less than operator",
			jsonpath:    "$.items[?(@.price < 20.00)]",
			expectedLen: 3,
			description: "Items under $20",
		},
		{
			name:        "Less than or equal operator",
			jsonpath:    "$.items[?(@.price <= 20.00)]",
			expectedLen: 3,
			description: "Items $20 or under",
		},
		{
			name:        "Number range with array",
			jsonpath:    "$.numbers[?(@ >= 10 && @ <= 20)]",
			expectedLen: 3,
			description: "Numbers between 10 and 20",
		},
		{
			name:        "Exact equality check",
			jsonpath:    "$.numbers[?(@ === 15)]",
			expectedLen: 1,
			description: "Exact match for number 15",
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
					t.Logf("Result %d: %v", i, result.Value)
				}
			}
		})
	}
}

// TestAdvancedFilterScenarios tests complex real-world filtering scenarios
func TestAdvancedFilterScenarios(t *testing.T) {
	jsonData := `{
		"orders": [
			{
				"id": "ORD001",
				"customer": {"name": "Alice", "type": "premium"},
				"items": [
					{"product": "laptop", "price": 999.99, "quantity": 1},
					{"product": "mouse", "price": 29.99, "quantity": 2}
				],
				"status": "shipped",
				"total": 1059.97
			},
			{
				"id": "ORD002", 
				"customer": {"name": "Bob", "type": "regular"},
				"items": [
					{"product": "keyboard", "price": 79.99, "quantity": 1}
				],
				"status": "pending",
				"total": 79.99
			},
			{
				"id": "ORD003",
				"customer": {"name": "Charlie", "type": "premium"},
				"items": [
					{"product": "monitor", "price": 299.99, "quantity": 1},
					{"product": "cable", "price": 19.99, "quantity": 3}
				],
				"status": "shipped",
				"total": 359.97
			}
		]
	}`

	tests := []struct {
		name        string
		jsonpath    string
		expectedLen int
		description string
	}{
		{
			name:        "Premium customers with shipped orders",
			jsonpath:    "$.orders[?(@.customer.type === 'premium' && @.status === 'shipped')]",
			expectedLen: 2,
			description: "Shipped orders from premium customers",
		},
		{
			name:        "Orders over $100",
			jsonpath:    "$.orders[?(@.total > 100)]",
			expectedLen: 2,
			description: "High-value orders",
		},
		{
			name:        "Orders with laptop products",
			jsonpath:    "$.orders[?(@.items[*].product === 'laptop')]",
			expectedLen: 1,
			description: "Orders containing laptops",
		},
		{
			name:        "Get all order items",
			jsonpath:    "$.orders[*].items[*]",
			expectedLen: 5,
			description: "All individual order items",
		},
		{
			name:        "Expensive items in orders",
			jsonpath:    "$.orders[*].items[?(@.price > 100)]",
			expectedLen: 2,
			description: "Order items over $100",
		},
		// Note: length property test would go here when implemented
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
					t.Logf("Result %d: %v", i, result.Value)
				}
			}
		})
	}
}
