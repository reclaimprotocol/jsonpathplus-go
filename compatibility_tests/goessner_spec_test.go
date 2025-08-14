package compatibility_tests

import (
	"testing"

	jp "github.com/reclaimprotocol/jsonpathplus-go"
)

// TestGoessnerSpecification tests all JSONPath expressions from the original Goessner specification
// Reference: http://goessner.net/articles/JsonPath/
func TestGoessnerSpecification(t *testing.T) {
	// Standard test data from Goessner's original article
	jsonData := `{
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
					"title": "Sword of Honour",
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
		}
	}`

	tests := []struct {
		name        string
		xpath       string
		jsonpath    string
		expectedLen int
		description string
		skip        bool
		skipReason  string
	}{
		{
			name:        "Authors of all books",
			xpath:       "/store/book/author",
			jsonpath:    "$.store.book[*].author",
			expectedLen: 4,
			description: "The authors of all books in the store",
		},
		{
			name:        "All authors anywhere",
			xpath:       "//author",
			jsonpath:    "$..author",
			expectedLen: 4,
			description: "All authors",
		},
		{
			name:        "All things in store",
			xpath:       "/store/*",
			jsonpath:    "$.store.*",
			expectedLen: 2,
			description: "All things in store (books array and bicycle object)",
		},
		{
			name:        "All prices in store",
			xpath:       "/store//price",
			jsonpath:    "$.store..price",
			expectedLen: 5,
			description: "The price of everything in the store",
		},
		{
			name:        "Third book",
			xpath:       "//book[3]",
			jsonpath:    "$..book[2]",
			expectedLen: 1,
			description: "The third book (0-indexed as [2])",
		},
		{
			name:        "Last book",
			xpath:       "//book[last()]",
			jsonpath:    "$..book[-1:]",
			expectedLen: 1,
			description: "The last book in order",
		},
		{
			name:        "First two books",
			xpath:       "//book[position()<3]",
			jsonpath:    "$..book[:2]",
			expectedLen: 2,
			description: "The first two books",
		},
		{
			name:        "First two books (union syntax)",
			xpath:       "//book[position()<3]",
			jsonpath:    "$..book[0,1]",
			expectedLen: 2,
			description: "The first two books using union syntax",
		},
		{
			name:        "Books with ISBN",
			xpath:       "//book[isbn]",
			jsonpath:    "$..book[?(@.isbn)]",
			expectedLen: 2,
			description: "Filter all books with an ISBN number",
		},
		{
			name:        "Books cheaper than 10",
			xpath:       "//book[price<10]",
			jsonpath:    "$..book[?(@.price<10)]",
			expectedLen: 2,
			description: "Filter all books cheaper than 10",
		},
		{
			name:        "Root object",
			xpath:       "/",
			jsonpath:    "$",
			expectedLen: 1,
			description: "The root of the JSON object",
		},
		{
			name:        "All elements beneath root",
			xpath:       "$..*",
			jsonpath:    "$..*",
			expectedLen: 27, // Approximate count of all nested elements
			description: "All members of JSON structure beneath root",
		},
		{
			name:        "All parent components",
			xpath:       "//*",
			jsonpath:    "$..",
			expectedLen: 28, // Including root
			description: "All parent components including root",
		},
		// Advanced JSONPath-Plus features
		{
			name:        "Property names of store",
			xpath:       "/store/*/name()",
			jsonpath:    "$.store.*~",
			expectedLen: 2,
			description: "Property names of store sub-object ('book' and 'bicycle')",
		},
		{
			name:        "Price properties not equal to 8.95",
			xpath:       "//*[name() = 'price' and . != 8.95]",
			jsonpath:    "$..*[?(@property === 'price' && @ !== 8.95)]",
			expectedLen: 4,
			description: "All price property values not equal to 8.95",
		},
		{
			name:        "Parent of items with price > 19",
			xpath:       "//*[price>19]/..",
			jsonpath:    "$..[?(@.price>19)]^",
			expectedLen: 2,
			description: "Parent of items with price > 19 (store and book array)",
		},
		{
			name:        "Books where parent has red bicycle",
			xpath:       "//book[parent::*/bicycle/color = \"red\"]/category",
			jsonpath:    "$..book[?(@parent.bicycle && @parent.bicycle.color === \"red\")].category",
			expectedLen: 4,
			description: "Categories of books where parent has red bicycle",
		},
		{
			name:        "Book children except category",
			xpath:       "//book/*[name() != 'category']",
			jsonpath:    "$..book.*[?(@property !== \"category\")]",
			expectedLen: 14, // All book properties except categories
			description: "All children of books except category",
		},
		{
			name:        "Books not at index 0",
			xpath:       "//book[position() != 1]",
			jsonpath:    "$..book[?(@property !== 0)]",
			expectedLen: 3,
			description: "All books except the first one",
		},
		{
			name:        "Store grandchildren where parent is not book",
			xpath:       "/store/*/*[name(parent::*) != 'book']",
			jsonpath:    "$.store.*[?(@parentProperty !== \"book\")]",
			expectedLen: 2,
			description: "Grandchildren of store where parent property is not 'book'",
		},
		{
			name:        "Book properties where parent index is not 0",
			xpath:       "//book[count(preceding-sibling::*) != 0]/*/text()",
			jsonpath:    "$..book.*[?(@parentProperty !== 0)]",
			expectedLen: 12, // Properties of books at indices 1, 2, 3
			description: "Properties of books where parent index is not 0",
		},
		{
			name:        "Books with price equal to third book price",
			xpath:       "//book[price = /store/book[3]/price]",
			jsonpath:    "$..book[?(@.price === @root.store.book[2].price)]",
			expectedLen: 1,
			description: "Books with price equal to third book price",
			skip:        true,
			skipReason:  "@root feature not yet implemented",
		},
		{
			name:        "Fiction books with regex",
			xpath:       "//book/*[name() = 'category' and matches(., 'tion')]",
			jsonpath:    "$..book.*[?(@property === \"category\" && @.match(/tion$/i))]",
			expectedLen: 3,
			description: "Categories ending in 'tion' (case insensitive)",
			skip:        true,
			skipReason:  "Regex matching not yet implemented",
		},
		{
			name:        "Books with ISBN property using regex",
			xpath:       "//book/*[matches(name(), 'isbn')]/parent::*",
			jsonpath:    "$..book.*[?(@property.match(/isbn$/i))]^",
			expectedLen: 2,
			description: "Books with properties matching 'isbn' regex",
			skip:        true,
			skipReason:  "Regex matching and parent selector combination not yet implemented",
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
				t.Logf("XPath: %s", tt.xpath)
				t.Logf("JSONPath: %s", tt.jsonpath)
				for i, result := range results {
					t.Logf("Result %d: %v", i, result.Value)
				}
			}
		})
	}
}

