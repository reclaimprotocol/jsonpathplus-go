const { JSONPath } = require('jsonpath-plus');

const data = {"store":{"book":[{"title":"Book0"},{"title":"Book1"},{"title":"Book2"}]}};

console.log('=== Testing String Index Access ===');

// Test different index access patterns
const tests = [
  "$.store.book[0]",           // Direct numeric index
  "$.store.book['0']",         // String index  
  "$['store']['book'][0]",     // Bracket notation with numeric
  "$['store']['book']['0']",   // Bracket notation with string
  "$.store['book'][0]",        // Mixed notation
];

tests.forEach(query => {
  console.log(`\nQuery: ${query}`);
  try {
    const values = JSONPath({path: query, json: data, resultType: 'value'});
    const paths = JSONPath({path: query, json: data, resultType: 'path'});
    console.log(`JS Count: ${values.length}`);
    console.log(`JS Values: ${JSON.stringify(values)}`);
    console.log(`JS Paths: ${JSON.stringify(paths)}`);
  } catch (error) {
    console.log(`JS Error: ${error.message}`);
  }
});