# JSONPath-Plus Compatibility Tests

This directory contains comprehensive test suites that verify compatibility with the [JSONPath-Plus](https://github.com/JSONPath-Plus/JSONPath) JavaScript library.

## Test Coverage

### Core Functionality
- **basic_api_test.go** - Basic JSONPath API operations
- **array_test.go** - Array operations and indexing
- **filter_test.go** - Filter expressions and conditions
- **recursive_test.go** - Recursive descent operations (`..`)
- **slice_test.go** - Array slicing with various patterns
- **union_test.go** - Union operator for multiple selections

## Running Tests

```bash
# Run all compatibility tests
go test ./compatibility_tests/...

# Run specific test file
go test ./compatibility_tests/basic_api_test.go

# Run with verbose output
go test -v ./compatibility_tests/...

# Run specific test case
go test -v ./compatibility_tests/ -run TestBasicAPI
```

## Test Data Structure

Most tests use the standard JSONPath test data structure:

```json
{
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
      }
    ],
    "bicycle": {
      "color": "red",
      "price": 19.95
    }
  }
}
```

## Test Categories

### 1. Basic API Tests (`basic_api_test.go`)
- Root element access (`$`)
- Property access (`$.store.book[*].author`)
- Array indexing (`$.store.book[0]`, `$.store.book[-1]`)
- Recursive descent (`$..price`)

### 2. Array Tests (`array_test.go`)
- Wildcard selection (`$[*]`)
- Array slicing (`$[0:2]`, `$[-2:]`)
- Nested array access (`$.matrix[*][1]`)
- Empty array handling

### 3. Filter Tests (`filter_test.go`)
- Comparison operators (`@.price < 10`)
- Existence checks (`@.isbn`)
- String equality (`@.category == 'fiction'`)
- Complex conditions (`@.category == 'electronics' && @.inStock == true`)

### 4. Recursive Tests (`recursive_test.go`)
- Simple recursive descent (`$..author`)
- Deep nesting (`$..details..publisher`)
- Recursive with filters (`$..book[*].title`)

### 5. Slice Tests (`slice_test.go`)
- Basic slicing (`[0:3]`, `[2:5]`)
- Negative indices (`[-3:]`, `[:-1]`)
- Step values (`[::2]`, `[1::3]`)
- Reverse slicing (`[::-1]`)

### 6. Union Tests (`union_test.go`)
- Property unions (`['title','author']`)
- Index unions (`[0,1]`)
- Mixed selections (`$.store['book','bicycle']`)
- Special character keys (`['weird-key','normal_key']`)

## Compatibility Notes

These tests are designed to match the behavior of JSONPath-Plus JavaScript library:

1. **Result Order**: Results maintain the order they appear in the JSON structure
2. **Data Types**: Numbers are preserved as `float64`, strings as `string`
3. **Array Handling**: Empty results return empty slice, not `nil`
4. **Path Strings**: Generated paths match JSONPath-Plus format

## Adding New Tests

When adding new test cases:

1. Follow the existing test structure
2. Use meaningful test names and descriptions
3. Include both positive and negative test cases
4. Test edge cases (empty arrays, missing properties, etc.)
5. Validate both result count and actual values

## Expected Compatibility

These tests verify compatibility with JSONPath-Plus version 7.x features:

- ✅ Basic JSONPath syntax
- ✅ Array operations
- ✅ Filter expressions
- ✅ Recursive descent
- ✅ Union operators
- ✅ Array slicing
- ✅ Negative indexing
- ⚠️  Script expressions (limited support)
- ⚠️  Parent references (limited support)

## Performance

All tests should complete within reasonable time:
- Individual test cases: < 1ms
- Complete test suite: < 1s

Run benchmarks with:
```bash
go test -bench=. ./compatibility_tests/...
```