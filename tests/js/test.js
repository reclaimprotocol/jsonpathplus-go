#!/usr/bin/env node

const { JSONPath } = require('jsonpath-plus');

if (process.argv.length < 4) {
  console.error(`Usage: ${process.argv[1]} <jsonpath-or-file> <json-data-or-file> [--query-file] [--data-file]`);
  process.exit(1);
}

const fs = require('fs');
const jsonpathOrFile = process.argv[2];
const jsonDataOrFile = process.argv[3];
let jsonpath, data;

// Check flags
const queryFromFile = process.argv.includes('--query-file');
const dataFromFile = process.argv.includes('--data-file') || process.argv.includes('--file');

try {
  // Read JSONPath query
  if (queryFromFile) {
    jsonpath = fs.readFileSync(jsonpathOrFile, 'utf8').trim();
  } else {
    jsonpath = jsonpathOrFile;
  }
  
  // Read JSON data
  if (dataFromFile) {
    const fileData = fs.readFileSync(jsonDataOrFile, 'utf8');
    data = JSON.parse(fileData);
  } else {
    data = JSON.parse(jsonDataOrFile);
  }
  
  const values = JSONPath({
    path: jsonpath,
    json: data,
    resultType: 'value'
  });
  
  const paths = JSONPath({
    path: jsonpath, 
    json: data,
    resultType: 'path'
  });
  
  const output = {
    count: values.length,
    values: values,
    paths: paths
  };
  
  console.log(JSON.stringify(output));
} catch (error) {
  const output = {
    error: error.message,
    count: 0,
    values: [],
    paths: []
  };
  
  console.log(JSON.stringify(output));
}