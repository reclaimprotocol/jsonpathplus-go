const { JSONPath } = require('jsonpath-plus');

const testData = {
  "users": [
    {"name": "Alice", "age": 25, "active": true},
    {"name": "Bob", "age": 30, "active": false}, 
    {"name": "Charlie", "age": 35, "active": true}
  ]
};

// Test different whitespace variations
const queries = [
  '$.users[?(@.name === "Alice")]',        // No spaces
  '$.users[?(@.name==="Alice")]',          // No spaces at all
  '$.users[?(@.name === "Alice")]',        // Normal spaces
  '$.users[?(@.name  ===  "Alice")]',      // Extra spaces
  '$.users[?( @.name === "Alice" )]',      // Spaces around expression
  '$.users[?(@.name === \'Alice\')]',      // Single quotes
];

console.log('Testing whitespace sensitivity in JavaScript JSONPath-Plus:');
console.log('='.repeat(60));

queries.forEach((query, i) => {
  console.log(`\nTest ${i + 1}: ${query}`);
  try {
    const results = JSONPath({path: query, json: testData});
    console.log(`✓ Results: ${results.length} (${results.length > 0 ? results[0].name : 'none'})`);
  } catch (error) {
    console.log(`✗ Error: ${error.message}`);
  }
});

// Test the specific failing nested filter with different whitespace
console.log('\n' + '='.repeat(60));
console.log('Testing nested filter whitespace variations:');
console.log('='.repeat(60));

const orderData = {
  "orders": [
    {"id": "ORD001", "items": [{"product": "laptop", "price": 999.99}, {"product": "mouse", "price": 29.99}]},
    {"id": "ORD002", "items": [{"product": "keyboard", "price": 79.99}]},
    {"id": "ORD003", "items": [{"product": "laptop", "price": 1299.99}]}
  ]
};

const nestedQueries = [
  '$.orders[?(@.items[?(@.product === "laptop")])]',
  '$.orders[?(@.items[?(@.product==="laptop")])]',
  '$.orders[?(@.items[?( @.product === "laptop" )])]',
  '$.orders[?(@.items[?(@.product === \'laptop\')])]'
];

nestedQueries.forEach((query, i) => {
  console.log(`\nNested Test ${i + 1}: ${query}`);
  try {
    const results = JSONPath({path: query, json: orderData});
    console.log(`✓ Results: ${results.length} orders`);
    if (results.length > 0) {
      console.log(`  First result ID: ${results[0].id}`);
    }
  } catch (error) {
    console.log(`✗ Error: ${error.message}`);
  }
});