package main

import (
	"log"
	"net/http"
	"time"

	"github.com/GGboya/adaptlimit"
	"github.com/gin-gonic/gin"
)

func main() {
	// 创建一个限流器，初始并发限制为 5
	limiter := adaptlimit.NewLimiter(5,
		adaptlimit.WithMinLimit(2),        // 最小并发数
		adaptlimit.WithMaxLimit(10),       // 最大并发数
		adaptlimit.WithIncreaseAmount(1),  // 每次成功增加 1
		adaptlimit.WithDecreaseRatio(0.9), // 失败时降低到 90%
	)

	// 创建 Gin 引擎
	r := gin.Default()

	// 不限流的路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "This route is not rate limited",
		})
	})

	// 使用限流中间件的路由
	r.GET("/limited", limiter.GinMiddleware(), func(c *gin.Context) {
		time.Sleep(time.Millisecond * 100)
		c.JSON(http.StatusOK, gin.H{
			"message":  "This route is rate limited",
			"limit":    limiter.GetLimit(),
			"inflight": limiter.GetInflight(),
		})
	})

	// 另一个限流的路由
	r.GET("/metrics", limiter.GinMiddleware(), func(c *gin.Context) {
		time.Sleep(time.Millisecond * 100)
		c.JSON(http.StatusOK, gin.H{
			"message":  "This metrics route is rate limited",
			"limit":    limiter.GetLimit(),
			"inflight": limiter.GetInflight(),
		})
	})

	// 启动服务器
	log.Println("Starting Gin server on :8080")
	log.Println("Routes:")
	log.Println("  - / (no limit)")
	log.Println("  - /limited (limited)")
	log.Println("  - /metrics (limited)")

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
