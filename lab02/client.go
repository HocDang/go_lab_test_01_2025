package main

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	numRequests = 2000
	numThreads  = 500
	serverURL   = "http://108.61.178.68:8080"
)

func sendRequests(endpoint string) int32 {
	start := time.Now()

	var wg sync.WaitGroup
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	var totalSuccess int32

	requestsPerThread := numRequests / numThreads

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count_success := 0

			for j := 0; j < requestsPerThread; j++ {
				_, err := client.Get(serverURL + endpoint)
				if err == nil {
					count_success++
				} else {
					// fmt.Println(err)
				}
			}

			// Đảm bảo việc cập nhật `totalSuccess` là thread-safe
			atomic.AddInt32(&totalSuccess, int32(count_success))
		}()
	}

	// Đợi tất cả các goroutines hoàn thành
	wg.Wait()
	duration := time.Since(start)
	fmt.Printf("Total successful requests for endpoint %s: %d\n", endpoint, totalSuccess)
	fmt.Printf("Completed %s requests in %v\n", endpoint, duration)

	return totalSuccess
}

func main() {
	fmt.Println("Sending requests to server")

	time.Sleep(2 * time.Second)
	sendRequests("/nolock")

	time.Sleep(2 * time.Second)
	sendRequests("/atomic")

	time.Sleep(2 * time.Second)
	sendRequests("/mutex")

	time.Sleep(2 * time.Second)
	sendRequests("/batch")
}
