# v1.2.0 - Nested Filter Property Access Support

## ✨ New Features

- **Nested Filter Property Access**: Added support for deeply nested property existence checks in filter expressions
  - Filters can now check for nested properties like `?(@.returnValue.returnValue.Contact)`
  - Enables complex queries: `$.actions[?(@.returnValue.returnValue.Contact)].returnValue.returnValue.Contact.LastName`
  - Full support for multi-level property paths in filter existence checks

## 🐛 Bug Fixes

- **Filter Existence Regex**: Updated the `tryExistenceFilter` regex pattern to support nested property paths
  - Changed from `/^\.(\w+)$/` to `/^\.([a-zA-Z_]\w*(?:\.[a-zA-Z_]\w*)*)$/`
  - Now correctly handles any depth of nested properties

## 🧪 Testing

- Added 2 new test cases for nested filter property access
- All 55 test cases passing (100% JavaScript compatibility maintained)
- Test data includes realistic nested structures (salesforce_actions)

## 📦 Configuration

- Updated `.golangci.yml` to include `version: 2` for compatibility with newer golangci-lint versions

## 🔗 Links

- Test cases demonstrating the new feature in [`testcases.json`](https://github.com/reclaimprotocol/jsonpathplus-go/blob/main/tests/shared/testcases.json)
- Implementation in [`filters.go`](https://github.com/reclaimprotocol/jsonpathplus-go/blob/main/internal/filters/filters.go)

## 🙏 Contributors

Thank you to all contributors who helped with this release!