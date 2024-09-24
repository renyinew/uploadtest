package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// 随机数据生成器：返回一个伪造的随机字节流
type RandomReader struct{}

func (r *RandomReader) Read(p []byte) (n int, err error) {
	// 填充伪随机数据
	for i := range p {
		p[i] = byte(rand.Intn(256))
	}
	return len(p), nil
}

// 模拟不同的客户端 HTTP headers
func generateRandomHeaders(req *http.Request) {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (iPad; CPU OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
	}

	// 设置随机的 User-Agent
	randomIndex := rand.Intn(len(userAgents))
	req.Header.Set("User-Agent", userAgents[randomIndex])

	// 设置其他动态 headers
	req.Header.Set("X-Request-ID", fmt.Sprintf("%d", rand.Intn(100000)))   // 随机 Request ID
	req.Header.Set("X-Client-Session", fmt.Sprintf("%d", rand.Intn(10000))) // 随机 Client Session
}

func sendOneGBData(client *http.Client, url string) error {
	// 1GB 数据大小
	const dataSize = 1024 * 1024 * 1024

	// 创建随机数据流
	randomData := io.LimitReader(&RandomReader{}, dataSize)

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", url, randomData)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	// 生成随机的客户端 headers
	generateRandomHeaders(req)

	// 开始计时
	start := time.Now()

	// 发送请求并处理响应
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 计算传输时间和速度
	duration := time.Since(start).Seconds()
	speed := float64(dataSize) / (1024 * 1024) / duration // 计算传输速度 MB/s

	log.Printf("Upload completed in %.2f seconds. Speed: %.2f MB/s\n", duration, speed)
	log.Printf("Response: %s\n", resp.Status)
	return nil
}

func main() {
	// 定义一个参数用于传递服务器 IP 地址
	ip := flag.String("ip", "localhost", "IP address of the server")
	port := flag.String("port", "18080", "Port of the server")
	iterations := flag.Int("n", 5, "Number of times to send 1GB data")

	// 解析命令行参数
	flag.Parse()

	// 构建服务器 URL
	serverURL := fmt.Sprintf("http://%s:%s/upload", *ip, *port)
	log.Printf("Sending data to %s\n", serverURL)

	client := &http.Client{}

	// 循环多次发送 1GB 数据
	for i := 0; i < *iterations; i++ { // 通过参数控制发送次数
		log.Printf("Sending 1GB of data... (Round %d)\n", i+1)
		err := sendOneGBData(client, serverURL)
		if err != nil {
			log.Fatalf("Failed to send data: %v", err)
		}
		time.Sleep(1 * time.Second) // 可选的延时，模拟不同的发送时间
	}
}
