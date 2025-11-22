# v1.1.5 - Parser Fix for `.[` Syntax

## 🐛 Bug Fixes

- **Fixed Parser Issue**: Resolved an issue where the parser would fail when encountering a dot followed immediately by a bracket (e.g., `$.data.[?(@.field)]`). This syntax is now correctly handled as a bracket access.

## 🧪 Testing

- Added specific test cases for `.[` syntax with filter expressions.
- Verified 100% compatibility with the JavaScript `jsonpath-plus` implementation.
- Full regression suite passed.

## 📦 Installation

```bash
go get github.com/reclaimprotocol/jsonpathplus-go@v1.1.5
```