#!/usr/bin/env node

const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');

// Load shared test cases
const testData = JSON.parse(fs.readFileSync('./shared/testcases.json', 'utf8'));

function runGoTest(jsonpath, data) {
  try {
    const dataJson = JSON.stringify(testData.testData[data]);
    const result = execSync(`cd go && go run main.go "${jsonpath}" '${dataJson}'`, 
      { encoding: 'utf8', timeout: 10000 });
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

function runJsTest(jsonpath, data) {
  try {
    const dataJson = JSON.stringify(testData.testData[data]);
    const result = execSync(`node js/test.js "${jsonpath}" '${dataJson}'`, 
      { encoding: 'utf8', timeout: 10000 });
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

function main() {
  // Allow filtering by category or test name
  const filterArg = process.argv[2];
  let testCases = testData.testCases;
  
  if (filterArg) {
    testCases = testCases.filter(tc => 
      tc.category.includes(filterArg) || 
      tc.name.toLowerCase().includes(filterArg.toLowerCase())
    );
  }

  console.log('='.repeat(80));
  console.log('JSONPATH LIBRARY COMPARISON TEST');
  console.log('='.repeat(80));
  console.log(`Timestamp: ${new Date().toISOString()}`);
  console.log(`Test Cases: ${testCases.length}${filterArg ? ` (filtered by: ${filterArg})` : ''}`);
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
    console.log(`   Data: ${testCase.data}`);
    
    // Execute both implementations
    console.log(`   Executing Go implementation...`);
    const goResult = runGoTest(testCase.jsonpath, testCase.data);
    
    console.log(`   Executing JS implementation...`);
    const jsResult = runJsTest(testCase.jsonpath, testCase.data);
    
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
  
  fs.writeFileSync('./test_results.json', JSON.stringify(reportData, null, 2));
  console.log('');
  console.log('📄 Detailed report saved to: ./test_results.json');
}

if (require.main === module) {
  main();
}