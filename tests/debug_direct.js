const { JSONPath } = require('jsonpath-plus');

const data = {"store":{"book":[{"title":"Book0"},{"title":"Book1"},{"title":"Book2"}]}};

console.log('=== Direct JS Test ===');
console.log('Query: $..book[?(@property !== 0)]');
console.log('Data:', JSON.stringify(data));

try {
  const values = JSONPath({
    path: '$..book[?(@property !== 0)]',
    json: data,
    resultType: 'value'
  });
  
  const paths = JSONPath({
    path: '$..book[?(@property !== 0)]', 
    json: data,
    resultType: 'path'
  });
  
  console.log('\nJS Results:');
  console.log('Count:', values.length);
  console.log('Values:', JSON.stringify(values));
  console.log('Paths:', JSON.stringify(paths));
} catch (error) {
  console.log('JS Error:', error.message);
}