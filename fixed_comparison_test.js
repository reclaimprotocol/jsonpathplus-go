#!/usr/bin/env node

/**
 * Fixed Comprehensive JSONPath Library Comparison Tool
 */

const { JSONPath } = require('jsonpath-plus');
const fs = require('fs');
const { execSync } = require('child_process');
const path = require('path');

// Load test data registry
const testDataRegistry = JSON.parse(fs.readFileSync('./test_data_registry.json', 'utf8'));

// Selected critical test cases for comparison
const testCases = [
  // === BASIC FUNCTIONALITY ===
  {
    name: "Authors of all books",
    jsonpath: "$.store.book[*].author",
    data: "goessner_spec_data",
    category: "basic",
    description: "Basic property access in arrays"
  },
  {
    name: "All authors", 
    jsonpath: "$..author",
    data: "goessner_spec_data",
    category: "recursive_descent",
    description: "Recursive descent for specific property"
  },
  {
    name: "All elements beneath root",
    jsonpath: "$..*",
    data: "goessner_spec_data", 
    category: "recursive_descent",
    description: "Complete recursive descent"
  },
  {
    name: "Third book",
    jsonpath: "$..book[2]",
    data: "goessner_spec_data",
    category: "array_access",
    description: "Array index access with recursive descent"
  },
  {
    name: "Books with ISBN",
    jsonpath: "$..book[?(@.isbn)]",
    data: "goessner_spec_data",
    category: "filters",
    description: "Filter by property existence"
  },
  {
    name: "Books cheaper than 10",
    jsonpath: "$..book[?(@.price<10)]", 
    data: "goessner_spec_data",
    category: "filters",
    description: "Numeric comparison filter"
  },
  {
    name: "Price properties not equal to 8.95",
    jsonpath: "$..*[?(@property === 'price' && @ !== 8.95)]",
    data: "goessner_spec_data",
    category: "property_filters",
    description: "Property name filter with value condition"
  },
  {
    name: "Books not at index 0", 
    jsonpath: "$..book[?(@property !== 0)]",
    data: "goessner_spec_data",
    category: "property_filters",
    description: "Array index filter"
  },
  {
    name: "Parent filter - simple",
    jsonpath: "$.store.book[?(@parent.bicycle)]",
    data: "goessner_spec_data",
    category: "parent_filters", 
    description: "Parent object property existence"
  },
  {
    name: "Orders with laptop products",
    jsonpath: "$.orders[?(@.items[*].product === 'laptop')]",
    data: "orders_data",
    category: "array_wildcards",
    description: "Array wildcard in filter expressions"
  }
];

/**
 * Execute JSONPath query using Go implementation
 */
