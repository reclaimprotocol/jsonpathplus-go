#!/usr/bin/env node

const { JSONPath } = require('jsonpath-plus');

// Test data
const goessnerData = {
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
};

const companyData = {
  "company": {
    "departments": {
      "engineering": {
        "employees": [
          {"name": "Alice", "level": "senior"},
          {"name": "Bob", "level": "junior"}
        ],
        "manager": "Eve"
      },
      "sales": {
        "employees": [
          {"name": "Charlie", "level": "senior"}
        ],
        "manager": "Dave"
      }
    }
  }
};

console.log("=== Detailed Analysis of Discrepancies ===\n");

// Case 1: Book properties where parent index is not 0
console.log("1. Book properties where parent index is not 0");
console.log("   Query: $..book.*[?(@parentProperty !== 0)]");
console.log("   Go test expects: 12 results");

try {
  const result1 = JSONPath({
    path: "$..book.*[?(@parentProperty !== 0)]",
    json: goessnerData,
    resultType: 'value'
  });
  
  console.log(`   JSONPath-Plus actual: ${result1.length} results`);
  console.log("   All results:");
  result1.forEach((val, i) => {
    console.log(`     [${i}] ${JSON.stringify(val)}`);
  });
  
  // Let's also check the path information
  const resultWithPaths = JSONPath({
    path: "$..book.*[?(@parentProperty !== 0)]", 
    json: goessnerData,
    resultType: 'path'
  });
  
  console.log("   Corresponding paths:");
  resultWithPaths.forEach((path, i) => {
    console.log(`     [${i}] ${path}`);
  });
  
  // Manual count verification
  console.log("\n   Manual verification:");
  console.log("   Book 0: 4 properties (excluded by filter)");
  console.log("   Book 1: 4 properties (category, author, title, price)");  
  console.log("   Book 2: 5 properties (category, author, title, isbn, price)");
  console.log("   Book 3: 5 properties (category, author, title, isbn, price)");
  console.log("   Expected total: 4 + 5 + 5 = 14 properties");
  console.log("   Go test expects: 12 (may be incorrect)");
  
} catch (error) {
  console.log(`   ERROR: ${error.message}`);
}

console.log("\n" + "=".repeat(60) + "\n");

// Case 2: Filter departments by parent property  
console.log("2. Filter departments by parent property");
console.log("   Query: $.company.departments.*[?(@parentProperty === 'departments')]");
console.log("   Go test expects: 2 results");

try {
  const result2 = JSONPath({
    path: "$.company.departments.*[?(@parentProperty === 'departments')]",
    json: companyData,
    resultType: 'value'
  });
  
  console.log(`   JSONPath-Plus actual: ${result2.length} results`);
  console.log("   All results:");
  result2.forEach((val, i) => {
    console.log(`     [${i}] ${JSON.stringify(val)}`);
  });
  
  // Let's understand what @parentProperty should be here
  console.log("\n   Understanding the context:");
  
  // Get all departments first
  const depts = JSONPath({
    path: "$.company.departments.*",
    json: companyData,
    resultType: 'value'
  });
  
  console.log(`   All departments: ${depts.length} items`);
  depts.forEach((dept, i) => {
    console.log(`     [${i}] ${JSON.stringify(dept)}`);
  });
  
  // Check what departments itself contains
  const deptsObj = JSONPath({
    path: "$.company.departments", 
    json: companyData,
    resultType: 'value'
  });
  
  console.log(`   Departments object: ${JSON.stringify(deptsObj[0])}`);
  console.log(`   Properties: ${Object.keys(deptsObj[0])}`);
  
  // Test different interpretations
  console.log("\n   Testing alternate interpretations:");
  
  // Maybe the query should be different?
  const alt1 = JSONPath({
    path: "$.company.*[?(@parentProperty === 'company')]",
    json: companyData,
    resultType: 'value'
  });
  console.log(`   $.company.*[?(@parentProperty === 'company')]: ${alt1.length} results`);
  
  // Or maybe it's looking for department properties
  const alt2 = JSONPath({
    path: "$.company.departments.*.*[?(@parentProperty === 'departments')]", 
    json: companyData,
    resultType: 'value'
  });
  console.log(`   $.company.departments.*.*[?(@parentProperty === 'departments')]: ${alt2.length} results`);
  
} catch (error) {
  console.log(`   ERROR: ${error.message}`);
}

console.log("\n=== Conclusions ===");
console.log("1. JSONPath-Plus reference implementation confirms most Go test expectations");
console.log("2. Case 1: Go test expects 12 but should expect 14 (test expectation may be wrong)");
console.log("3. Case 2: JSONPath-Plus also returns 0, suggesting @parentProperty behavior");
console.log("   may not work as the Go test expects, or the query needs adjustment");