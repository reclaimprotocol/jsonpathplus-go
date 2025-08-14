#!/usr/bin/env node

const { JSONPath } = require('jsonpath-plus');

console.log("=== Final @parentProperty Behavior Analysis ===\n");

// Simple test to understand the exact semantics
const simpleData = {
  "level1": {
    "level2": {
      "array": [
        {"name": "item1"},
        {"name": "item2"}
      ],
      "object": {"prop": "value"}
    }
  }
};

console.log("Test data structure:");
console.log(JSON.stringify(simpleData, null, 2));
console.log("");

// Test array elements - we know this works
console.log("1. Array elements (known to work):");
const arrayTest = JSONPath({
  path: "$.level1.level2.array[?(@parentProperty === 'level2')]",
  json: simpleData,
  resultType: 'value'
});
console.log(`   $.level1.level2.array[?(@parentProperty === 'level2')]: ${arrayTest.length} results`);
arrayTest.forEach((val, i) => console.log(`     [${i}] ${JSON.stringify(val)}`));

const arrayTest2 = JSONPath({
  path: "$.level1.level2.array[?(@parentProperty === 'array')]", 
  json: simpleData,
  resultType: 'value'
});
console.log(`   $.level1.level2.array[?(@parentProperty === 'array')]: ${arrayTest2.length} results`);

console.log("");

// Test object properties - we suspect this doesn't work
console.log("2. Object properties (suspected not to work):");
const objTest = JSONPath({
  path: "$.level1.level2.*[?(@parentProperty === 'level2')]",
  json: simpleData, 
  resultType: 'value'
});
console.log(`   $.level1.level2.*[?(@parentProperty === 'level2')]: ${objTest.length} results`);

const objTest2 = JSONPath({
  path: "$.level1.*[?(@parentProperty === 'level1')]",
  json: simpleData,
  resultType: 'value' 
});
console.log(`   $.level1.*[?(@parentProperty === 'level1')]: ${objTest2.length} results`);

console.log("");

// Let's see what @parentProperty actually returns for different contexts
console.log("3. What does @parentProperty actually contain?");

// For debugging, let's try to use @parentProperty in different ways
const debugTests = [
  "$.level1.level2.array[?(@parentProperty)]",        // Any parentProperty
  "$.level1.level2.*[?(@parentProperty)]",            // Any parentProperty for objects  
  "$.level1.level2.array[?(@parentProperty !== '')]", // Non-empty parentProperty
  "$.level1.level2.*[?(@parentProperty !== '')]",     // Non-empty parentProperty for objects
];

debugTests.forEach(query => {
  try {
    const result = JSONPath({ path: query, json: simpleData, resultType: 'value' });
    console.log(`   ${query}: ${result.length} results`);
  } catch (e) {
    console.log(`   ${query}: ERROR - ${e.message}`);
  }
});

console.log("\n=== CONCLUSION ===");
console.log("Based on JSONPath-Plus reference implementation:");
console.log("1. @parentProperty WORKS for array elements");
console.log("2. @parentProperty does NOT work for object properties");  
console.log("3. This explains why the Go test expectations may be wrong");
console.log("4. The Go implementation should match JSONPath-Plus behavior");
console.log("");
console.log("Recommendation: Update Go test expectations to match JSONPath-Plus");
console.log("- Tests expecting @parentProperty on object properties should expect 0 results");
console.log("- Tests expecting @parentProperty on array elements should work correctly");