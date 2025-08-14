#!/usr/bin/env node

// Cross-validation script against JSONPath-Plus reference implementation
// Run: npm install jsonpath-plus && node validate_against_reference.js

const { JSONPath } = require('jsonpath-plus');

// Test data from Goessner specification
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

// JSONPath-Plus test data
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

// Test cases that are currently failing
const testCases = [
  {
    name: "All elements beneath root",
    query: "$..*",
    data: goessnerData,
    expectedCount: 27,
    description: "Should include all nested elements"
  },
  {
    name: "Price properties not equal to 8.95", 
    query: "$..*[?(@property === 'price' && @ !== 8.95)]",
    data: goessnerData,
    expectedCount: 4,
    description: "Should find all price properties except 8.95"
  },
  {
    name: "Book children except category",
    query: "$..book.*[?(@property !== \"category\")]", 
    data: goessnerData,
    expectedCount: 14,
    description: "All book properties except category"
  },
  {
    name: "Books not at index 0",
    query: "$..book[?(@property !== 0)]",
    data: goessnerData, 
    expectedCount: 3,
    description: "Books at indices 1, 2, 3"
  },
  {
    name: "Store grandchildren where parent is not book",
    query: "$.store.*[?(@parentProperty !== \"book\")]",
    data: goessnerData,
    expectedCount: 2,
    description: "Store children where parent property != 'book'"
  },
  {
    name: "Book properties where parent index is not 0", 
    query: "$..book.*[?(@parentProperty !== 0)]",
    data: goessnerData,
    expectedCount: 12,
    description: "Properties of books at non-zero indices"
  },
  {
    name: "Filter departments by parent property",
    query: "$.company.departments.*[?(@parentProperty === 'departments')]",
    data: companyData,
    expectedCount: 2, 
    description: "Department objects where parent property is 'departments'"
  }
];

console.log("=== JSONPath-Plus Reference Implementation Validation ===\n");

let totalTests = 0;
let matchingResults = 0;
let differentResults = 0;

testCases.forEach((testCase, index) => {
  totalTests++;
  
  console.log(`${index + 1}. ${testCase.name}`);
  console.log(`   Query: ${testCase.query}`);
  console.log(`   Expected Count: ${testCase.expectedCount}`);
  
  try {
    const result = JSONPath({
      path: testCase.query,
      json: testCase.data,
      resultType: 'value'
    });
    
    const actualCount = result.length;
    console.log(`   Actual Count: ${actualCount}`);
    
    if (actualCount === testCase.expectedCount) {
      console.log(`   ✅ MATCHES expected count`);
      matchingResults++;
    } else {
      console.log(`   ❌ DIFFERENT from expected count`);
      differentResults++;
      
      // Show first few results for debugging
      console.log(`   First 5 results:`);
      result.slice(0, 5).forEach((val, i) => {
        console.log(`     [${i}] ${JSON.stringify(val)}`);
      });
      if (result.length > 5) {
        console.log(`     ... and ${result.length - 5} more`);
      }
    }
    
    console.log(`   Description: ${testCase.description}`);
    
  } catch (error) {
    console.log(`   ❌ ERROR: ${error.message}`);
    differentResults++;
  }
  
  console.log("");
});

console.log("=== Summary ===");
console.log(`Total test cases: ${totalTests}`);
console.log(`Matching expected counts: ${matchingResults}`);
console.log(`Different from expected: ${differentResults}`);
console.log(`Match rate: ${((matchingResults / totalTests) * 100).toFixed(1)}%`);

if (differentResults > 0) {
  console.log("\n=== Analysis ===");
  console.log("Some test expectations may be incorrect, or JSONPath-Plus");
  console.log("behavior may differ from the expected values in the Go tests.");
  console.log("This validation helps determine if the Go implementation");
  console.log("should match JSONPath-Plus exactly or if test expectations");
  console.log("need to be updated.");
}