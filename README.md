# JSONPath-Plus Go

[![CI](https://github.com/reclaimprotocol/jsonpathplus-go/actions/workflows/ci.yml/badge.svg)](https://github.com/reclaimprotocol/jsonpathplus-go/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/reclaimprotocol/jsonpathplus-go)](https://goreportcard.com/report/github.com/reclaimprotocol/jsonpathplus-go)
[![GoDoc](https://godoc.org/github.com/reclaimprotocol/jsonpathplus-go?status.svg)](https://godoc.org/github.com/reclaimprotocol/jsonpathplus-go)

A high-performance Go implementation of JSONPath with **string character position tracking** - perfect for JSON editors, linters, and applications requiring precise location information.

## 🚀 Quick Start

```bash
go get github.com/reclaimprotocol/jsonpathplus-go
```

```go
package main

import (
    "fmt"
    jp "github.com/reclaimprotocol/jsonpathplus-go"
)

func main() {
    jsonStr := `{"users":[{"name":"Alice","age":30},{"name":"Bob","age":25}]}`
    
    // Query with character position tracking
    results, err := jp.Query("$.users[*].name", jsonStr)
    if err != nil {
        panic(err)
    }
    
    for _, result := range results {
        fmt.Printf("Value: %v, Position: %d, Length: %d\n", 
            result.Value, result.OriginalIndex, result.Length)
    }
}
```

## ✨ Features

- 🎯 **String Position Tracking** - Get exact character positions in original JSON
- 🏭 **Production Ready** - Built-in caching, logging, metrics, and security
- 🧵 **Thread Safe** - Concurrent operations with context support
- 🔒 **Secure** - Input validation and rate limiting
- ⚡ **High Performance** - Optimized with LRU caching and minimal allocations
- 📝 **JSONPath-Plus Compatible** - Full feature compatibility

## 📁 Project Structure

```
├── README.md                    # Main documentation
├── go.mod                      # Go module configuration
├── *.go                        # Core library source code
├── cmd/                        # Command line tools and examples
│   ├── basic/                  # Basic usage examples
│   ├── production/             # Production setup examples  
│   └── showcase/               # Feature demonstration
├── tests/                      # All test files
│   ├── *_test.go              # Unit tests
│   └── benchmarks/            # Performance benchmarks
├── docs/                       # Documentation
│   ├── README.md              # Detailed docs
│   └── *.md                   # Additional documentation
└── .github/                    # CI/CD configuration
    └── workflows/ci.yml        # GitHub Actions
```

## 🔧 Advanced Usage

### Production Engine

```go
engine, err := jp.NewEngine(jp.DefaultConfig())
if err != nil {
    log.Fatal(err)
}
defer engine.Close()

// Thread-safe queries with timeout
results, err := engine.QueryData("$.store.book[*]", data)
```

### String Position Tracking

```go
jsonStr := `{"id": 123, "name": "test"}`
results, err := jp.Query("$.name", jsonStr)

// Result contains:
// - Value: "test" 
// - OriginalIndex: 15 (character position of "name" key)
// - Length: 6 (length of "name" in JSON)
// - Path: "$.name"
```

## 📊 Performance

```
BenchmarkSimplePath-12                  1,761,480 ops    683.0 ns/op
BenchmarkStringIndexPreservation-12       342,960 ops  3,647.0 ns/op  
BenchmarkCachedQuery-12                     6,973 ops 171,663.0 ns/op
```

## 🧪 Testing

```bash
go test -v ./tests/...           # Run all tests
go test -bench=. ./tests/...     # Run benchmarks  
go test -race ./tests/...        # Race condition testing
```

## 📖 Examples

See [`cmd/`](cmd/) directory for comprehensive examples:
- [`cmd/basic/`](cmd/basic/) - Basic JSONPath operations
- [`cmd/production/`](cmd/production/) - Production configuration
- [`cmd/showcase/`](cmd/showcase/) - Advanced features demo

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by [JSONPath-Plus](https://github.com/JSONPath-Plus/JSONPath) JavaScript library
- Built with ❤️ for the Go community

---

**⚡ Generated with [Claude Code](https://claude.ai/code)**