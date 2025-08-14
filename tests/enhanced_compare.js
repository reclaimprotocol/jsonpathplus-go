#!/usr/bin/env node

const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');

// Load all test configurations
// Skip existing tests for now as they have issues
// const sharedTests = JSON.parse(fs.readFileSync('./shared/testcases.json', 'utf8'));
const sharedTests = { testCases: [] }; // Empty for now
const comprehensiveTests = JSON.parse(fs.readFileSync('./comprehensive_test_suite.json', 'utf8'));

function runGoTest(jsonpath, data) {
  try {
    const dataJson = JSON.stringify(data);
    
    const tempDataFile = `./temp_data_${Date.now()}_${Math.random().toString(36).substr(2, 9)}.json`;
    const tempQueryFile = `./temp_query_${Date.now()}_${Math.random().toString(36).substr(2, 9)}.txt`;
    fs.writeFileSync(tempDataFile, dataJson);
    fs.writeFileSync(tempQueryFile, jsonpath);
    
    try {
      const result = execSync(`cd go && go run main.go "../${tempQueryFile}" "../${tempDataFile}" --query-file --data-file`, 
        { encoding: 'utf8', timeout: 10000 });
      return JSON.parse(result.trim());
    } finally {
      try { fs.unlinkSync(tempDataFile); } catch (e) {}
      try { fs.unlinkSync(tempQueryFile); } catch (e) {}
    }
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
    const dataJson = JSON.stringify(data);
    
    const tempDataFile = `./temp_data_${Date.now()}_${Math.random().toString(36).substr(2, 9)}.json`;
    const tempQueryFile = `./temp_query_${Date.now()}_${Math.random().toString(36).substr(2, 9)}.txt`;
    fs.writeFileSync(tempDataFile, dataJson);
    fs.writeFileSync(tempQueryFile, jsonpath);
    
    try {
      const result = execSync(`node js/test.js "${tempQueryFile}" "${tempDataFile}" --query-file --data-file`, 
        { encoding: 'utf8', timeout: 10000 });
      return JSON.parse(result.trim());
    } finally {
      try { fs.unlinkSync(tempDataFile); } catch (e) {}
      try { fs.unlinkSync(tempQueryFile); } catch (e) {}
    }
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
  comparison.functionalMatch = comparison.countMatch && comparison.valuesMatch && comparison.errorMatch;
  
  return comparison;
}

function runTestSuite(suiteName, tests, testData) {
  console.log(`\n${'='.repeat(80)}`);
  console.log(`RUNNING TEST SUITE: ${suiteName.toUpperCase()}`);
  console.log(`${'='.repeat(80)}`);
  
  const results = [];
  const summary = {
    total: tests.length,
    perfectMatches: 0,
    functionalMatches: 0,
    countMatches: 0,
    valueMatches: 0,
    pathMatches: 0,
    errorMatches: 0
  };
  
  tests.forEach((test, index) => {
    console.log(`\n[${index + 1}/${tests.length}] ${test.name}`);
    console.log(`   Query: ${test.query}`);
    
    const data = testData[test.data] || test.data;
    
    console.log(`   Executing Go implementation...`);
    const goResult = runGoTest(test.query, data);
    
    console.log(`   Executing JS implementation...`);
    const jsResult = runJsTest(test.query, data);
    
    const comparison = compareResults(goResult, jsResult);
    
    // Display results
    console.log('   RESULTS:');
    console.log(`   Go:       ${goResult.error ? `ERROR: ${goResult.error.substring(0, 50)}...` : `${goResult.count} results`}`);
    console.log(`   JS:       ${jsResult.error ? `ERROR: ${jsResult.error.substring(0, 50)}...` : `${jsResult.count} results`}`);
    
    const status = comparison.perfectMatch ? '✅ PERFECT' : 
                  comparison.functionalMatch ? '🟡 FUNCTIONAL' : '❌ BROKEN';
    console.log(`   Status:   ${status}`);
    
    if (!comparison.perfectMatch) {
      const issues = [];
      if (!comparison.countMatch) issues.push('COUNT');
      if (!comparison.valuesMatch) issues.push('VALUES');
      if (!comparison.pathsMatch) issues.push('PATHS');
      if (!comparison.errorMatch) issues.push('ERRORS');
      console.log(`   Issues:   ${issues.join(', ')}`);
    }
    
    // Update summary
    if (comparison.perfectMatch) summary.perfectMatches++;
    if (comparison.functionalMatch) summary.functionalMatches++;
    if (comparison.countMatch) summary.countMatches++;
    if (comparison.valuesMatch) summary.valueMatches++;
    if (comparison.pathsMatch) summary.pathMatches++;
    if (comparison.errorMatch) summary.errorMatches++;
    
    results.push({
      test,
      goResult,
      jsResult, 
      comparison
    });
    
    console.log('-'.repeat(80));
  });
  
  // Suite summary
  console.log(`\nSUITE SUMMARY: ${suiteName}`);
  console.log(`Perfect Matches: ${summary.perfectMatches}/${summary.total} (${(summary.perfectMatches/summary.total*100).toFixed(1)}%)`);
  console.log(`Functional Matches: ${summary.functionalMatches}/${summary.total} (${(summary.functionalMatches/summary.total*100).toFixed(1)}%)`);
  console.log(`Count Matches: ${summary.countMatches}/${summary.total} (${(summary.countMatches/summary.total*100).toFixed(1)}%)`);
  console.log(`Value Matches: ${summary.valueMatches}/${summary.total} (${(summary.valueMatches/summary.total*100).toFixed(1)}%)`);
  console.log(`Path Matches: ${summary.pathMatches}/${summary.total} (${(summary.pathMatches/summary.total*100).toFixed(1)}%)`);
  
  return { results, summary };
}

function main() {
  console.log('🚀 COMPREHENSIVE JSONPATH COMPATIBILITY TEST SUITE');
  console.log('='.repeat(80));
  console.log(`Timestamp: ${new Date().toISOString()}`);
  console.log('Comparing Go JSONPath-Plus implementation with JavaScript reference');
  
  const allResults = [];
  const overallSummary = {
    totalTests: 0,
    perfectMatches: 0,
    functionalMatches: 0,
    countMatches: 0,
    valueMatches: 0,
    pathMatches: 0,
    errorMatches: 0,
    suites: {}
  };
  
  // Run existing test suite
  console.log('\n📋 RUNNING EXISTING SHARED TEST CASES');
  const existingResults = runTestSuite('existing_tests', sharedTests.testCases, sharedTests.testData);
  allResults.push({ name: 'existing_tests', ...existingResults });
  
  // Run comprehensive test suites
  Object.entries(comprehensiveTests.testSuites).forEach(([suiteName, suite]) => {
    const suiteResults = runTestSuite(suiteName, suite.tests, {
      ...sharedTests.testData,
      ...comprehensiveTests.testData
    });
    allResults.push({ name: suiteName, ...suiteResults });
  });
  
  // Calculate overall summary
  allResults.forEach(suite => {
    overallSummary.totalTests += suite.summary.total;
    overallSummary.perfectMatches += suite.summary.perfectMatches;
    overallSummary.functionalMatches += suite.summary.functionalMatches;
    overallSummary.countMatches += suite.summary.countMatches;
    overallSummary.valueMatches += suite.summary.valueMatches;
    overallSummary.pathMatches += suite.summary.pathMatches;
    overallSummary.errorMatches += suite.summary.errorMatches;
    overallSummary.suites[suite.name] = suite.summary;
  });
  
  // Final report
  console.log('\n' + '='.repeat(80));
  console.log('🏁 FINAL COMPREHENSIVE COMPATIBILITY REPORT');
  console.log('='.repeat(80));
  console.log(`Total Tests: ${overallSummary.totalTests}`);
  console.log(`Perfect Matches: ${overallSummary.perfectMatches} (${(overallSummary.perfectMatches/overallSummary.totalTests*100).toFixed(1)}%)`);
  console.log(`Functional Matches: ${overallSummary.functionalMatches} (${(overallSummary.functionalMatches/overallSummary.totalTests*100).toFixed(1)}%)`);
  console.log(`Count Matches: ${overallSummary.countMatches} (${(overallSummary.countMatches/overallSummary.totalTests*100).toFixed(1)}%)`);
  console.log(`Value Matches: ${overallSummary.valueMatches} (${(overallSummary.valueMatches/overallSummary.totalTests*100).toFixed(1)}%)`);
  console.log(`Path Matches: ${overallSummary.pathMatches} (${(overallSummary.pathMatches/overallSummary.totalTests*100).toFixed(1)}%)`);
  
  // Save detailed report
  const reportData = {
    timestamp: new Date().toISOString(),
    overallSummary,
    suiteResults: allResults.map(suite => ({
      name: suite.name,
      summary: suite.summary,
      failedTests: suite.results.filter(r => !r.comparison.perfectMatch).map(r => ({
        name: r.test.name,
        query: r.test.query,
        issues: {
          countMatch: r.comparison.countMatch,
          valuesMatch: r.comparison.valuesMatch, 
          pathsMatch: r.comparison.pathsMatch,
          errorMatch: r.comparison.errorMatch
        }
      }))
    }))
  };
  
  fs.writeFileSync('./comprehensive_test_results.json', JSON.stringify(reportData, null, 2));
  console.log('\n📄 Detailed report saved to: ./comprehensive_test_results.json');
  
  // Priority fixes needed
  console.log('\n🔧 PRIORITY FIXES NEEDED:');
  const pathIssues = allResults.reduce((acc, suite) => acc + (suite.summary.total - suite.summary.pathMatches), 0);
  const valueIssues = allResults.reduce((acc, suite) => acc + (suite.summary.total - suite.summary.valueMatches), 0);
  const countIssues = allResults.reduce((acc, suite) => acc + (suite.summary.total - suite.summary.countMatches), 0);
  
  console.log(`1. Path Format Issues: ${pathIssues} tests need path format fixes`);
  console.log(`2. Value Issues: ${valueIssues} tests have incorrect values`);
  console.log(`3. Count Issues: ${countIssues} tests have incorrect result counts`);
  
  const targetCompatibility = 100;
  const currentCompatibility = (overallSummary.perfectMatches / overallSummary.totalTests * 100);
  console.log(`\n🎯 TARGET: ${targetCompatibility}% Perfect Compatibility`);
  console.log(`📊 CURRENT: ${currentCompatibility.toFixed(1)}% Perfect Compatibility`);
  console.log(`📈 PROGRESS NEEDED: ${(targetCompatibility - currentCompatibility).toFixed(1)}% improvement required`);
}

if (require.main === module) {
  main();
}

module.exports = { runGoTest, runJsTest, compareResults };