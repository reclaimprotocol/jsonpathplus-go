# JSONPath-Plus Go - Test Suite

This directory contains the organized test infrastructure for comparing the Go JSONPath implementation against the JavaScript JSONPath-Plus reference implementation.

## Structure

```
tests/
├── shared/
│   └── testcases.json          # Common test data and cases
├── go/
│   ├── go.mod                  # Go module for testing
│   └── main.go                 # Go test binary
├── js/
│   └── test.js                 # JavaScript test binary  
├── compare.js                  # Main comparison script
├── single.sh                   # Single query testing
└── test_results.json          # Latest test results
```

## Usage

### Run Full Comparison Test
```bash
cd tests
node compare.js
```

### Filter Tests by Category
```bash
node compare.js debug           # Run only debug tests
node compare.js property        # Run property filter tests  
node compare.js recursive       # Run recursive descent tests
```

### Test Single Query
```bash
./single.sh '$..book[*]' 'goessner_spec_data'
./single.sh '$..book[?(@property !== 0)]' 'simple_books'
```

### Test Go Implementation Only
```bash
cd go
go run main.go '$..book[*]' '{"store":{"book":[{"title":"A"}]}}'
```

### Test JS Implementation Only  
```bash
node js/test.js '$..book[*]' '{"store":{"book":[{"title":"A"}]}}'
```

## Test Data

All test data is centralized in `shared/testcases.json` with:
- **testData**: JSON datasets (goessner_spec_data, company_data, etc.)
- **testCases**: Array of test cases with query, data reference, and metadata

## Current Status

Last run compatibility: **10.0% perfect matches**
- Count matches: 80% (correct result counts)
- Value matches: 30% (exact value matching) 
- Main issues: Property ordering and array index filtering

See `test_results.json` for detailed results.