// TestBasicJSONPathOperations tests fundamental JSONPath operations
func TestBasicJSONPathOperations(t *testing.T) {
	jsonData := `{
		"store": {
			"book": [
				{"title": "Book 1", "price": 10.50},
				{"title": "Book 2", "price": 15.75}
			],
			"bicycle": {"color": "red", "price": 19.95}
		}
	}`

	tests := []struct {
		name        string
		jsonpath    string
		expectedLen int
		description string
	}{
		{
			name:        "Root access",
			jsonpath:    "$",
			expectedLen: 1,
			description: "Access root object",
		},
		{
			name:        "Direct property access",
			jsonpath:    "$.store",
			expectedLen: 1,
			description: "Access store property",
		},
		{
			name:        "Array access by index",
			jsonpath:    "$.store.book[0]",
			expectedLen: 1,
			description: "Access first book",
		},
		{
			name:        "Array access all elements",
			jsonpath:    "$.store.book[*]",
			expectedLen: 2,
			description: "Access all books",
		},
		{
			name:        "Nested property access",
			jsonpath:    "$.store.book[*].title",
			expectedLen: 2,
			description: "Access all book titles",
		},
		{
			name:        "Wildcard property access",
			jsonpath:    "$.store.*",
			expectedLen: 2,
			description: "Access all store properties",
		},
		{
			name:        "Recursive descent",
			jsonpath:    "$..price",
			expectedLen: 3,
			description: "Find all price properties recursively",
		},
		{
			name:        "Array slice",
			jsonpath:    "$.store.book[0:1]",
			expectedLen: 1,
			description: "Get books from index 0 to 1",
		},
		{
			name:        "Union selection",
			jsonpath:    "$.store.book[0,1]",
			expectedLen: 2,
			description: "Get books at indices 0 and 1",
		},
		{
			name:        "Filter expression",
			jsonpath:    "$.store.book[?(@.price > 12)]",
			expectedLen: 1,
			description: "Filter books by price greater than 12",
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

// TestStringIndexFinding tests string index functionality
func TestStringIndexFinding(t *testing.T) {
	jsonData := `{
		"users": ["alice", "bob", "charlie", "alice"],
		"tags": ["go", "json", "path", "go"],
		"mixed": [1, "test", 3, "test", 5]
	}`

	tests := []struct {
		name        string
		jsonpath    string
		expectedLen int
		description string
	}{
		{
			name:        "Find exact string match",
			jsonpath:    "$.users[?(@==='alice')]",
			expectedLen: 2,
			description: "Find all occurrences of 'alice'",
		},
		{
			name:        "Find string in mixed array",
			jsonpath:    "$.mixed[?(@==='test')]",
			expectedLen: 2,
			description: "Find string 'test' in mixed type array",
		},
		{
			name:        "Find first occurrence index",
			jsonpath:    "$.users[0]",
			expectedLen: 1,
			description: "Get first user",
		},
		{
			name:        "Find last occurrence",
			jsonpath:    "$.users[-1]",
			expectedLen: 1,
			description: "Get last user",
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

// TestAdvancedFilterExpressions tests complex filter scenarios
func TestAdvancedFilterExpressions(t *testing.T) {
	jsonData := `{
		"products": [
			{"id": 1, "name": "Widget", "price": 10.0, "inStock": true},
			{"id": 2, "name": "Gadget", "price": 25.0, "inStock": false},
			{"id": 3, "name": "Tool", "price": 15.0, "inStock": true},
			{"id": 4, "name": "Device", "price": 30.0, "inStock": false}
		]
	}`

	tests := []struct {
		name        string
		jsonpath    string
		expectedLen int
		description string
	}{
		{
			name:        "Filter by boolean property",
			jsonpath:    "$.products[?(@.inStock)]",
			expectedLen: 2,
			description: "Products that are in stock",
		},
		{
			name:        "Filter by negated boolean",
			jsonpath:    "$.products[?(!@.inStock)]",
			expectedLen: 2,
			description: "Products that are not in stock",
		},
		{
			name:        "Filter by price range",
			jsonpath:    "$.products[?(@.price >= 15 && @.price <= 25)]",
			expectedLen: 2,
			description: "Products with price between 15 and 25",
		},
		{
			name:        "Filter by string comparison",
			jsonpath:    "$.products[?(@.name === 'Widget')]",
			expectedLen: 1,
			description: "Product named 'Widget'",
		},
		{
			name:        "Filter by multiple conditions",
			jsonpath:    "$.products[?(@.price > 20 && @.inStock === false)]",
			expectedLen: 2,
			description: "Expensive products not in stock",
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

// TestSliceOperations tests various array slicing scenarios
func TestSliceOperations(t *testing.T) {
	jsonData := `{
		"numbers": [0, 1, 2, 3, 4, 5, 6, 7, 8, 9],
		"letters": ["a", "b", "c", "d", "e"]
	}`

	tests := []struct {
		name        string
		jsonpath    string
		expectedLen int
		description string
	}{
		{
			name:        "Slice first three",
			jsonpath:    "$.numbers[0:3]",
			expectedLen: 3,
			description: "First three numbers",
		},
		{
			name:        "Slice last three",
			jsonpath:    "$.numbers[-3:]",
			expectedLen: 3,
			description: "Last three numbers",
		},
		{
			name:        "Slice with step",
			jsonpath:    "$.numbers[::2]",
			expectedLen: 5,
			description: "Every second number",
		},
		{
			name:        "Slice middle range",
			jsonpath:    "$.numbers[2:7]",
			expectedLen: 5,
			description: "Numbers from index 2 to 6",
		},
		{
			name:        "Reverse slice",
			jsonpath:    "$.numbers[::-1]",
			expectedLen: 10,
			description: "All numbers in reverse order",
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
