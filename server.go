package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// 获取客户端的 IP 地址
func getClientIP(r *http.Request) string {
	// 先检查 X-Forwarded-For header，防止客户端通过代理或负载均衡发送请求
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		// 可能包含多个 IP 地址，以逗号分隔，取第一个
		return strings.Split(forwarded, ",")[0]
	}
	// 如果没有 X-Forwarded-For，直接使用 RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // 返回原始的 RemoteAddr
	}
	return ip
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// 获取客户端 IP 地址
	clientIP := getClientIP(r)

	// 创建缓冲区，分块读取数据
	buffer := make([]byte, 1024*1024) // 1MB 缓冲区
	totalBytes := 0

	start := time.Now() // 开始计时

	for {
		n, err := r.Body.Read(buffer)
		if n > 0 {
			totalBytes += n
			duration := time.Since(start).Seconds() // 计算接收到现在的时间
			speed := float64(totalBytes) / (1024 * 1024) / duration // 计算接收速度 MB/s
			fmt.Printf("Client IP: %s - Received %d MB... Speed: %.2f MB/s\n", clientIP, totalBytes/(1024*1024), speed)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
	}

	totalMB := totalBytes / (1024 * 1024)
	fmt.Printf("Client IP: %s - Total received: %d MB in %.2f seconds. Average speed: %.2f MB/s\n",
		clientIP, totalMB, time.Since(start).Seconds(), float64(totalBytes)/(1024*1024)/time.Since(start).Seconds())

	w.Write([]byte("Upload complete"))
}

func main() {
	http.HandleFunc("/upload", uploadHandler) // 路由处理
	fmt.Println("Server is listening on port 18081...")
	log.Fatal(http.ListenAndServe(":18081", nil)) // 启动服务器，监听18081端口
}
