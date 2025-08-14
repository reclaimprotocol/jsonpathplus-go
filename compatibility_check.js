const { JSONPath } = require('jsonpath-plus');
const fs = require('fs');

// Test data
const testData = {
  "store": {
    "book": [
      {"category": "reference", "author": "Nigel Rees", "title": "Sayings of the Century", "price": 8.95},
      {"category": "fiction", "author": "Evelyn Waugh", "title": "Sword of Honour", "price": 12.99},
      {"category": "fiction", "author": "Herman Melville", "title": "Moby Dick", "isbn": "0-553-21311-3", "price": 8.99},
      {"category": "fiction", "author": "J. R. R. Tolkien", "title": "The Lord of the Rings", "isbn": "0-395-19395-8", "price": 22.99}
    ],
    "bicycle": {"color": "red", "price": 19.95}
  },
  "orders": [
    {"id": "ORD001", "items": [{"product": "laptop", "price": 999.99}, {"product": "mouse", "price": 29.99}]},
    {"id": "ORD002", "items": [{"product": "keyboard", "price": 79.99}]},
    {"id": "ORD003", "items": [{"product": "laptop", "price": 1299.99}]}
  ],
  "users": ["alice", "bob", "charlie", "alice"],
  "mixed": [42, "hello", true, null, {"key": "value"}, [1, 2, 3]]
};

// Test queries
const testQueries = [
  // Basic queries
  '$.store.book[*].author',
  '$..author', 
  '$.store.book[0]',
  '$.store.book[-1]', // Negative index
  
  // Filter queries
  '$.store.book[?(@.isbn)]',
  '$.store.book[?(@.price < 10)]',
  '$.orders[?(@.items[?(@.product === "laptop")])]', // Nested filter
  '$.orders[?(@.items[*].product === "laptop")]', // Wildcard in filter
  
  // Context filters  
  '$.store.book[?(@property !== 0)]',
  '$.users[?(@property === "1")]',
  
  // Function filters
  '$.mixed[?(@.length > 3)]', // Length on null should error
  
  // Path filters
  '$.users[?(@path === "$[\'users\'][1]")]',
  
  // Parent filters
  '$.store.book[?(@parent.bicycle)]'
];

function testQuery(query, data) {
  try {
    const results = JSONPath({path: query, json: data});
    return {
      success: true,
      count: results.length,
      results: results.slice(0, 3), // First 3 results for brevity
      error: null
    };
  } catch (error) {
    return {
      success: false,
      count: 0,
      results: [],
      error: error.message
    };
  }
}

console.log('JavaScript JSONPath-Plus Compatibility Check');
console.log('===========================================');

const results = {};
testQueries.forEach(query => {
  console.log(`\nTesting: ${query}`);
  const result = testQuery(query, testData);
  results[query] = result;
  
  if (result.success) {
    console.log(`✓ Count: ${result.count}`);
    if (result.results.length > 0) {
      console.log(`  Sample: ${JSON.stringify(result.results[0]).substring(0, 100)}${result.results[0] && JSON.stringify(result.results[0]).length > 100 ? '...' : ''}`);
    }
  } else {
    console.log(`✗ Error: ${result.error}`);
  }
});

// Save results for Go comparison
fs.writeFileSync('js_results.json', JSON.stringify(results, null, 2));
console.log('\n\nResults saved to js_results.json for Go comparison.');