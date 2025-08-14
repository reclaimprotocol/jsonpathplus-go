#!/usr/bin/env node

/**
 * Comprehensive JSONPath Library Comparison Tool
 * 
 * This script compares the Go JSONPath implementation against the reference
 * JSONPath-Plus JavaScript implementation using identical inputs and generates
 * a detailed analysis report.
 */

const { JSONPath } = require('jsonpath-plus');
const fs = require('fs');
const { execSync } = require('child_process');

// Load test data registry
const testDataRegistry = JSON.parse(fs.readFileSync('./test_data_registry.json', 'utf8'));

// Comprehensive test cases covering all major functionality
const testCases = [
  // === GOESSNER SPECIFICATION TESTS ===
  {
    name: "Authors of all books",
    jsonpath: "$.store.book[*].author",
    data: "goessner_spec_data",
    category: "goessner_spec",
    description: "Basic property access in arrays"
  },
  {
    name: "All authors", 
    jsonpath: "$..author",
    data: "goessner_spec_data",
    category: "goessner_spec",
    description: "Recursive descent for specific property"
  },
  {
    name: "All elements beneath root",
    jsonpath: "$..*",
    data: "goessner_spec_data", 
    category: "goessner_spec",
    description: "Complete recursive descent"
  },
  {
    name: "All prices in store",
    jsonpath: "$.store..price",
    data: "goessner_spec_data",
    category: "goessner_spec", 
    description: "Recursive descent for price property"
  },
  {
    name: "Third book",
    jsonpath: "$..book[2]",
    data: "goessner_spec_data",
    category: "goessner_spec",
    description: "Array index access with recursive descent"
  },
  {
    name: "Last book",
    jsonpath: "$..book[-1]",
    data: "goessner_spec_data",
    category: "goessner_spec",
    description: "Negative array indexing"
  },
  {
    name: "First two books",
    jsonpath: "$..book[0,1]",
    data: "goessner_spec_data", 
    category: "goessner_spec",
    description: "Multiple array indices (union)"
  },
  {
    name: "Books with ISBN",
    jsonpath: "$..book[?(@.isbn)]",
    data: "goessner_spec_data",
    category: "goessner_spec",
    description: "Filter by property existence"
  },
  {
    name: "Books cheaper than 10",
    jsonpath: "$..book[?(@.price<10)]", 
    data: "goessner_spec_data",
    category: "goessner_spec",
    description: "Numeric comparison filter"
  },
  {
    name: "All books except first",
    jsonpath: "$..book[1:]",
    data: "goessner_spec_data",
    category: "goessner_spec", 
    description: "Array slice notation"
  },
  {
    name: "Price properties not equal to 8.95",
    jsonpath: "$..*[?(@property === 'price' && @ !== 8.95)]",
    data: "goessner_spec_data",
    category: "goessner_spec",
    description: "Property name filter with value condition"
  },
  {
    name: "Book children except category",
    jsonpath: "$..book.*[?(@property !== \"category\")]",
    data: "goessner_spec_data",
    category: "goessner_spec",
    description: "Property name exclusion filter"
  },
  {
    name: "Books not at index 0", 
    jsonpath: "$..book[?(@property !== 0)]",
    data: "goessner_spec_data",
    category: "goessner_spec",
    description: "Array index filter"
  },
  {
    name: "Book properties where parent index is not 0",
    jsonpath: "$..book.*[?(@parentProperty !== 0)]",
    data: "goessner_spec_data",
    category: "goessner_spec",
    description: "Parent property filter"
  },

  // === JSONPATH-PLUS SPECIFIC FEATURES ===
  {
    name: "Filter by property name",
    jsonpath: "$.company.departments.*[?(@property === 'engineering')]",
    data: "company_data",
    category: "jsonpath_plus_features",
    description: "Property name matching"
  },
  {
    name: "Filter by parent property",
    jsonpath: "$.company.departments.engineering.employees[?(@parentProperty === 'engineering')]", 
    data: "company_data",
    category: "jsonpath_plus_features",
    description: "Parent property context"
  },
  {
    name: "Filter departments by parent property", 
    jsonpath: "$.company.departments.*[?(@parentProperty === 'departments')]",
    data: "company_data",
    category: "jsonpath_plus_features",
    description: "Object property parent filtering"
  },
  {
    name: "Parent filter - simple",
    jsonpath: "$.store.book[?(@parent.bicycle)]",
    data: "goessner_spec_data",
    category: "jsonpath_plus_features", 
    description: "Parent object property existence"
  },
  {
    name: "Parent filter - property value",
    jsonpath: "$.store.book[?(@parent.bicycle.color === 'red')]",
    data: "goessner_spec_data",
    category: "jsonpath_plus_features",
    description: "Parent object nested property value"
  },

  // === ADVANCED FILTER SCENARIOS ===
  {
    name: "Premium customers with shipped orders",
    jsonpath: "$.orders[?(@.customer.type === 'premium' && @.status === 'shipped')]",
    data: "orders_data",
    category: "advanced_filters",
    description: "Logical AND in filters"
  },
  {
    name: "Orders over $100",
    jsonpath: "$.orders[?(@.total > 100)]", 
    data: "orders_data",
    category: "advanced_filters",
    description: "Numeric comparison filter"
  },
  {
    name: "Orders with laptop products",
    jsonpath: "$.orders[?(@.items[*].product === 'laptop')]",
    data: "orders_data",
    category: "advanced_filters",
    description: "Array wildcard in filter expressions"
  },

  // === FUNCTION PREDICATES ===
  {
    name: "String contains function",
    jsonpath: "$.store.book[?(@.title.contains('Sword'))]",
    data: "goessner_spec_data", 
    category: "function_predicates",
    description: "String contains method"
  },
  {
    name: "String startsWith function", 
    jsonpath: "$.store.book[?(@.title.startsWith('Moby'))]",
    data: "goessner_spec_data",
    category: "function_predicates",
    description: "String startsWith method"
  },
  {
    name: "Regex match function",
    jsonpath: "$.store.book[?(@.title.match(/.*Century.*/))]",
    data: "goessner_spec_data",
    category: "function_predicates",
    description: "Regular expression matching"
  },

  // === EDGE CASES ===
  {
    name: "Empty array filter",
    jsonpath: "$[?(@.nonexistent)]",
    data: "goessner_spec_data",
    category: "edge_cases", 
    description: "Filter on non-existent property"
  },
  {
    name: "Null value access",
    jsonpath: "$.store.book[0].nonexistent",
    data: "goessner_spec_data",
    category: "edge_cases",
    description: "Access non-existent property"
  }
];

