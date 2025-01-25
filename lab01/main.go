package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func handleLogic(start, end int, wg *sync.WaitGroup) {
	defer wg.Done()

	// Thực hiện công việc trong khoảng start đến end
	for i := start; i <= end; i++ {
		// giả sử tính toán mất thời gian 10 nanosecond
		time.Sleep(10 * time.Nanosecond)
	}
}

func testRoutines(numRoutines, totalWork int) time.Duration {
	var wg sync.WaitGroup
	wg.Add(numRoutines)

	startTime := time.Now()
	workPerRoutine := totalWork / numRoutines
	remainingWork := totalWork % numRoutines

	// Chia công việc cho các goroutines sao cho phần việc gần như đồng đều
	for i := 0; i < numRoutines; i++ {
		start := i*workPerRoutine + 1
		end := (i + 1) * workPerRoutine

		// Nếu có phần công việc dư, cấp thêm 1 nhiệm vụ cho những goroutines cuối
		if i == numRoutines-1 {
			end += remainingWork
		}

		go handleLogic(start, end, &wg)
	}

	wg.Wait() // Đợi tất cả goroutines hoàn thành

	return time.Since(startTime)
}

func main() {
	cores := runtime.NumCPU()
	runtime.GOMAXPROCS(cores)
	fmt.Printf("Số CPU: %d\n", cores)

	// Tạo danh sách các số luồng từ 1 đến 50, mỗi lần tăng 1
	for threads := 1; threads <= 50; threads++ {
		// Tính toán tổng công việc là 8 triệu (ví dụ) chia cho số threads
		duration := testRoutines(threads, 8000000)
		fmt.Printf("%d threads - Time: %v\n", threads, duration)
	}
}
