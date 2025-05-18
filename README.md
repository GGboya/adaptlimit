# adaptlimit

[English](README.md) | [中文](README_zh-CN.md)

A lightweight, adaptive concurrency limiter for Go, inspired by Netflix's concurrency-limits library. It automatically adjusts concurrency limits based on application performance and system load.

## Features

- 🚀 Adaptive concurrency control using AIMD algorithm
- 🔄 Self-tuning limits based on request success/failure
- 🎯 Precise per-route rate limiting in HTTP middleware
- 🌟 Easy integration with standard Go HTTP and Gin framework
- ⚡ Zero external dependencies for core functionality
- 🛠️ Configurable through functional options pattern

## Installation

```bash
go get github.com/yourusername/adaptlimit
```

## Quick Start

```go
// Create a limiter with initial limit of 10
limiter := NewLimiter(10,
    WithMinLimit(5),
    WithMaxLimit(20),
)

// Use with Gin
r := gin.Default()
r.GET("/limited", limiter.GinMiddleware(), func(c *gin.Context) {
    // Your handler logic
})
```


## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License 