const { JSONPath } = require('jsonpath-plus');

// Test data
const testData = {
  "store": {
    "book": [
      {"category": "reference", "author": "Nigel Rees", "title": "Sayings of the Century", "price": 8.95},
      {"category": "fiction", "author": "Evelyn Waugh", "title": "Sword of Honour", "price": 12.99},
      {"category": "fiction", "author": "Herman Melville", "title": "Moby Dick", "isbn": "0-553-21311-3", "price": 8.99}
    ]
  },
  "orders": [
    {"id": "ORD001", "items": [{"product": "laptop", "price": 999.99}, {"product": "mouse", "price": 29.99}]},
    {"id": "ORD002", "items": [{"product": "keyboard", "price": 79.99}]},
    {"id": "ORD003", "items": [{"product": "laptop", "price": 1299.99}]}
  ],
  "users": ["alice", "bob", "charlie"],
  "mixed": [42, "hello", true, null, {"key": "value"}]
};

// Key test queries
const queries = [
  '$.store.book[*].author',                             // Basic wildcard
  '$.store.book[?(@.isbn)]',                           // Filter existence
  '$.orders[?(@.items[?(@.product === "laptop")])]',   // Nested filter
  '$.orders[?(@.items[*].product === "laptop")]',      // Wildcard in filter (should error)
  '$.mixed[?(@.length > 3)]',                         // Length on null (should error)
  '$.users[-1]'                                       // Negative index
];

console.log('JavaScript JSONPath-Plus Results:');
console.log('==============================');

queries.forEach(query => {
  console.log(`\nQuery: ${query}`);
  try {
    const results = JSONPath({path: query, json: testData});
    console.log(`✓ Count: ${results.length}`);
    if (results.length > 0) {
      console.log(`  First: ${JSON.stringify(results[0])}`);
    }
  } catch (error) {
    console.log(`✗ Error: ${error.message}`);
  }
});