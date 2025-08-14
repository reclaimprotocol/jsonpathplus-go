#!/usr/bin/env node

const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');

// Load shared test cases
const testData = JSON.parse(fs.readFileSync('./shared/testcases.json', 'utf8'));

function runGoTest(jsonpath, data) {
  try {
    // Load data from individual files in data directory
    let dataJson;
    if (testData.testData[data]) {
      dataJson = JSON.stringify(testData.testData[data]);
    } else {
      // Try to load from data directory
      try {
        const dataFile = `./data/${data}.json`;
        dataJson = fs.readFileSync(dataFile, 'utf8');
      } catch (e) {
        throw new Error(`Data not found: ${data}`);
      }
    }
    
    // Write both query and data to temporary files to avoid shell escaping issues
    const tempDataFile = `./temp_data_${Date.now()}_${Math.random().toString(36).substr(2, 9)}.json`;
    const tempQueryFile = `./temp_query_${Date.now()}_${Math.random().toString(36).substr(2, 9)}.txt`;
    fs.writeFileSync(tempDataFile, dataJson);
    fs.writeFileSync(tempQueryFile, jsonpath);
    
    try {
      const result = execSync(`cd go && go run main.go "../${tempQueryFile}" "../${tempDataFile}" --query-file --data-file`, 
        { encoding: 'utf8', timeout: 10000 });
      return JSON.parse(result.trim());
    } finally {
      // Clean up temp files
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
    // Load data from individual files in data directory
    let dataJson;
    if (testData.testData[data]) {
      dataJson = JSON.stringify(testData.testData[data]);
    } else {
      // Try to load from data directory
      try {
        const dataFile = `./data/${data}.json`;
        dataJson = fs.readFileSync(dataFile, 'utf8');
      } catch (e) {
        throw new Error(`Data not found: ${data}`);
      }
    }
    
    // Write both query and data to temporary files to avoid shell escaping issues
    const tempDataFile = `./temp_data_${Date.now()}_${Math.random().toString(36).substr(2, 9)}.json`;
    const tempQueryFile = `./temp_query_${Date.now()}_${Math.random().toString(36).substr(2, 9)}.txt`;
    fs.writeFileSync(tempDataFile, dataJson);
    fs.writeFileSync(tempQueryFile, jsonpath);
    
    try {
      const result = execSync(`node js/test.js "${tempQueryFile}" "${tempDataFile}" --query-file --data-file`, 
        { encoding: 'utf8', timeout: 10000 });
      return JSON.parse(result.trim());
    } finally {
      // Clean up temp files
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
  
  // Detailed Analysis Summary
  console.log('');
  console.log('🔍 COMPATIBILITY ANALYSIS:');
  console.log('');
  
  // Categorize issues
  const zeroResultIssues = results.filter(r => !r.comparison.perfectMatch && r.goResult.count === 0 && r.jsResult.count > 0);
  const countMismatches = results.filter(r => !r.comparison.countMatch);
  const valueMismatches = results.filter(r => r.comparison.countMatch && !r.comparison.valuesMatch);
  const errorMismatches = results.filter(r => !r.comparison.errorMatch);
  
  console.log(`📊 Issue Breakdown:`);
  console.log(`   • Zero Results (Go=0, JS>0): ${zeroResultIssues.length} tests`);
  console.log(`   • Count Mismatches: ${countMismatches.length} tests`);
  console.log(`   • Value Mismatches: ${valueMismatches.length} tests`);  
  console.log(`   • Error Handling Mismatches: ${errorMismatches.length} tests`);
  console.log('');
  
  // Priority Issues (Zero results from Go)
  if (zeroResultIssues.length > 0) {
    console.log('🚨 HIGH PRIORITY - Zero Result Issues:');
    zeroResultIssues.slice(0, 5).forEach((issue, i) => {
      console.log(`   ${i+1}. ${issue.testCase.name}`);
      console.log(`      Query: ${issue.testCase.jsonpath}`);
      console.log(`      Category: ${issue.testCase.category}`);
    });
    if (zeroResultIssues.length > 5) {
      console.log(`   ... and ${zeroResultIssues.length - 5} more`);
    }
    console.log('');
  }
  
  // Working Categories
  const workingCategories = Object.entries(summary.categories).filter(([cat, stats]) => stats.matches === stats.total);
  if (workingCategories.length > 0) {
    console.log('✅ FULLY WORKING CATEGORIES:');
    workingCategories.forEach(([category, stats]) => {
      console.log(`   • ${category}: ${stats.matches}/${stats.total}`);
    });
    console.log('');
  }
  
  // Problematic Categories  
  const problematicCategories = Object.entries(summary.categories)
    .filter(([cat, stats]) => stats.matches / stats.total < 0.5)
    .sort(([,a], [,b]) => (a.matches/a.total) - (b.matches/b.total));
    
  if (problematicCategories.length > 0) {
    console.log('❌ NEEDS ATTENTION:');
    problematicCategories.forEach(([category, stats]) => {
      const percentage = (stats.matches / stats.total * 100).toFixed(1);
      console.log(`   • ${category}: ${stats.matches}/${stats.total} (${percentage}%)`);
    });
    console.log('');
  }
  
  // Next Steps
  console.log('🎯 RECOMMENDED NEXT STEPS:');
  if (zeroResultIssues.length > 0) {
    const topCategories = [...new Set(zeroResultIssues.map(r => r.testCase.category))].slice(0, 3);
    console.log(`   1. Debug filter evaluation for: ${topCategories.join(', ')}`);
  }
  if (errorMismatches.length > 0) {
    console.log(`   2. Implement missing error handling patterns`);
  }
  if (valueMismatches.length > 0) {
    console.log(`   3. Fix remaining value serialization differences`);
  }
  console.log('');

  // Save detailed report
  const reportData = {
    timestamp: new Date().toISOString(),
    summary,
    analysis: {
      zeroResultIssues: zeroResultIssues.length,
      countMismatches: countMismatches.length,
      valueMismatches: valueMismatches.length,
      errorMismatches: errorMismatches.length,
      workingCategories: workingCategories.map(([cat, stats]) => cat),
      problematicCategories: problematicCategories.map(([cat, stats]) => ({ category: cat, success_rate: (stats.matches/stats.total*100).toFixed(1) + '%' }))
    },
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
  console.log('📄 Detailed report saved to: ./test_results.json');
}

if (require.main === module) {
  main();
}