function executeGoQuery(jsonpath, data, testIndex) {
  try {
    // Create unique temporary directory for this test
    const tempDir = `/tmp/jsonpath_test_${testIndex}_${Date.now()}`;
    fs.mkdirSync(tempDir, { recursive: true });
    
    // Write test data to file
    const dataStr = JSON.stringify(testDataRegistry[data]);
    fs.writeFileSync(path.join(tempDir, 'test_data.json'), dataStr);
    
    // Create Go test program with proper escaping
    const escapedJsonPath = jsonpath.replace(/"/g, '\\"');
    const goProgram = `package main

import (
	"encoding/json"
	"fmt"
	"log"
	jp "github.com/reclaimprotocol/jsonpathplus-go"
	"os"
)

func main() {
	data, err := os.ReadFile("test_data.json")
	if err != nil {
		log.Fatal(err)
	}
	
	results, err := jp.Query("${escapedJsonPath}", string(data))
	if err != nil {
		fmt.Printf("ERROR: %v\\n", err)
		return
	}
	
	values := make([]interface{}, len(results))
	paths := make([]string, len(results))
	for i, result := range results {
		values[i] = result.Value
		paths[i] = result.Path
	}
	
	output, err := json.Marshal(map[string]interface{}{
		"count": len(results),
		"values": values,
		"paths": paths,
	})
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println(string(output))
}`;

    fs.writeFileSync(path.join(tempDir, 'main.go'), goProgram);
    
    // Initialize Go module and execute
    execSync(`cd ${tempDir} && go mod init test_${testIndex}`, { timeout: 5000 });
    execSync(`cd ${tempDir} && go mod edit -require=github.com/reclaimprotocol/jsonpathplus-go@v0.0.0`, { timeout: 5000 });
    execSync(`cd ${tempDir} && go mod edit -replace=github.com/reclaimprotocol/jsonpathplus-go@v0.0.0=/Users/abdul/Desktop/code/cc_exp/jsonpathplus-go`, { timeout: 5000 });
    
    const result = execSync(`cd ${tempDir} && go run main.go`, 
      { timeout: 10000, encoding: 'utf8' });
    
    // Clean up
    execSync(`rm -rf ${tempDir}`);
    
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
function runComparisonTest() {
  console.log('='.repeat(80));
  console.log('JSONPATH LIBRARY COMPARISON TEST');
  console.log('='.repeat(80));
  console.log(`Timestamp: ${new Date().toISOString()}`);
  console.log(`Go Implementation: /Users/abdul/Desktop/code/cc_exp/jsonpathplus-go`);
  console.log(`Reference Implementation: JSONPath-Plus v${require('jsonpath-plus/package.json').version}`);
  console.log(`Test Cases: ${testCases.length}`);
  console.log('='.repeat(80));
  console.log('');

  const results = [];
  const summary = {
    total: testCases.length,
    perfectMatches: 0,
    countMatches: 0,
    valueMatches: 0,
    categories: {}
  };

  testCases.forEach((testCase, index) => {
    console.log(`[${index + 1}/${testCases.length}] ${testCase.name}`);
    console.log(`   Query: ${testCase.jsonpath}`);
    console.log(`   Category: ${testCase.category}`);
    console.log(`   Description: ${testCase.description}`);
    
    // Execute both implementations
    console.log(`   Executing Go implementation...`);
    const goResult = executeGoQuery(testCase.jsonpath, testCase.data, index);
    
    console.log(`   Executing JSONPath-Plus...`);
    const jsResult = executeJSONPathPlusQuery(testCase.jsonpath, testCase.data);
    
    // Compare results
    const comparison = compareResults(goResult, jsResult);
    
    // Display results
    console.log('');
    console.log('   RESULTS:');
    console.log(`   Go:       ${goResult.error ? `ERROR: ${goResult.error.substring(0, 50)}...` : `${goResult.count} results`}`);
    console.log(`   JS:       ${jsResult.error ? `ERROR: ${jsResult.error}` : `${jsResult.count} results`}`);
    console.log(`   Match:    ${comparison.perfectMatch ? '✅ PERFECT' : '❌ DIFFERENT'}`);
    
    if (!comparison.perfectMatch) {
      console.log(`   Issues:   ${!comparison.countMatch ? 'COUNT ' : ''}${!comparison.valuesMatch ? 'VALUES ' : ''}${!comparison.errorMatch ? 'ERRORS' : ''}`);
      
      if (!goResult.error && !jsResult.error && goResult.count > 0 && jsResult.count > 0) {
        console.log(`   Go Values:  ${JSON.stringify(goResult.values.slice(0, 3))}${goResult.values.length > 3 ? '...' : ''}`);
        console.log(`   JS Values:  ${JSON.stringify(jsResult.values.slice(0, 3))}${jsResult.values.length > 3 ? '...' : ''}`);
      }
    }
    
    // Update summary
    if (comparison.perfectMatch) summary.perfectMatches++;
    if (comparison.countMatch) summary.countMatches++;
    if (comparison.valuesMatch) summary.valueMatches++;
    
    if (!summary.categories[testCase.category]) {
      summary.categories[testCase.category] = { total: 0, matches: 0 };
    }
    summary.categories[testCase.category].total++;
    if (comparison.perfectMatch) {
      summary.categories[testCase.category].matches++;
    }
    
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
  console.log('='.repeat(80));
  console.log('FINAL COMPARISON REPORT');
  console.log('='.repeat(80));
  console.log(`Total Tests: ${summary.total}`);
  console.log(`Perfect Matches: ${summary.perfectMatches} (${(summary.perfectMatches/summary.total*100).toFixed(1)}%)`);
  console.log(`Count Matches: ${summary.countMatches} (${(summary.countMatches/summary.total*100).toFixed(1)}%)`);
  console.log(`Value Matches: ${summary.valueMatches} (${(summary.valueMatches/summary.total*100).toFixed(1)}%)`);
  console.log('');
  
  console.log('CATEGORY BREAKDOWN:');
  Object.entries(summary.categories).forEach(([category, stats]) => {
    const percentage = (stats.matches / stats.total * 100).toFixed(1);
    console.log(`  ${category.padEnd(20)}: ${stats.matches}/${stats.total} (${percentage}%)`);
  });
  
  console.log('');
  console.log('DETAILED FAILURES:');
  const failures = results.filter(r => !r.comparison.perfectMatch);
  if (failures.length === 0) {
    console.log('  🎉 All tests passed perfectly!');
  } else {
    failures.forEach((failure, index) => {
      console.log(`  ${index + 1}. ${failure.testCase.name}`);
      console.log(`     Query: ${failure.testCase.jsonpath}`);
      
      if (failure.goResult.error) {
        console.log(`     Go Error: ${failure.goResult.error.substring(0, 80)}...`);
      } else {
        console.log(`     Go: ${failure.goResult.count} results`);
      }
      
      if (failure.jsResult.error) {
        console.log(`     JS Error: ${failure.jsResult.error}`);
      } else {
        console.log(`     JS: ${failure.jsResult.count} results`);
      }
      
      console.log('');
    });
  }
  
  console.log('='.repeat(80));
  console.log(`OVERALL COMPATIBILITY RATE: ${(summary.perfectMatches/summary.total*100).toFixed(1)}%`);
  console.log('='.repeat(80));

  // Save detailed report
  const reportData = {
    timestamp: new Date().toISOString(),
    summary,
    results: results.map(r => ({
      testCase: r.testCase,
      goResult: {
        count: r.goResult.count,
        error: r.goResult.error,
        hasValues: r.goResult.values.length > 0
      },
      jsResult: {
        count: r.jsResult.count,
        error: r.jsResult.error,
        hasValues: r.jsResult.values.length > 0
      },
      comparison: r.comparison
    }))
  };
  
  fs.writeFileSync('./detailed_comparison_report.json', JSON.stringify(reportData, null, 2));
  console.log('');
  console.log('📄 Detailed report saved to: detailed_comparison_report.json');
  
  return reportData;
}

// Execute the comparison
if (require.main === module) {
  runComparisonTest();
}