package main

import (
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

func WebThread() {
	// Start the web server
	router := gin.Default()
	router.Use(gin.Recovery())
	// 实际应该将全局管道放在router作用域外，这样才能控制并发数
	var limiter = make(chan struct{}, 3)
	// 可以通过i的值验证是否对每个请求单独创建了一个线程
	var i int64
	router.GET("/thread", func(c *gin.Context) {
		// 试图创建一个全局管道，用于控制并发数
		// var limiter = make(chan struct{}, 3)
		// i的值也可用于验证是否对每个请求单独创建了一个线程
		// var i int64
		go func() {
			index := atomic.AddInt64(&i, 1)
			limiter <- struct{}{}
			defer func() {
				<-limiter
			}()
			// 模拟耗时操作
			for j := 0; j < 3; j++ {
				time.Sleep(1 * time.Second)
				fmt.Printf("创建线程%v: %v\n", index, time.Now())
			}
		}()
	})
	router.Run() // listen and serve on
}

func Request() {
	//发送web请求
	for i := 0; i < 5; i++ {
		resp, err := http.Get("http://localhost:8080/thread")
		if err != nil {
			fmt.Println("http.Get err:", err)
			return
		}
		defer resp.Body.Close()
		// 读取响应内容
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			// 处理读取错误
			fmt.Println("读取响应失败:", err)
			return
		}
		// 打印响应状态码和正文
		fmt.Println("状态码:", resp.StatusCode)
		fmt.Println("响应正文:", string(body))
	}
}

func main() {
	go func() {
		time.Sleep(3 * time.Second)
		Request()
	}()
	WebThread()
}
