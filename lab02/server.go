package main

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	counterAtomic  int64
	counterNoLock  int64
	counterMutex   int64
	counterBatch   int64
	mutex          sync.Mutex
	requestChannel = make(chan int, 10000) // channel để xử lý batch (buffered)
	workerCount    = 2
)

func main() {
	// Sử dụng tất cả CPU
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Tạo ứng dụng Fiber
	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	// Endpoint tăng counter không sử dụng lock
	app.Get("/nolock", func(c *fiber.Ctx) error {
		counterNoLock++
		return c.SendString("OK")
	})

	// Endpoint tăng counter dùng atomic
	app.Get("/atomic", func(c *fiber.Ctx) error {
		atomic.AddInt64(&counterAtomic, 1)
		return c.SendString("OK")
	})

	// Endpoint tăng counter sử dụng Mutex
	app.Get("/mutex", func(c *fiber.Ctx) error {
		mutex.Lock()
		counterMutex++
		mutex.Unlock()
		return c.SendString("OK")
	})

	// Endpoint tăng counter sử dụng batch processing
	app.Get("/batch", func(c *fiber.Ctx) error {
		select {
		case requestChannel <- 1: // Thêm request vào channel
			return c.SendString("OK")
		default:
			return c.Status(fiber.StatusTooManyRequests).SendString("Request queue is full")
		}
	})

	// Endpoint hiển thị trạng thái counters
	app.Get("/status", func(c *fiber.Ctx) error {
		status := "Counter without lock: " + strconv.FormatInt(counterNoLock, 10) +
			"\nCounter with atomic: " + strconv.FormatInt(counterAtomic, 10) +
			"\nCounter with mutex: " + strconv.FormatInt(counterMutex, 10) +
			"\nCounter with batch processing: " + strconv.FormatInt(counterBatch, 10)
		return c.SendString(status)
	})

	// Khởi tạo batchProcessor với worker count
	go batchProcessor(workerCount)

	// Bắt đầu server Fiber
	fmt.Println("Server running on port 8080")
	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}
}

// Hàm này sẽ khởi tạo một số lượng worker để xử lý requests từ channel
func batchProcessor(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go worker(i)
	}
}

// Worker để xử lý công việc từ kênh `requestChannel`
func worker(i int) {
	fmt.Println("Worker", i, "started")
	for {
		select {
		case <-requestChannel:
			// Khi có request, thực hiện công việc xử lý (ví dụ: tăng counter)
			atomic.AddInt64(&counterBatch, 1)
		}
	}
}
