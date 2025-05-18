# adaptlimit

[English](README.md) | [中文](README_zh-CN.md)

一个轻量级的 Go 语言自适应并发限流器，灵感来自 Netflix 的 concurrency-limits 库。它能够根据应用程序性能和系统负载自动调整并发限制。

## 特性

- 🚀 使用 AIMD 算法的自适应并发控制
- 🔄 基于请求成功/失败的自动调节限制
- 🎯 HTTP 中间件中精确的路由级限流
- 🌟 易于与标准 Go HTTP 和 Gin 框架集成
- ⚡ 核心功能零外部依赖
- 🛠️ 通过函数选项模式进行配置

## 安装

```bash
go get github.com/yourusername/adaptlimit
```

## 快速开始

```go
// 创建一个初始限制为 10 的限流器
limiter := NewLimiter(10,
    WithMinLimit(5),    // 最小并发数
    WithMaxLimit(20),   // 最大并发数
)

// 在 Gin 中使用
r := gin.Default()
r.GET("/limited", limiter.GinMiddleware(), func(c *gin.Context) {
    // 你的处理逻辑
})
```


## 贡献

欢迎贡献！请随时提交 Pull Request。

## 许可证

MIT 许可证 