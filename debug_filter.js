const { JSONPath } = require('jsonpath-plus');

const data = {"store":{"book":[{"title":"Book0"},{"title":"Book1"},{"title":"Book2"}]}};

console.log('Testing: $..book[?(@property !== 0)]');

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
  
  console.log('JS Results:');
  console.log('Count:', values.length);
  console.log('Values:', JSON.stringify(values));
  console.log('Paths:', JSON.stringify(paths));
} catch (error) {
  console.log('JS Error:', error.message);
}

console.log('\nFor comparison, testing: $..book[*]');
const allBooks = JSONPath({path: '$..book[*]', json: data, resultType: 'value'});
const allBooksPaths = JSONPath({path: '$..book[*]', json: data, resultType: 'path'});
console.log('All books count:', allBooks.length);
console.log('All books paths:', JSON.stringify(allBooksPaths));