const JSONPath = require('jsonpath-plus');

const data = {
  "store": {
    "book": [
      {"category": "reference", "title": "Book0"},
      {"category": "fiction", "title": "Book1"} 
    ],
    "info": {
      "0": "string_zero",
      "1": "string_one"
    }
  }
};

console.log('=== Testing @property type behavior in JS ===');

const tests = [
  // Array index tests
  {query: "$..book[?(@property === 0)]", desc: "Array: @property === 0"},
  {query: "$..book[?(@property === '0')]", desc: "Array: @property === '0'"},
  {query: "$..book[?(@property !== 0)]", desc: "Array: @property !== 0"},
  {query: "$..book[?(@property != 0)]", desc: "Array: @property != 0"},
  
  // Object key tests  
  {query: "$..info[?(@property === 0)]", desc: "Object: @property === 0"},
  {query: "$..info[?(@property === '0')]", desc: "Object: @property === '0'"},
  {query: "$..info[?(@property !== 0)]", desc: "Object: @property !== 0"},
  {query: "$..info[?(@property != 0)]", desc: "Object: @property != 0"},
];

for (const test of tests) {
  console.log(`\n${test.desc}`);
  console.log(`Query: ${test.query}`);
  try {
    const results = JSONPath.JSONPath({path: test.query, json: data, resultType: 'all'});
    console.log(`Count: ${results.length}`);
    if (results.length > 0) {
      results.forEach((r, i) => {
        console.log(`  [${i}] Path: ${r.path}, Value: ${JSON.stringify(r.value)}`);
      });
    }
  } catch (error) {
    console.log(`Error: ${error.message}`);
  }
}