/**
 * Execute JSONPath query using Go implementation
 */
function executeGoQuery(jsonpath, data) {
  try {
    // Write test data to temporary file
    const dataStr = JSON.stringify(testDataRegistry[data]);
    fs.writeFileSync('/tmp/test_data.json', dataStr);
    
    // Create Go test program
    const goProgram = `
package main

import (
	"encoding/json"
	"fmt"
	"log"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
	"os"
)

func main() {
	data, err := os.ReadFile("/tmp/test_data.json")
	if err != nil {
		log.Fatal(err)
	}
	
	results, err := jp.Query("${jsonpath}", string(data))
	if err != nil {
		fmt.Printf("ERROR: %v\\n", err)
		return
	}
	
	values := make([]interface{}, len(results))
	for i, result := range results {
		values[i] = result.Value
	}
	
	output, err := json.Marshal(map[string]interface{}{
		"count": len(results),
		"values": values,
		"paths": func() []string {
			paths := make([]string, len(results))
			for i, result := range results {
				paths[i] = result.Path
			}
			return paths
		}(),
	})
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println(string(output))
}`;

    fs.writeFileSync('/tmp/go_test.go', goProgram);
    
    // Execute Go program
    const result = execSync('cd /tmp && go mod init temp && go mod edit -require=github.com/reclaimprotocol/jsonpathplus-go@v0.0.0 && go mod edit -replace=github.com/reclaimprotocol/jsonpathplus-go@v0.0.0=/Users/abdul/Desktop/code/cc_exp/jsonpathplus-go && go run go_test.go', 
      { timeout: 10000, encoding: 'utf8' });
    
    return JSON.parse(result.trim());
  } catch (error) {
    return {
      error: error.message,
      count: 0,
      values: [],
      paths: []
    };
  }
}

