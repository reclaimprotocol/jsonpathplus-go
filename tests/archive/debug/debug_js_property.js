const JSONPath = require('jsonpath-plus');

const data = {
  "store": {
    "book": [
      {"category": "reference", "author": "Nigel Rees", "title": "Sayings of the Century", "price": 8.95},
      {"category": "fiction", "author": "Evelyn Waugh", "title": "Sword of Honour", "price": 12.99},
      {"category": "fiction", "author": "Herman Melville", "title": "Moby Dick", "isbn": "0-553-21311-3", "price": 8.99},
      {"category": "fiction", "author": "J. R. R. Tolkien", "title": "The Lord of the Rings", "isbn": "0-395-19395-8", "price": 22.99}
    ],
    "bicycle": {"color": "red", "price": 19.95}
  }
};

console.log('=== Testing JS Property Filter Queries ===');

const queries = [
  "$..book[*]",                        // All books to see @property values
  "$..book[?(@property !== 0)]",       // The failing query
  "$..book[?(@property != '0')]",      // Alternative with string
  "$..book[?(@property > 0)]",         // Alternative comparison
  "$..*[?(@property === 'price' && @ !== 8.95)]", // The other failing query
];

for (const query of queries) {
  console.log(`\nQuery: ${query}`);
  try {
    const results = JSONPath.JSONPath({path: query, json: data, resultType: 'all'});
    console.log(`JS Count: ${results.length}`);
    
    const paths = results.map(r => r.path);
    console.log(`JS Paths: [${paths.map(p => `"${p}"`).join(', ')}]`);
    
    if (results.length > 0 && results.length <= 5) {
      const values = results.map(r => r.value);
      console.log(`JS Values: [${values.map(v => JSON.stringify(v)).join(', ')}]`);
    }
    
    // For the book wildcard query, show what @property evaluates to
    if (query === "$..book[*]") {
      results.forEach((r, i) => {
        console.log(`  [${i}] Path: ${r.path}, ParentProperty: '${r.parentProperty}'`);
      });
    }
  } catch (error) {
    console.log(`JS Error: ${error.message}`);
  }
}