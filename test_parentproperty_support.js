#!/usr/bin/env node

const { JSONPath } = require('jsonpath-plus');

// Test if @parentProperty is supported in JSONPath-Plus at all
const testData = {
  "parent": {
    "child1": {"value": 1},
    "child2": {"value": 2},
    "items": [
      {"name": "item1"},
      {"name": "item2"} 
    ]
  }
};

console.log("=== Testing @parentProperty Support in JSONPath-Plus ===\n");

const testCases = [
  {
    name: "Basic @parentProperty test",
    query: "$.parent.*[?(@parentProperty === 'parent')]",
    description: "Should find children where parent property is 'parent'"
  },
  {
    name: "Array element @parentProperty",  
    query: "$.parent.items[?(@parentProperty === 'items')]",
    description: "Should find array elements where parent property is 'items'"
  },
  {
    name: "@property test for comparison",
    query: "$.parent.*[?(@property === 'child1')]", 
    description: "Should find element with property name 'child1'"
  },
  {
    name: "Array @property test",
    query: "$.parent.items[?(@property === 0)]",
    description: "Should find first array element" 
  }
];

testCases.forEach((testCase, index) => {
  console.log(`${index + 1}. ${testCase.name}`);
  console.log(`   Query: ${testCase.query}`);
  console.log(`   Description: ${testCase.description}`);
  
  try {
    const result = JSONPath({
      path: testCase.query,
      json: testData,
      resultType: 'value'
    });
    
    console.log(`   Results: ${result.length} items`);
    result.forEach((val, i) => {
      console.log(`     [${i}] ${JSON.stringify(val)}`);
    });
    
    if (result.length === 0) {
      console.log("   ⚠️  No results - may indicate @parentProperty is not supported");
    }
    
  } catch (error) {
    console.log(`   ❌ ERROR: ${error.message}`);
    if (error.message.includes('parentProperty')) {
      console.log("   💡 This suggests @parentProperty is not supported in JSONPath-Plus");
    }
  }
  
  console.log("");
});

// Check JSONPath-Plus version and documentation
console.log("=== JSONPath-Plus Version Info ===");
try {
  const pkg = require('jsonpath-plus/package.json');
  console.log(`Version: ${pkg.version}`);
  console.log(`Description: ${pkg.description}`);
} catch (e) {
  console.log("Could not retrieve version info");
}

console.log("\n=== Additional Context Tests ===");

// Test what context variables are actually available
const contextTests = [
  "@",        // Current item
  "@.value",  // Current item property  
  "@property", // Property name
  "@parent",   // Parent object
  "@parentProperty", // Parent property name
  "@root",     // Root object
  "@path",     // Current path
];

contextTests.forEach(ctx => {
  try {
    const result = JSONPath({
      path: `$.parent.child1[?(${ctx})]`,
      json: testData,
      resultType: 'value'
    });
    console.log(`${ctx}: ${result.length > 0 ? '✅ Recognized' : '❌ Not recognized or false'}`);
  } catch (error) {
    if (error.message.includes(ctx.replace('@', ''))) {
      console.log(`${ctx}: ❌ Not supported (${error.message.split(':')[0]})`);
    } else {
      console.log(`${ctx}: ❌ Error`);
    }
  }
});