/**
 * Execute JSONPath query using JSONPath-Plus reference implementation
 */
function executeJSONPathPlusQuery(jsonpath, data) {
  try {
    const testData = testDataRegistry[data];
    
    const values = JSONPath({
      path: jsonpath,
      json: testData,
      resultType: 'value'
    });
    
    const paths = JSONPath({
      path: jsonpath, 
      json: testData,
      resultType: 'path'
    });
    
    return {
      count: values.length,
      values: values,
      paths: paths
    };
  } catch (error) {
    return {
      error: error.message,
      count: 0,
      values: [],
      paths: []
    };
  }
}

/**
 * Compare two result sets
 */
function compareResults(goResult, jsResult) {
  const comparison = {
    countMatch: goResult.count === jsResult.count,
    valuesMatch: JSON.stringify(goResult.values) === JSON.stringify(jsResult.values),
    pathsMatch: JSON.stringify(goResult.paths) === JSON.stringify(jsResult.paths),
    goError: !!goResult.error,
    jsError: !!jsResult.error,
    errorMatch: (!!goResult.error) === (!!jsResult.error)
  };
  
  comparison.perfectMatch = comparison.countMatch && comparison.valuesMatch && comparison.pathsMatch && comparison.errorMatch;
  
  return comparison;
}

/**
 * Main test execution
 */
function runComprehensiveComparison() {
  console.log('='.repeat(80));
  console.log('COMPREHENSIVE JSONPATH LIBRARY COMPARISON');
  console.log('='.repeat(80));
  console.log(`Timestamp: ${new Date().toISOString()}`);
  console.log(`Go Implementation: /Users/abdul/Desktop/code/cc_exp/jsonpathplus-go`);
  console.log(`Reference Implementation: JSONPath-Plus v${require('jsonpath-plus/package.json').version}`);
  console.log(`Total Test Cases: ${testCases.length}`);
  console.log('='.repeat(80));
  console.log('');

  const results = [];
  const summary = {
    total: testCases.length,
    perfectMatches: 0,
    countMatches: 0,
    valueMatches: 0,
    pathMatches: 0,
    errorMatches: 0,
    categories: {}
  };

  testCases.forEach((testCase, index) => {
    console.log(`[${index + 1}/${testCases.length}] ${testCase.name}`);
    console.log(`   Query: ${testCase.jsonpath}`);
    console.log(`   Data: ${testCase.data}`);
    console.log(`   Category: ${testCase.category}`);
    console.log(`   Description: ${testCase.description}`);
    
    // Execute both implementations
    const goResult = executeGoQuery(testCase.jsonpath, testCase.data);
    const jsResult = executeJSONPathPlusQuery(testCase.jsonpath, testCase.data);
    
    // Compare results
    const comparison = compareResults(goResult, jsResult);
    
    // Display results
    console.log('');
    console.log('   GO IMPLEMENTATION:');
    if (goResult.error) {
      console.log(`     ❌ ERROR: ${goResult.error}`);
    } else {
      console.log(`     Count: ${goResult.count}`);
      console.log(`     Values: ${JSON.stringify(goResult.values).substring(0, 100)}${goResult.values.length > 2 ? '...' : ''}`);
    }
    
    console.log('');
    console.log('   JSONPATH-PLUS REFERENCE:');
    if (jsResult.error) {
      console.log(`     ❌ ERROR: ${jsResult.error}`);
    } else {
      console.log(`     Count: ${jsResult.count}`);
      console.log(`     Values: ${JSON.stringify(jsResult.values).substring(0, 100)}${jsResult.values.length > 2 ? '...' : ''}`);
    }
    
    console.log('');
    console.log('   COMPARISON:');
    console.log(`     Perfect Match: ${comparison.perfectMatch ? '✅ YES' : '❌ NO'}`);
    console.log(`     Count Match: ${comparison.countMatch ? '✅' : '❌'} (${goResult.count} vs ${jsResult.count})`);
    console.log(`     Values Match: ${comparison.valuesMatch ? '✅' : '❌'}`);
    console.log(`     Paths Match: ${comparison.pathsMatch ? '✅' : '❌'}`);
    console.log(`     Error Match: ${comparison.errorMatch ? '✅' : '❌'}`);
    
    // Update summary
    if (comparison.perfectMatch) summary.perfectMatches++;
    if (comparison.countMatch) summary.countMatches++;
    if (comparison.valuesMatch) summary.valueMatches++;
    if (comparison.pathsMatch) summary.pathMatches++;
    if (comparison.errorMatch) summary.errorMatches++;
    
    if (!summary.categories[testCase.category]) {
      summary.categories[testCase.category] = { total: 0, matches: 0 };
    }
    summary.categories[testCase.category].total++;
    if (comparison.perfectMatch) {
      summary.categories[testCase.category].matches++;
    }
    
    // Store detailed result
    results.push({
      testCase,
      goResult,
      jsResult,
      comparison
    });
    
    console.log('');
    console.log('-'.repeat(80));
    console.log('');
  });

  // Generate summary report
  console.log('');
  console.log('='.repeat(80));
  console.log('SUMMARY REPORT');
  console.log('='.repeat(80));
  console.log(`Total Tests: ${summary.total}`);
  console.log(`Perfect Matches: ${summary.perfectMatches} (${(summary.perfectMatches/summary.total*100).toFixed(1)}%)`);
  console.log(`Count Matches: ${summary.countMatches} (${(summary.countMatches/summary.total*100).toFixed(1)}%)`);
  console.log(`Value Matches: ${summary.valueMatches} (${(summary.valueMatches/summary.total*100).toFixed(1)}%)`);
  console.log(`Path Matches: ${summary.pathMatches} (${(summary.pathMatches/summary.total*100).toFixed(1)}%)`);
  console.log(`Error Handling Matches: ${summary.errorMatches} (${(summary.errorMatches/summary.total*100).toFixed(1)}%)`);
  console.log('');
  
  console.log('BY CATEGORY:');
  Object.entries(summary.categories).forEach(([category, stats]) => {
    const percentage = (stats.matches / stats.total * 100).toFixed(1);
    console.log(`  ${category}: ${stats.matches}/${stats.total} (${percentage}%)`);
  });
  
  console.log('');
  console.log('FAILED TESTS:');
  const failures = results.filter(r => !r.comparison.perfectMatch);
  if (failures.length === 0) {
    console.log('  🎉 All tests passed perfectly!');
  } else {
    failures.forEach((failure, index) => {
      console.log(`  ${index + 1}. ${failure.testCase.name}`);
      console.log(`     Query: ${failure.testCase.jsonpath}`);
      console.log(`     Go: ${failure.goResult.count} results${failure.goResult.error ? ' (ERROR)' : ''}`);
      console.log(`     JS: ${failure.jsResult.count} results${failure.jsResult.error ? ' (ERROR)' : ''}`);
      console.log(`     Issues: ${!failure.comparison.countMatch ? 'COUNT ' : ''}${!failure.comparison.valuesMatch ? 'VALUES ' : ''}${!failure.comparison.pathsMatch ? 'PATHS ' : ''}${!failure.comparison.errorMatch ? 'ERRORS' : ''}`);
    });
  }
  
  console.log('');
  console.log('='.repeat(80));
  console.log(`OVERALL COMPATIBILITY: ${(summary.perfectMatches/summary.total*100).toFixed(1)}%`);
  console.log('='.repeat(80));

  // Save detailed report
  const reportData = {
    timestamp: new Date().toISOString(),
    summary,
    results,
    testCases
  };
  
  fs.writeFileSync('./comparison_report.json', JSON.stringify(reportData, null, 2));
  console.log('');
  console.log('📄 Detailed report saved to: comparison_report.json');
  
  return reportData;
}

// Execute the comprehensive comparison
if (require.main === module) {
  runComprehensiveComparison();
}